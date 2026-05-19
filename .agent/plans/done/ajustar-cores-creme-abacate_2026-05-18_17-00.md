# Plano: Ajustar Cores (Fundo Creme e Cards Abacate)

## Pedido do usuário
O usuário quer refinar o esquema de cores: o fundo geral (body) deve ser creme, e os botões e divs (cards/containers) devem ser verde abacate/escuro.

## Objetivo técnico
Reestruturar a hierarquia de cores da dashboard (tema Light) para que o fundo seja claro (creme) e os elementos de superfície (cards, botões) sejam escuros (verde abacate).

## Estratégia de implementação
Como os `cards` e `divs` agora serão escuros (verde abacate) e o fundo será claro (creme), precisamos garantir que o texto dentro desses cards seja legível (branco/claro). 
Atualmente, a variável `--text` é global. Se o fundo é claro e o card é escuro, precisamos garantir que o texto seja claro sobre o fundo escuro.

**Proposta de Cores:**
- `--bg` (Fundo): `#F5EFE6` (Creme)
- `--card` e `--surface` (Divs principais): `#3D5229` (Verde Abacate Escuro)
- `--accent` (Botões/Destaques): `#2B3D18` ou um abacate mais vibrante `#6b8e23`.
- O texto global (`--text`) será claro (`#FFFFFF`), pois quase todo o conteúdo da dashboard está dentro de cards. Para os títulos que ficarem soltos no fundo creme (se houver), farei ajustes específicos se necessário.

## Passos detalhados
1.  **Modificar `dashboard/src/index.css` (Variáveis):**
    - Atualizar `[data-theme="light"]`:
        - `--bg`: `#F5EFE6`
        - `--card`, `--card-elevated`, `--surface`: `#3D5229` (Abacate Escuro)
        - `--text`: `#F0F4E8` (Texto claro)
        - `--text-secondary`, `--hint`: `#C9D6B8`
        - `--accent`: `#6B8E23` (Abacate vibrante)
        - `--border`: `#2B3D18`
2.  **Atualizar Telegram SDK:**
    - Ajustar `useTheme.ts` para que o background seja o creme e o cabeçalho seja o creme também (ou o verde abacate, dependendo de como ficar melhor visualmente).
3.  **Aprovação e Build:**
    - Pedir aprovação, implementar e rodar `npm run build`.

## Riscos
A inversão (fundo claro com cards muito escuros) essencialmente transforma o layout num híbrido de dark mode com bordas claras. Precisarei garantir que textos que ficam fora dos cards (se existirem) não sumam no fundo creme.