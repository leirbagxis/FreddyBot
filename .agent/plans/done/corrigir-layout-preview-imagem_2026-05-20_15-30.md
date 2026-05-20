# Plano: Corrigir Recorte da Imagem no Preview

## Pedido do usuário
A imagem no preview do dashboard está saindo "pela metade" (cortada).

## Objetivo técnico
Ajustar o estilo CSS da imagem de prévia no componente `AdminNoticeTab.tsx` para garantir que a imagem inteira seja visível, sem cortes agressivos, mantendo a proporção original.

## Contexto atual
O componente `AdminNoticeTab.tsx` utiliza a classe `object-cover` e um `max-h-[200px]`. A propriedade `object-cover` redimensiona a imagem para preencher todo o espaço disponível, cortando as bordas se a proporção da imagem for diferente do container (320px de largura por 200px de altura).

## Arquivos analisados
- `dashboard/src/components/AdminNoticeTab.tsx`

## Arquivos que poderão ser modificados
- `dashboard/src/components/AdminNoticeTab.tsx`

## Estratégia de implementação
1. Alterar `object-cover` para `object-contain`. Isso garantirá que a imagem caiba inteira dentro do container.
2. Adicionar `w-auto` e `mx-auto` se necessário, ou manter `w-full` com `object-contain`.
3. Alternativamente, aumentar o `max-h` ou permitir que a altura seja automática até um limite maior (ex: `max-h-[350px]`).

## Passos detalhados

1. **Modificar `dashboard/src/components/AdminNoticeTab.tsx`**:
    - Localizar a tag `<img>` dentro de `renderPreview`.
    - Substituir `object-cover` por `object-contain`.
    - Aumentar `max-h-[200px]` para `max-h-[350px]` para dar mais espaço a imagens verticais.

## Riscos
Nenhum. É apenas uma mudança estética no frontend.

## Impactos esperados
- A imagem de prévia será exibida por completo, permitindo ao admin verificar o conteúdo total da mídia antes do disparo.

## Como testar

### Build
```bash
cd dashboard && npm run build
```

### Testes
1. Ir na dashboard admin -> Aba Broadcast.
2. Colar uma URL de imagem ou File ID (especialmente uma imagem vertical ou muito larga).
3. Verificar se a imagem aparece inteira no preview, sem partes cortadas.

## Rollback
Reverter a classe CSS para `object-cover` e `max-h-[200px]`.
