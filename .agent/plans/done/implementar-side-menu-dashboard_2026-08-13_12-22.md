# Plano: Implementar Menu Lateral (Side Menu / Drawer) na Dashboard

## Pedido do usuário
Criar um menu lateral (side menu/drawer) acionado por um ícone no header da dashboard com as seguintes características:
- **Topo**: Foto de perfil do usuário em formato circular (`rounded-full`), ao lado do nome do usuário e, abaixo do nome, o ID do usuário (alinhados verticalmente ao centro ao lado da foto).
- **Bloco "Minha Conta" (My Account)**: Um container com cantos arredondados (`rounded-xl`/`2xl`) contendo:
  - Título "Minha Conta" à esquerda.
  - Exibição do número de canais que o usuário possui.
- **Navegação Rápida**: Ícones e links para Meus Canais, Meus Templates, Conta Telegram e Suporte/Ouvidoria.

## Objetivo
Criar o componente `SideMenu.tsx` (drawer deslizante com overlay) e integrá-lo ao header da Dashboard (`App.tsx`), garantindo um design fluido, transparente e responsivo.

## Contexto atual
- O header (`top-bar`) atualmente exibe o botão de voltar e o alternador de tema.
- Os dados do usuário (`displayName`, `photo_url`, `id`, `channels`) já estão disponíveis no estado do `App.tsx`.

## Arquivos analisados
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`
- `dashboard/src/components/ui/dialog.tsx` (ou Sheet para drawers)

## Arquivos que poderão ser criados/modificados
- `dashboard/src/components/SideMenu.tsx` (NOVO)
- `dashboard/src/App.tsx`
- `dashboard/src/index.css`

## Estratégia de implementação

1. **Criar `dashboard/src/components/SideMenu.tsx`**:
   - Overlay translúcido com fechamento ao clicar fora ou na tecla ESC.
   - Painel lateral deslizante (slide-in) com bordas arredondadas e efeito de transparência/blur.
   - **Cabeçalho**:
     - Avatar circular (`rounded-full` 48x48px com a foto ou inicial).
     - Nome do usuário em negrito + ID em fonte menor (`text-muted-foreground`) ao lado da foto.
   - **Card "Minha Conta"**:
     - Container `rounded-2xl` com fundo sutil (`bg-muted/40 border border-border p-4`).
     - Ícone `User` + Texto "Minha Conta" + Contador de canais (ex: `3 canais`).
   - **Lista de Links e Atalhos**:
     - 📢 Meus Canais (`/me/channels`)
     - 📄 Meus Templates (abre modal de templates)
     - 👤 Conta Telegram (abre modal de conexão)
     - 🎧 Suporte / Ouvidoria (link externo `t.me`)

2. **Integrar no `dashboard/src/App.tsx`**:
   - Adicionar estado `showSideMenu` (boolean).
   - Adicionar o botão de menu (`<Menu size={22} />`) no header `.top-bar`.
   - Renderizar o `<SideMenu />` repassando os dados do usuário (`tgUser`, `displayName`, `userID`, `channelsCount`, etc.).

3. **Estilos em `dashboard/src/index.css`**:
   - Animações de entrada/saída do drawer (`slide-in-left`, `fade-in`).

4. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Criar `dashboard/src/components/SideMenu.tsx`.
2. Atualizar `dashboard/src/App.tsx` com o gatilho e estado do menu.
3. Adicionar regras no `dashboard/src/index.css` se necessário.
4. Executar `npm run build`.

## Riscos
- **Baixo**: Componente novo e desacoplado, sem breaking changes.

## Impactos esperados
- Menu lateral intuitivo, elegante e com todas as informações da conta do usuário acessíveis com 1 clique.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
rm -f dashboard/src/components/SideMenu.tsx
git checkout dashboard/src/App.tsx
git checkout dashboard/src/index.css
```
