# Plano: Refatorar Fundo "Minha Conta" e Criar Card de Navegação em Div com Divisórias

## Pedido do usuário
1. **Fundo da Div "Minha Conta"**: Mudar a cor de fundo interna para o mesmo tom de branco leve das bordas (`bg-white/10` ou `bg-white/5`).
2. **Nova Div de Navegação**:
   - Criar uma div estilizada abaixo da "Minha Conta", com o mesmo fundo branco leve e bordas suavemente arredondadas (`rounded-2xl bg-white/5 border border-white/15`).
   - Dentro da div incluir 4 opções:
     - **Meus Canais**
     - **Meus Templates**
     - **Meus Agendamentos**
     - **Suporte**
   - Ícones monocromáticos (sem cor personalizada) à esquerda do nome correspondente a cada opção.
   - Ícone de seta `>` (`ChevronRight`) proporcional à direita de cada opção.
   - Linha divisória fina (`border-b border-white/10`) separando cada uma das opções.

## Objetivo
Visual unificado no `SideMenu.tsx`: card "Minha Conta" com fundo branco translúcido leve e novo container unificado de navegação com linhas divisórias entre os itens, ícones limpos sem cores fortes e setas `>` à direita.

## Contexto atual
- `SideMenu.tsx` possui `bg-muted/30` na div "Minha Conta" e botões soltos na lista de navegação.

## Arquivos analisados
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/App.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/SideMenu.tsx`
- `dashboard/src/App.tsx`

## Estratégia de implementação

1. **`dashboard/src/components/SideMenu.tsx`**:
   - Alterar fundo de "Minha Conta" de `bg-muted/30` para `bg-white/5 border border-white/15`.
   - Encapsular os 4 itens de navegação dentro de um container com `bg-white/5 border border-white/15 rounded-2xl overflow-hidden`.
   - Itens de navegação:
     - **Meus Canais**: Ícone `<Hash size={16} />`, texto "Meus Canais", ícone `<ChevronRight size={14} />`.
     - **Meus Templates**: Ícone `<FileText size={16} />`, texto "Meus Templates", ícone `<ChevronRight size={14} />`.
     - **Meus Agendamentos**: Ícone `<Calendar size={16} />`, texto "Meus Agendamentos", ícone `<ChevronRight size={14} />`.
     - **Suporte**: Ícone `<Headphones size={16} />`, texto "Suporte", ícone `<ChevronRight size={14} />`.
   - Remover classes de cores vibrantes nos ícones (usar a cor padrão monocromática do texto `text-foreground` / `text-muted-foreground`).
   - Adicionar borda inferior separadora `border-b border-white/10` em cada item, exceto no último.

2. **Validação**:
   - Executar `npm run build` na pasta `dashboard`.

## Passos detalhados
1. Editar `dashboard/src/components/SideMenu.tsx`.
2. Executar `npm run build`.

## Riscos
- **Baixo**: Alteração visual de UI no menu lateral.

## Impactos esperados
- Menu elegante com cartões em tom translúcido unificado (branco leve), lista de navegação limpa, ícones neutros e divisórias sutis.

## Compatibilidade
- Linux, macOS, Windows, Docker, CI/CD

## Como testar
```bash
cd dashboard && npm run build
```

## Rollback
```bash
git checkout dashboard/src/components/SideMenu.tsx
```
