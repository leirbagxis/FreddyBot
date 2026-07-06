package executor

import (
	"context"
	"sync"
)

// Provider e a interface que a factory usa para consultar se um usuario
// possui conta conectada ativa. Isso evita dependencia circular com o
// pacote de servicos.
type Provider interface {
	// HasConnectedAccount retorna true se o usuario possui uma conta
	// conectada ativa e com sessao valida.
	HasConnectedAccount(ctx context.Context, userID int64) bool
}

// ExecutorFactory cria a implementacao correta de TelegramExecutor
// com base na existencia de conta conectada para o usuario.
type ExecutorFactory struct {
	botAPI   TelegramExecutor
	mtproto  *MTProtoExecutor
	provider Provider
	cache    map[int64]TelegramExecutor
	mu       sync.RWMutex
}

// NewExecutorFactory cria uma nova factory.
// botAPI: executor via Bot API (sempre disponivel).
// mtproto: executor via MTProto (opcional, nil se nao configurado).
// provider: interface para consultar se usuario tem conta ativa.
func NewExecutorFactory(
	botAPI TelegramExecutor,
	mtproto *MTProtoExecutor,
	provider Provider,
) *ExecutorFactory {
	return &ExecutorFactory{
		botAPI:   botAPI,
		mtproto:  mtproto,
		provider: provider,
		cache:    make(map[int64]TelegramExecutor),
	}
}

// ForUser retorna o executor apropriado para o usuario.
// Se o usuario tiver uma conta conectada ativa, retorna UserExecutor com MTProto.
// Caso contrario, retorna apenas o BotAPIExecutor.
func (f *ExecutorFactory) ForUser(ctx context.Context, userID int64) TelegramExecutor {
	f.mu.RLock()
	exec, ok := f.cache[userID]
	f.mu.RUnlock()

	if ok {
		return exec
	}

	var chosen TelegramExecutor
	if f.mtproto != nil && f.provider.HasConnectedAccount(ctx, userID) {
		chosen = NewUserExecutor(userID, f.botAPI, f.mtproto)
	} else {
		chosen = f.botAPI
	}

	f.mu.Lock()
	f.cache[userID] = chosen
	f.mu.Unlock()

	return chosen
}

// InvalidateCache limpa o cache de executors para um usuario.
// Deve ser chamado quando uma conta e conectada ou desconectada.
func (f *ExecutorFactory) InvalidateCache(userID int64) {
	f.mu.Lock()
	delete(f.cache, userID)
	f.mu.Unlock()
}

// ForChannel e um atalho que busca o owner do canal e retorna o executor.
func (f *ExecutorFactory) ForChannel(ctx context.Context, channelID int64, getOwnerID func(context.Context, int64) (int64, error)) (TelegramExecutor, error) {
	ownerID, err := getOwnerID(ctx, channelID)
	if err != nil {
		return f.botAPI, err
	}
	return f.ForUser(ctx, ownerID), nil
}
