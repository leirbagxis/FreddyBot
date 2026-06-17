package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"strconv"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/leirbagxis/FreddyBot/internal/core/crypto"
	"github.com/leirbagxis/FreddyBot/internal/database/models"
	"github.com/leirbagxis/FreddyBot/internal/database/repositories"
	"github.com/leirbagxis/FreddyBot/pkg/config"
	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

var ErrNeeds2FA = errors.New("SESSION_PASSWORD_NEEDED")

type authPhase int

const (
	phaseNone authPhase = iota
	phaseRunning
	phaseAwaitingCode
	phaseAwaiting2FA
	phaseDone
)

type authState struct {
	phone      string
	phase      authPhase
	codeCh     chan string
	passwordCh chan string
	codeSentCh chan struct{} // closed when Telegram confirms code was sent
	needs2FA   chan struct{}
	errCh      chan error
	cancel     context.CancelFunc
}

const clientIdleTimeout = 10 * time.Minute

type userClient struct {
	client   *telegram.Client
	cancel   context.CancelFunc
	lastUsed time.Time
	userID   int64
}

type channelUserAuth struct {
	phone      string
	codeCh     <-chan string
	passwordCh <-chan string
	codeSentCh chan<- struct{}
	needs2FA   chan<- struct{}
}

func mtprotoChannelID(botAPIID int64) int64 {
	s := strconv.FormatInt(botAPIID, 10)
	if strings.HasPrefix(s, "-100") {
		if cleaned, err := strconv.ParseInt(s[4:], 10, 64); err == nil {
			return cleaned
		}
	}
	return botAPIID
}

func (a *channelUserAuth) Phone(ctx context.Context) (string, error) {
	logger.Info("TGC", "Flow: Phone() returning %s", a.phone)
	return a.phone, nil
}

func (a *channelUserAuth) Password(ctx context.Context) (string, error) {
	logger.Info("TGC", "Flow: Password() called — 2FA required")
	select {
	case a.needs2FA <- struct{}{}:
	default:
	}
	select {
	case pwd := <-a.passwordCh:
		logger.Info("TGC", "Flow: Password() received 2FA input")
		return pwd, nil
	case <-ctx.Done():
		logger.Info("TGC", "Flow: Password() cancelled")
		return "", ctx.Err()
	}
}

func (a *channelUserAuth) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	logger.Info("TGC", "Flow: Code() called — code was sent via Telegram")
	select {
	case a.codeSentCh <- struct{}{}:
	default:
	}
	select {
	case code := <-a.codeCh:
		logger.Info("TGC", "Flow: Code() received user input")
		return code, nil
	case <-ctx.Done():
		logger.Info("TGC", "Flow: Code() cancelled")
		return "", ctx.Err()
	}
}

func (a *channelUserAuth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	logger.Info("TGC", "Flow: AcceptTermsOfService() called — auto-accepting")
	return nil
}

func (a *channelUserAuth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	logger.Info("TGC", "Flow: SignUp() called — not supported")
	return auth.UserInfo{}, errors.New("sign up not supported")
}

type TelegramClientService struct {
	mu           sync.Mutex
	activeStates map[int64]*authState
	activeClient map[int64]*userClient

	sessionRepo *repositories.TelegramSessionRepository
	masterKey   []byte
	apiID       int
	apiHash     string
}

func NewTelegramClientService(sessionRepo *repositories.TelegramSessionRepository) *TelegramClientService {
	masterKey := []byte(config.EncryptionKey)
	if len(masterKey) == 64 {
		if decoded, err := hex.DecodeString(config.EncryptionKey); err == nil {
			masterKey = decoded
		}
	}
	if len(masterKey) != 32 {
		logger.Error("TGC", "ENC_KEY must be 32 bytes (got %d) — use 64 hex chars or 32 raw bytes", len(masterKey))
	}

	return &TelegramClientService{
		activeStates: make(map[int64]*authState),
		activeClient: make(map[int64]*userClient),
		sessionRepo:  sessionRepo,
		masterKey:    masterKey,
		apiID:        config.TelegramAPIID,
		apiHash:      config.TelegramAPIHash,
	}
}

func (s *TelegramClientService) StartPhoneFlow(ctx context.Context, userID int64, phone string) error {
	s.mu.Lock()
	if _, exists := s.activeStates[userID]; exists {
		s.mu.Unlock()
		return fmt.Errorf("auth flow already in progress for user %d", userID)
	}

	codeCh := make(chan string, 1)
	passwordCh := make(chan string, 1)
	codeSentCh := make(chan struct{})
	needs2FA := make(chan struct{}, 1)

	state := &authState{
		phone:      phone,
		phase:      phaseRunning,
		codeCh:     codeCh,
		passwordCh: passwordCh,
		codeSentCh: codeSentCh,
		needs2FA:   needs2FA,
		errCh:      make(chan error, 1),
	}
	s.activeStates[userID] = state
	s.mu.Unlock()

	logger.Info("TGC", "Starting auth flow for user %d, phone=%s, apiID=%d", userID, phone, s.apiID)

	client := telegram.NewClient(s.apiID, s.apiHash, telegram.Options{
		SessionStorage: s.dbStorage(userID),
	})

	authCtx, cancel := context.WithCancel(context.Background())
	state.cancel = cancel

	userAuth := &channelUserAuth{
		phone:      phone,
		codeCh:     codeCh,
		passwordCh: passwordCh,
		codeSentCh: codeSentCh,
		needs2FA:   needs2FA,
	}

	go func() {
		defer func() {
			s.mu.Lock()
			if curr, ok := s.activeStates[userID]; ok && curr == state {
				delete(s.activeStates, userID)
			}
			s.mu.Unlock()
		}()

		start := time.Now()
		err := client.Run(authCtx, func(runCtx context.Context) error {
			api := client.API()
			authClient := auth.NewClient(api, rand.Reader, s.apiID, s.apiHash)
			flow := auth.NewFlow(userAuth, auth.SendCodeOptions{})
			return flow.Run(runCtx, authClient)
		})

		s.mu.Lock()
		if curr, ok := s.activeStates[userID]; ok && curr == state {
			if err != nil {
				logger.Error("TGC", "Auth flow for user %d failed: %v", userID, err)
				select {
				case state.errCh <- err:
				default:
				}
			} else {
				state.phase = phaseDone
				logger.Info("TGC", "User %d authenticated in %v", userID, time.Since(start))

				s.activeClient[userID] = &userClient{
					client:   client,
					cancel:   cancel,
					lastUsed: time.Now(),
					userID:   userID,
				}
				state.errCh <- nil
			}
		}
		s.mu.Unlock()
	}()

	// Wait for code to be sent, or error
	logger.Info("TGC", "Waiting for code to be sent for user %d...", userID)
	select {
	case <-state.codeSentCh:
		logger.Info("TGC", "Code sent successfully for user %d", userID)
		return nil
	case err := <-state.errCh:
		logger.Error("TGC", "Failed to send code for user %d: %v", userID, err)
		return fmt.Errorf("send code: %w", err)
	case <-ctx.Done():
		logger.Warn("TGC", "Auth flow cancelled for user %d", userID)
		return ctx.Err()
	}
}

func (s *TelegramClientService) SubmitCode(ctx context.Context, userID int64, code string) error {
	s.mu.Lock()
	state, ok := s.activeStates[userID]
	s.mu.Unlock()
	if !ok {
		return fmt.Errorf("no auth flow in progress")
	}

	select {
	case state.codeCh <- code:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *TelegramClientService) Submit2FA(ctx context.Context, userID int64, password string) error {
	s.mu.Lock()
	state, ok := s.activeStates[userID]
	s.mu.Unlock()
	if !ok {
		return fmt.Errorf("no auth flow for user %d", userID)
	}

	select {
	case state.passwordCh <- password:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *TelegramClientService) WaitAuthResult(ctx context.Context, userID int64) error {
	s.mu.Lock()
	state, ok := s.activeStates[userID]
	s.mu.Unlock()

	if !ok {
		saved, err := s.sessionRepo.GetByUserID(ctx, userID)
		if err == nil && saved.IsActive {
			return nil
		}
		return fmt.Errorf("no auth flow for user %d", userID)
	}

	select {
	case <-state.needs2FA:
		return ErrNeeds2FA
	case err := <-state.errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *TelegramClientService) ConnectUser(ctx context.Context, userID int64) (*telegram.Client, error) {
	s.mu.Lock()
	if uc, ok := s.activeClient[userID]; ok {
		uc.lastUsed = time.Now()
		s.mu.Unlock()
		return uc.client, nil
	}
	s.mu.Unlock()

	sess, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil || !sess.IsActive {
		return nil, fmt.Errorf("no active session for user %d", userID)
	}

	client := telegram.NewClient(s.apiID, s.apiHash, telegram.Options{
		SessionStorage: s.dbStorage(userID),
	})

	clientCtx, cancel := context.WithCancel(context.Background())
	ready := make(chan struct{})

	go func() {
		if err := client.Run(clientCtx, func(ctx context.Context) error {
			close(ready) // Session ready — safe to make API calls
			if _, err := client.API().AccountUpdateStatus(ctx, true); err != nil {
				logger.Warn("TGC", "Falha ao setar offline user %d: %v", userID, err)
			}
			<-ctx.Done()
			return ctx.Err()
		}); err != nil && err != context.Canceled {
			logger.Error("TGC", "Client run for user %d failed: %v", userID, err)
		}
	}()

	// Wait for session to be ready before returning.
	// Without this, the first API call gets "waitSession: context deadline exceeded".
	select {
	case <-ready:
	case <-ctx.Done():
		cancel()
		return nil, fmt.Errorf("wait session ready: %w", ctx.Err())
	}

	uc := &userClient{
		client:   client,
		cancel:   cancel,
		lastUsed: time.Now(),
		userID:   userID,
	}

	s.mu.Lock()
	s.activeClient[userID] = uc
	s.mu.Unlock()

	// Auto-disconnect after idle timeout
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.mu.Lock()
				if time.Since(uc.lastUsed) > clientIdleTimeout {
					uc.cancel()
					delete(s.activeClient, userID)
					s.mu.Unlock()
					logger.Info("TGC", "Auto-disconnected idle client for user %d", userID)
					return
				}
				s.mu.Unlock()
			case <-clientCtx.Done():
				return
			}
		}
	}()

	return client, nil
}

func (s *TelegramClientService) DisconnectUser(ctx context.Context, userID int64) error {
	s.mu.Lock()
	if uc, ok := s.activeClient[userID]; ok {
		uc.cancel()
		delete(s.activeClient, userID)
	}
	if state, ok := s.activeStates[userID]; ok {
		state.cancel()
		delete(s.activeStates, userID)
	}
	s.mu.Unlock()

	phoneHash := ""
	if existing, err := s.sessionRepo.GetByUserID(ctx, userID); err == nil {
		phoneHash = existing.EncryptedPhoneHash
	}

	if err := s.sessionRepo.Upsert(ctx, &models.UserTelegramSession{
		UserID:             userID,
		EncryptedSession:   "",
		EncryptedPhoneHash: phoneHash,
		IsActive:           false,
	}); err != nil {
		return err
	}

	return s.sessionRepo.SetActive(ctx, userID, false)
}

func (s *TelegramClientService) IsConnected(ctx context.Context, userID int64) (bool, error) {
	sess, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, nil
	}
	return sess.IsActive, nil
}

func (s *TelegramClientService) SendMessageAsUser(ctx context.Context, userID int64, chatID int64, text string) error {
	client, err := s.ConnectUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("connect user: %w", err)
	}

	api := client.API()
	_, err = api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:    &tg.InputPeerUser{UserID: chatID},
		Message: text,
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	logger.Info("TGC", "Message sent as user %d to %d", userID, chatID)
	return nil
}

func (s *TelegramClientService) ResolveChannelAccessHash(ctx context.Context, userID, channelID int64) (int64, error) {
	client, err := s.ConnectUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("connect user: %w", err)
	}

	mtprotoID := mtprotoChannelID(channelID)
	api := client.API()
	result, err := api.ChannelsGetChannels(ctx, []tg.InputChannelClass{
		&tg.InputChannel{ChannelID: mtprotoID},
	})
	if err != nil {
		return 0, fmt.Errorf("get channels: %w", err)
	}

	chats := result.GetChats()
	for _, chat := range chats {
		if ch, ok := chat.(*tg.Channel); ok && ch.ID == mtprotoID {
			return ch.AccessHash, nil
		}
	}

	return 0, fmt.Errorf("channel %d (mtproto: %d) not found in response", channelID, mtprotoID)
}

func (s *TelegramClientService) ResolvePeerByUsername(ctx context.Context, userID int64, username string) (*tg.InputPeerChannel, error) {
	client, err := s.ConnectUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("connect user: %w", err)
	}

	api := client.API()
	resolved, err := api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("resolve username: %w", err)
	}

	peer, ok := resolved.Peer.(*tg.PeerChannel)
	if !ok {
		return nil, fmt.Errorf("username %s is not a channel", username)
	}

	for _, chat := range resolved.Chats {
		if ch, ok := chat.(*tg.Channel); ok && ch.ID == peer.ChannelID {
			return &tg.InputPeerChannel{
				ChannelID:  peer.ChannelID,
				AccessHash: ch.AccessHash,
			}, nil
		}
	}

	return nil, fmt.Errorf("channel %d not found in resolved chats", peer.ChannelID)
}

func (s *TelegramClientService) CancelFlow(ctx context.Context, userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, ok := s.activeStates[userID]; ok {
		state.cancel()
		delete(s.activeStates, userID)
	}
}

func (s *TelegramClientService) dbStorage(userID int64) telegram.SessionStorage {
	return &dbSessionStorage{
		userID:    userID,
		repo:      s.sessionRepo,
		masterKey: s.masterKey,
	}
}

type dbSessionStorage struct {
	userID    int64
	repo      *repositories.TelegramSessionRepository
	masterKey []byte
}

func (d *dbSessionStorage) LoadSession(ctx context.Context) ([]byte, error) {
	sess, err := d.repo.GetByUserID(ctx, d.userID)
	if err != nil {
		return nil, nil
	}
	if !sess.IsActive || sess.EncryptedSession == "" {
		return nil, nil
	}
	return crypto.Decrypt(sess.EncryptedSession, d.masterKey, d.userID)
}

func (d *dbSessionStorage) StoreSession(ctx context.Context, data []byte) error {
	encrypted, err := crypto.Encrypt(data, d.masterKey, d.userID)
	if err != nil {
		return err
	}

	var phoneHash string
	if existing, err := d.repo.GetByUserID(ctx, d.userID); err == nil {
		phoneHash = existing.EncryptedPhoneHash
	}

	return d.repo.Upsert(ctx, &models.UserTelegramSession{
		UserID:             d.userID,
		EncryptedSession:   encrypted,
		EncryptedPhoneHash: phoneHash,
		IsActive:           true,
	})
}
