package executor

import (
	"context"
	"fmt"
	"sync"
)

// Provider e a interface que a factory usa para consultar se um usuario
// possui conta conectada ativa ou assinatura premium. Isso evita
// dependencia circular com o pacote de servicos.
type Provider interface {
	// HasConnectedAccount retorna true se o usuario possui uma conta
	// conectada ativa e com sessao valida.
	HasConnectedAccount(ctx context.Context, userID int64) bool

	// HasPremiumManagedAccount retorna true se o usuario possui uma
	// assinatura premium ativa com o recurso ManagedPremiumAccount habilitado.
	HasPremiumManagedAccount(ctx context.Context, userID int64) bool
}

// ExecutorFactory cria a implementacao correta de TelegramExecutor
// com base na existencia de conta conectada ou premium para o usuario.
//
// A ordem de prioridade e:
//  1. ConnectedAccount (MTProto proprio do usuario) -> UserExecutor
//  2. Premium (conta admin gerenciada) -> PremiumExecutor
//  3. Padrao -> BotAPIExecutor
type ExecutorFactory struct {
	botAPI   TelegramExecutor
	mtproto  *MTProtoExecutor
	adminSP  AdminSessionProvider
	provider Provider
	cache    map[int64]TelegramExecutor
	mu       sync.RWMutex
}

// NewExecutorFactory cria uma nova factory.
// botAPI: executor via Bot API (sempre disponivel).
// mtproto: executor via MTProto (opcional, nil se nao configurado).
// adminSP: provedor de sessao admin para PremiumExecutor (opcional, nil se nao configurado).
// provider: interface para consultar se usuario tem conta ativa ou premium.
func NewExecutorFactory(
	botAPI TelegramExecutor,
	mtproto *MTProtoExecutor,
	adminSP AdminSessionProvider,
	provider Provider,
) *ExecutorFactory {
	return &ExecutorFactory{
		botAPI:   botAPI,
		mtproto:  mtproto,
		adminSP:  adminSP,
		provider: provider,
		cache:    make(map[int64]TelegramExecutor),
	}
}

// ForUser retorna o executor apropriado para o usuario.
// A ordem de prioridade:
//  1. Se o usuario tiver uma conta conectada ativa -> UserExecutor com MTProto.
//  2. Se o usuario tiver premium (ManagedPremiumAccount) -> PremiumExecutor com admin MTProto.
//  3. Caso contrario -> BotAPIExecutor.
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
	} else if f.mtproto != nil && f.adminSP != nil && f.provider.HasPremiumManagedAccount(ctx, userID) {
		chosen = NewPremiumExecutor(userID, f.botAPI, f.mtproto, f.adminSP)
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
		return nil, fmt.Errorf("get owner for channel %d: %w", channelID, err)
	}
	return f.ForUser(ctx, ownerID), nil
}
