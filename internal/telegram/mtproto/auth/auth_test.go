package auth_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/leirbagxis/FreddyBot/internal/telegram/mtproto/auth"
	"github.com/redis/go-redis/v9"
)

var errMockSave = fmt.Errorf("mock save error")

type mockAccountSaver struct {
	savedUserID       int64
	savedTelegramID   int64
	savedUsername     string
	savedFirstName    string
	saveCallCount     int
	failSave          bool
}

func (m *mockAccountSaver) SaveSession(ctx context.Context, userID int64, telegramUserID int64, username string, firstName string, sessionData []byte) error {
	m.saveCallCount++
	if m.failSave {
		return errMockSave
	}
	m.savedUserID = userID
	m.savedTelegramID = telegramUserID
	m.savedUsername = username
	m.savedFirstName = firstName
	return nil
}

func setupAuthTest(t *testing.T) (*miniredis.Miniredis, *auth.Service, *mockAccountSaver) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	saver := &mockAccountSaver{}
	// appID=0 e appHash="" fazem isConfigured() retornar false,
	// entao o servico nao tenta conectar MTProto real nos testes unitarios.
	svc := auth.NewService(redisClient, 0, "", saver)

	t.Cleanup(func() {
		redisClient.Close()
		mr.Close()
	})

	return mr, svc, saver
}

// saveState helper para criar estado no Redis diretamente (simula o que SendCode faria).
func saveState(t *testing.T, mr *miniredis.Miniredis, userID int64, phoneNumber string) {
	t.Helper()
	key := fmt.Sprintf("mtproto_auth:%d", userID)
	val := fmt.Sprintf(`{"user_id":%d,"phone_number":"%s","phone_code_hash":"test_hash_%d"}`, userID, phoneNumber, userID)
	mr.Set(key, val)
	mr.SetTTL(key, 5*time.Minute)
}

func TestAuthService_SendCode_Stub(t *testing.T) {
	mr, svc, _ := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	status, err := svc.SendCode(ctx, 100, "+5511999999999")
	if err != nil {
		t.Fatalf("SendCode() error = %v", err)
	}
	// Sem MTProto, o servico usa stub e retorna "code"
	if status.Step != "code" {
		t.Errorf("Step = %q, want %q", status.Step, "code")
	}
	if status.Error != "" {
		t.Errorf("Error = %q, want empty", status.Error)
	}

	// Deve salvar estado no Redis (stub tambem salva)
	if !mr.Exists("mtproto_auth:100") {
		t.Error("Auth state should exist in Redis after SendCode stub")
	}
}

func TestAuthService_SendCode_PhoneTooLong(t *testing.T) {
	mr, svc, _ := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	longPhone := "+55" + string(make([]byte, 30))
	status, err := svc.SendCode(ctx, 102, longPhone)
	if err != nil {
		t.Fatalf("SendCode() for long phone error = %v", err)
	}
	if status.Step != "error" {
		t.Errorf("Step = %q, want %q for invalid phone", status.Step, "error")
	}
}

func TestAuthService_VerifyCode_NoState(t *testing.T) {
	mr, svc, _ := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	// Tentar verificar codigo sem ter iniciado o fluxo
	status, err := svc.VerifyCode(ctx, 200, "12345")
	if err != nil {
		t.Fatalf("VerifyCode() error = %v", err)
	}
	if status.Step != "error" {
		t.Errorf("Step = %q, want %q", status.Step, "error")
	}
	if status.Error != "Sessão expirada. Inicie a conexão novamente." {
		t.Errorf("Error = %q, want expiration message", status.Error)
	}
}

func TestAuthService_FullFlow_NoPassword(t *testing.T) {
	mr, svc, saver := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	// Simular que SendCode foi chamado (criar estado no Redis)
	saveState(t, mr, 300, "+5511999999999")

	// VerifyCode com o estado existente — stub retorna "done" e salva sessao
	status, err := svc.VerifyCode(ctx, 300, "12345")
	if err != nil {
		t.Fatalf("VerifyCode() error = %v", err)
	}
	if status.Step != "done" {
		t.Errorf("Step = %q, want %q", status.Step, "done")
	}

	// O saver deve ter sido chamado (stub salva sessao)
	if saver.saveCallCount != 1 {
		t.Errorf("saveCallCount = %d, want 1", saver.saveCallCount)
	}
	if saver.savedUserID != 300 {
		t.Errorf("savedUserID = %d, want 300", saver.savedUserID)
	}
	if saver.savedTelegramID == 0 {
		t.Error("savedTelegramID should not be 0")
	}
	if saver.savedUsername != "stub_user" {
		t.Errorf("savedUsername = %q, want 'stub_user'", saver.savedUsername)
	}
}

func TestAuthService_PasswordFlow(t *testing.T) {
	mr, svc, saver := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	// Simular estado apos SendCode + VerifyCode com 2FA detectado
	key := fmt.Sprintf("mtproto_auth:%d", 400)
	val := `{"user_id":400,"phone_number":"+5511999999999","phone_code_hash":"test_hash_400","has_password":true}`
	mr.Set(key, val)
	mr.SetTTL(key, 5*time.Minute)

	// VerifyPassword com 2FA — stub retorna "done" e salva sessao
	status, err := svc.VerifyPassword(ctx, 400, "minha_senha")
	if err != nil {
		t.Fatalf("VerifyPassword() error = %v", err)
	}
	if status.Step != "done" {
		t.Errorf("Step = %q, want %q", status.Step, "done")
	}

	// O saver deve ter sido chamado (stub salva sessao)
	if saver.saveCallCount != 1 {
		t.Errorf("saveCallCount = %d, want 1", saver.saveCallCount)
	}
	if saver.savedUserID != 400 {
		t.Errorf("savedUserID = %d, want 400", saver.savedUserID)
	}
}

func TestAuthService_VerifyPassword_NoState(t *testing.T) {
	mr, svc, _ := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	status, err := svc.VerifyPassword(ctx, 500, "senha")
	if err != nil {
		t.Fatalf("VerifyPassword() error = %v", err)
	}
	if status.Step != "error" {
		t.Errorf("Step = %q, want %q", status.Step, "error")
	}
	if status.Error != "Sessão expirada. Inicie a conexão novamente." {
		t.Errorf("Error = %q, want expiration message", status.Error)
	}
}

func TestAuthService_GetStatus_NoState(t *testing.T) {
	mr, svc, _ := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	status, err := svc.GetStatus(ctx, 600)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}
	if status.Step != "phone" {
		t.Errorf("Step = %q, want %q", status.Step, "phone")
	}
}

func TestAuthService_StateExpiration(t *testing.T) {
	mr, svc, _ := setupAuthTest(t)
	defer mr.Close()
	ctx := context.Background()

	// Salvar estado com TTL curto (manualmente no miniredis)
	key := fmt.Sprintf("mtproto_auth:%d", 700)
	val := `{"user_id":700,"phone_number":"+5511999999999","phone_code_hash":"test_hash_700"}`
	mr.Set(key, val)
	mr.SetTTL(key, 1*time.Millisecond)

	// Avancar o tempo para expirar
	mr.FastForward(10 * time.Millisecond)

	// Tentar verificar codigo com estado expirado
	status, err := svc.VerifyCode(ctx, 700, "12345")
	if err != nil {
		t.Fatalf("VerifyCode() error = %v", err)
	}
	if status.Step != "error" {
		t.Errorf("Step = %q, want %q", status.Step, "error")
	}
	if status.Error != "Sessão expirada. Inicie a conexão novamente." {
		t.Errorf("Error = %q, want expiration message", status.Error)
	}

	// GetStatus deve retornar "phone" quando nao ha estado
	status, err = svc.GetStatus(ctx, 700)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}
	if status.Step != "phone" {
		t.Errorf("Step = %q, want %q", status.Step, "phone")
	}
}
