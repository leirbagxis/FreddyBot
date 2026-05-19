# Plano: User Profile Card Minimalista

## Pedido do usuário
O usuário achou o "Painel de Telemetria" muito exagerado/feio e pediu algo mais "básico", "account do usuário", uma "identidade" mais simples e elegante.

## Objetivo técnico
Substituir o complexo componente `.telemetry-hud` por um perfil de usuário minimalista e limpo, aproveitando a classe `.card` padrão (que já possui o efeito liquid-glass e cantos arredondados), com foco apenas nas informações essenciais: Avatar (inicial), Nome, ID e Link de acesso.

## Estratégia de implementação

**1. Design Minimalista (Profile Card):**
- Usar a classe `.card` padrão que já tem o efeito glass.
- Estrutura:
  - **Cabeçalho/Perfil:** Um círculo grande com a inicial do usuário (usando a cor Laranja McLaren de fundo), ao lado do Nome (grande) e o User ID (pequeno, cinza).
  - **Rodapé/Acesso:** Uma caixa simples com fundo escurecido/translúcido mostrando o link de convite, com um ícone sutil para abrir.
- Nenhuma animação exagerada, nenhum jargão técnico (como "Uplink" ou "Chief Engineer"). Apenas "Sua Conta" ou "Perfil".

**2. Alterações no React (`DashboardInicioTab.tsx`):**
Remover todo o bloco `.telemetry-hud` e substituir por:
```tsx
<div className="card animate-stagger-in">
    <div className="flex items-center gap-4 mb-5">
        <div className="w-12 h-12 rounded-full flex items-center justify-center bg-[var(--accent)] text-white font-bold text-xl flex-shrink-0">
            {displayName.charAt(0).toUpperCase()}
        </div>
        <div className="min-w-0">
            <h3 className="text-lg font-bold text-[var(--text)] truncate">{displayName}</h3>
            <p className="text-xs text-[var(--hint)] truncate">ID: {channel.ownerId}</p>
        </div>
    </div>
    <div className="bg-[var(--surface)] border border-[var(--border)] rounded-xl p-3 flex items-center justify-between cursor-pointer hover:border-[var(--accent)] transition-colors"
         onClick={() => { ...open link... }}>
        <div className="min-w-0 flex-1">
            <p className="text-[10px] font-semibold text-[var(--hint)] uppercase tracking-wider mb-1">Link de Acesso</p>
            <p className="text-sm font-mono text-[var(--text)] truncate">{channel.inviteUrl.replace('https://', '')}</p>
        </div>
        <ExternalLink size={16} className="text-[var(--accent)] ml-3 flex-shrink-0" />
    </div>
</div>
```

**3. Alterações no CSS (`index.css`):**
- Remover todo o bloco `/* ===== TELEMETRY HUD (McLAREN) ===== */` (classes `.telemetry-hud`, `.telemetry-pulse`, etc.).
- O novo design usará principalmente classes utilitárias do Tailwind e as variáveis globais que já definimos, mantendo o arquivo CSS extremamente limpo.

## Riscos
Nenhum risco funcional. A simplificação reduz a complexidade do código e melhora a manutenibilidade, agradando o gosto do usuário por designs mais limpos e diretos.