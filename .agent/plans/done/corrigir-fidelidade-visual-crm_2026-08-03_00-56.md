# Plano: Corrigir fidelidade visual do CRM administrativo

## Pedido do usuário
Corrigir o CRM administrativo já implementado porque, apesar de ter ficado visualmente agradável, suas cores, proporções e composição ainda não estão parecidas com a dashboard mostrada nas imagens de referência.

## Objetivo
Refazer o acabamento visual do overview em `/admin/dash` como uma reconstrução fiel da interface interna BizLink apresentada nas imagens, preservando os dados, as ações e os módulos administrativos reais do FreddyBot.

O alvo de fidelidade é a janela branca da dashboard — sidebar, topbar, faixa creme de indicadores e quadro de quatro colunas. O fundo promocional escuro e a moldura do navegador vistos nas referências não fazem parte do produto e não serão reproduzidos.

## Contexto atual
- O primeiro redesenho já integrou o CRM ao frontend React/Vite existente e preservou as dez áreas administrativas.
- Busca, ordenação, filtros, métricas e cards estão ligados a dados reais; não é necessário refazer a arquitetura ou o backend.
- A implementação anterior tratou a imagem como inspiração visual, quando o novo pedido exige correspondência mais literal de composição e densidade.
- O tema global é escolhido automaticamente por horário em `useTheme.ts`; entre 18h e 6h o admin abre no tema escuro. Esse comportamento contradiz diretamente a referência, que define uma interface clara com `#1f1f1f`, `#f6f7ed`, `#f4f4f4` e `#ffffff`.
- O hero atual possui três grandes blocos, um gauge de traço grosso e quatro KPIs em grade 2×2. A referência usa quatro zonas horizontais: barras duplas, gauge semicircular de traços finos e dois KPIs amplos.
- A sidebar atual tem 224 px, ícone de marca em bloco preto, subtítulo e grupos mais densos. A referência usa uma sidebar mais estreita, wordmark simples, navegação compacta e separação vertical mais leve.
- Os cards atuais têm pouco conteúdo e aplicam o fundo escuro a todos os usuários em atenção. Na referência o card escuro é um destaque pontual, enquanto os demais alternam branco e cinza suave.
- `General Sans` é indicada explicitamente na imagem, mas seus arquivos não estão presentes no repositório. A pilha atual usa `Geist Variable` como aproximação local.
- O worktree já contém alterações do usuário e trabalhos anteriores. A correção será cirúrgica e não usará reset ou sobrescrita ampla.

## Análise comparativa objetiva

### Referência
- Sidebar ocupa aproximadamente 20% da janela interna; área útil ocupa aproximadamente 80%.
- Topbar e hero formam uma grande faixa creme contínua na área útil, sem uma divisória horizontal forte entre eles.
- Busca fica à esquerda; `Sort by`, `Filters`, perfil textual e botão preto ficam alinhados à direita.
- Gráfico usa pares de barras — uma hachurada e uma sólida — para cada dia.
- Gauge é formado por dezenas de marcas radiais finas, não por um arco espesso.
- Após o gauge aparecem somente dois KPIs grandes, com rótulo pequeno e seta.
- O quadro começa quase imediatamente após o hero, sem um título geral adicional.
- Cabeçalhos de coluna são grandes e leves, com contador compacto à direita.
- Cards têm cerca de 110–120 px de altura, título, descrição em duas linhas, data e metadados; borda fina e raio pequeno.
- Cores declaradas na imagem: `#1f1f1f`, `#f6f7ed`, `#f4f4f4` e `#ffffff`.

### Implementação atual
- A sidebar é mais larga e visualmente mais pesada.
- O topbar é branco e separado do hero creme por padding e bordas, quebrando a faixa contínua da referência.
- Há seletor de tema e perfil com avatar/duas linhas, elementos ausentes na composição original.
- O gráfico mostra apenas uma barra por dia, alternando padrões.
- O gauge é um arco SVG grosso.
- Quatro KPIs em grade comprimem a última zona do hero.
- Existe um cabeçalho extra “Ciclo do usuário / Relacionamento operacional”, que afasta o quadro e muda a hierarquia.
- O destaque escuro está associado ao estágio “Atenção”, criando vários cards escuros e alterando a distribuição visual.

## Arquivos analisados
- `AGENTS.md`
- `.agent/context.md`
- `.agent/memory/memory.md`
- `.agent/plans/done/redesenhar-crm-admin-bizlink_2026-08-03_00-22.md`
- `.agent/skills/fable-method/SKILL.md`
- `.agent/skills/fable-method/references/domains/design-ux.md`
- `original-67e6a497c042e3e5028704313e0738c1.webp`
- `original-5618f2daa44a40756d821e94d38b6c67.webp`
- `original-68641804ff1b3e20d9f1fb6d0a521878.webp`
- `dashboard/src/hooks/useTheme.ts`
- `dashboard/src/index.css`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/CrmOverviewHero.tsx`
- `dashboard/src/components/admin/CustomerPipeline.tsx`
- `dashboard/src/components/admin/CustomerCard.tsx`
- `dashboard/src/components/admin/crmSelectors.ts`

## Arquivos que poderão ser modificados
- `dashboard/src/index.css`
- `dashboard/src/components/AdminDashboard.tsx`
- `dashboard/src/components/admin/AdminLayout.tsx`
- `dashboard/src/components/admin/AdminSidebar.tsx`
- `dashboard/src/components/admin/AdminTopbar.tsx`
- `dashboard/src/components/admin/CrmOverviewHero.tsx`
- `dashboard/src/components/admin/CustomerPipeline.tsx`
- `dashboard/src/components/admin/CustomerCard.tsx`
- `dashboard/src/components/admin/crmSelectors.ts`, somente se a apresentação exigir um seletor derivado adicional sem alterar regras de negócio.
- `dashboard/src/mockData.ts`, somente para demonstrar com clareza os estados visuais existentes.
- `.agent/memory/memory.md`
- `.agent/decisions.md`

## Estratégia de implementação
Usar a referência como uma especificação de layout e reconstruir o shell do overview com dimensões e relações mensuráveis. O admin terá sua paleta clara fixa e isolada dentro de `.admin-layout-v2`, independentemente do tema global da Mini App. Os demais módulos continuarão com os mesmos componentes e comportamentos, recebendo os tokens claros do shell.

A faixa superior será tratada como uma única superfície creme: topbar visualmente contínuo com o hero. O hero passará a ter quatro áreas horizontais nas mesmas proporções da imagem. O pipeline perderá o cabeçalho intermediário e ganhará cartões mais densos, com conteúdo derivado exclusivamente dos campos reais de usuário.

O destaque preto será único e semântico: o usuário selecionado terá prioridade; na ausência de seleção, o primeiro administrador visível poderá representar o card destacado. Blacklist continuará indicada por metadado/estado próprio, sem transformar uma coluna inteira em cards pretos.

Não será criado conteúdo comercial falso, como receita, propostas ou negócios. Os rótulos serão equivalentes do FreddyBot, mas posição, escala, densidade e estilo seguirão a referência.

## Passos detalhados

1. Preservar o worktree atual.
   - Revalidar `git status` e limitar os diffs aos componentes CRM e registros técnicos.
   - Não alterar as correções preexistentes em autenticação, cookies, API ou regras de feature.

2. Fixar o sistema visual claro do admin.
   - Declarar diretamente em `.admin-layout-v2` os quatro tokens da imagem: `#1f1f1f`, `#f6f7ed`, `#f4f4f4` e `#ffffff`.
   - Remover a variante escura do escopo administrativo e impedir que o horário ou o tema Telegram alterem o CRM.
   - Manter o tema global da Mini App intacto fora do admin.
   - Usar bordas quentes finas, raios de 6–9 px e nenhuma sombra decorativa.

3. Corrigir a geometria geral.
   - Reduzir a sidebar desktop para uma proporção próxima da referência.
   - Diminuir paddings do conteúdo e eliminar intervalos que quebram a continuidade do topo.
   - Fazer topbar e hero compartilharem visualmente a superfície creme.
   - Manter drawer no mobile e colapso acessível no desktop, sem sacrificar a composição de referência em 1024×768.

4. Simplificar a sidebar.
   - Usar wordmark textual simples “FreddyBot”, sem bloco preto de ícone e sem subtítulo dominante.
   - Aproximar altura, tamanho de ícone, recuo e item ativo da navegação BizLink.
   - Distribuir módulos reais em grupos visuais equivalentes aos blocos principal, projetos e sistema.
   - Manter todas as dez rotas administrativas acessíveis.

5. Reconstruir o topbar.
   - Preservar busca, ordenação e filtros funcionais na mesma sequência visual da referência.
   - Exibir perfil como controle compacto “Eu”, alinhado em uma linha.
   - Remover o alternador de tema do topbar administrativo.
   - Manter o botão preto ligado à ação real de broadcast, com rótulo curto compatível com o espaço disponível.
   - Garantir equivalentes compactos no mobile.

6. Reconstruir o gráfico de novos usuários.
   - Renderizar dois valores por dia: uma barra hachurada de comparação e uma barra sólida do período atual, ambos derivados de janelas reais de datas.
   - Reproduzir escala, alinhamento inferior, largura e espaçamento observados.
   - Tratar séries vazias e valores iguais sem altura enganosa.

7. Reconstruir o gauge.
   - Substituir o arco grosso por marcas radiais finas em SVG/CSS, preenchendo visualmente a porcentagem real de ativação.
   - Centralizar percentual e rótulo com peso e escala equivalentes à imagem.
   - Manter descrição acessível do valor completo.

8. Reduzir os KPIs do hero para dois blocos.
   - Exibir `Usuários` e `Canais`, métricas reais mais próximas dos dois grandes indicadores da referência.
   - Mover admin/blacklist para o pipeline, filtros ou metadados existentes, sem perder funcionalidade.
   - Posicionar rótulo, número e seta conforme a composição original.

9. Ajustar o pipeline.
   - Remover o título intermediário geral e iniciar diretamente nos quatro cabeçalhos de coluna.
   - Manter as segmentações reais atuais, mas usar títulos curtos e distribuição visual comparável.
   - Aproximar largura, gap, contador e alinhamento dos cabeçalhos da referência.
   - Preservar overflow horizontal em telas menores.

10. Tornar os cards mais fiéis.
    - Adicionar descrição curta derivada de dados reais: identificação Telegram, quantidade/nomes de canais e estado administrativo.
    - Organizar data e metadados em uma linha inferior compacta.
    - Alternar superfícies branca/cinza suave de modo controlado, mantendo borda fina e altura da referência.
    - Aplicar card preto somente ao usuário selecionado ou a um único destaque semântico.
    - Manter clique abrindo o detalhe real e estados de foco/hover acessíveis.

11. Refinar tipografia.
    - Usar pesos, tamanhos e tracking medidos a partir da referência.
    - Manter `Geist Variable` como fallback local seguro.
    - Não baixar nem incorporar General Sans sem ativo/licença confirmados; se o arquivo da fonte for fornecido depois, a troca ficará restrita a `@font-face` e ao token tipográfico.

12. Verificar por comparação visual.
    - Renderizar obrigatoriamente em 1024×768, mesmo tamanho da imagem principal, além de 1440×1024 e 390×844.
    - Comparar sidebar, início e fim do hero, divisões internas, posição dos cabeçalhos e densidade dos cards.
    - Testar busca, filtros, ordenação, broadcast, navegação, card selecionado, vazio e drawer mobile.
    - Registrar diferenças residuais que dependam da fonte General Sans ou de dados reais em maior volume.

13. Executar build e passe `fable-judge`.
    - Rodar Vite build e TypeScript, separando erros novos de erros preexistentes.
    - Confirmar a interface renderizada no navegador e não apenas a compilação.
    - Atualizar memória/decisões e mover o plano para `done` apenas depois das evidências.

## Riscos
- **Fonte exata:** General Sans não está no repositório; sem o arquivo licenciado, Geist continuará sendo uma aproximação.
- **Dados do domínio:** a referência mostra negócios e descrições comerciais; o FreddyBot possui usuários/canais. O conteúdo será semanticamente diferente, embora a estrutura visual seja equivalente.
- **Quantidade de cards:** o mock atual tem poucos registros; a densidade da referência precisa ser validada também com dados de demonstração, sem poluir dados reais.
- **Tema claro fixo:** usuários que preferem tema escuro verão o admin claro. Essa é uma decisão deliberada para cumprir a referência; a Mini App continuará respeitando o tema escolhido.
- **CSS legado:** regras antigas continuam no mesmo arquivo. A correção deve consolidar o bloco CRM, não acrescentar mais camadas conflitantes.
- **Worktree sujo:** há alterações relacionadas e não relacionadas ainda não commitadas; rollback amplo poderia apagar trabalho do usuário.
- **TypeScript preexistente:** o projeto possui erros de tipagem anteriores fora dos arquivos CRM. A conclusão deverá demonstrar ausência de novos erros, sem alegar que toda a base está limpa se não estiver.

## Impactos esperados
- O overview abrirá sempre com as cores exatas da imagem de referência.
- A primeira dobra em 1024×768 terá proporções, hierarquia e densidade muito mais próximas do BizLink.
- O hero deixará de parecer um dashboard genérico e passará a reproduzir barras duplas, gauge de marcas e dois KPIs.
- Todas as ações e configurações administrativas existentes continuarão disponíveis.
- Não haverá alteração de API, banco, autenticação, regras de negócio ou deploy.
- A Mini App de usuários permanecerá visualmente independente.

## Compatibilidade
- Linux: suportado em navegador moderno e no fluxo Vite atual.
- macOS: suportado em Chrome, Firefox e Safari modernos.
- Windows: suportado no frontend; desenvolvimento pode usar WSL/Git Bash conforme o projeto.
- Docker: sem mudança de imagem, dependências ou pipeline.
- CI/CD: preservado; a alteração é de frontend e usa o build existente.
- Mobile/Telegram WebView: shell responsivo preservado, embora a referência de fidelidade seja desktop.

## Como testar

### Build
```bash
cd dashboard && /home/gabriel/.local/share/mise/installs/node/26.5.1/bin/node ./node_modules/vite/bin/vite.js build
```

### Testes
```bash
cd dashboard && /home/gabriel/.local/share/mise/installs/node/26.5.1/bin/node ./node_modules/typescript/bin/tsc --noEmit
```

Validar também manualmente:
- tema global salvo como `dark` e admin ainda renderizando na paleta clara;
- busca por nome, username e ID;
- ordenação e todos os filtros;
- navegação para broadcast e demais módulos;
- card selecionado, blacklist, admin, listas vazias e overflow;
- foco por teclado e contraste.

### Execução
```bash
cd dashboard && /home/gabriel/.local/share/mise/installs/node/26.5.1/bin/node ./node_modules/vite/bin/vite.js --host 127.0.0.1
```

Abrir:
```txt
http://127.0.0.1:5173/admin/dash
```

Capturar e comparar em:
- `1024×768` — referência principal;
- `1440×1024` — desktop amplo;
- `390×844` — mobile.

## Rollback
Reverter apenas os hunks introduzidos por este plano nos componentes CRM, preservando as alterações anteriores do usuário e o primeiro redesenho. O plano não inclui migration, alteração de API ou dados, portanto o rollback é exclusivamente de frontend.

Não usar `git reset --hard`, `git checkout --` amplo nem exclusão de diretórios.

## Observações
- Esta correção assume que “parecido” significa fidelidade da interface interna, não copiar o cenário promocional com fundo rochoso e moldura Safari.
- As configurações reais do admin têm prioridade sobre rótulos fictícios da imagem.
- A implementação somente começará após aprovação explícita deste plano.

## Resultado da implementação
- Admin isolado em tema claro fixo com os quatro tons declarados pela referência.
- Sidebar reduzida para 184 px em 1024×768, com wordmark simples e navegação compacta.
- Topbar e hero formam uma faixa creme contínua.
- Gráfico reconstruído com cinco pares de barras comparando o período atual aos cinco dias anteriores.
- Gauge reconstruído com 45 marcas radiais finas.
- Hero reduzido para dois KPIs reais: usuários e canais.
- Pipeline iniciado imediatamente após o hero, com quatro colunas, cards densos e um único destaque preto semântico.
- Busca, filtros, ordenação, configurações e drawer mobile revalidados no navegador.

## Evidências de verificação
- `vite build`: aprovado, 2023 módulos transformados.
- `git diff --check`: aprovado.
- Renderização observada em 1024×768, 1440×1024 e 390×844.
- Em 1024×768: sidebar `184×768`, topbar `840×70`, hero `840×196`, pipeline iniciado em `y=266`.
- Mesmo com `data-theme="dark"`, o layout computado permaneceu branco e o hero `rgb(246, 247, 237)`.
- Testes de interação: 4 colunas, 45 marcas, 5 pares de barras, 1 card destacado, filtro de blacklist com 1 resultado, busca vazia com 4 estados vazios, Configurações aberta com 41 controles e drawer mobile aberto com overlay.
- `tsc --noEmit`: ainda falha em erros preexistentes espalhados pelo projeto; nenhum erro foi reportado nos componentes CRM criados ou alterados nesta correção.

## Veredito fable-judge
**VERIFIED WITH CAVEATS**

- A correspondência estrutural, as cores, o build e as interações foram reproduzidos.
- General Sans não está disponível no repositório; Geist Variable permanece como fallback tipográfico.
- A validação de navegador usou os dados mock do ambiente Vite; uma sessão autenticada contra a API real não foi necessária para a mudança visual e não foi executada.
