---
name: add-processing-node
description: >-
  Guia o fluxo completo para adicionar um novo nó de processamento de imagem
  ao FotoFácil, do teste TDD no backend Go até o registro na UI do frontend e
  a atualização da documentação. Use quando o usuário pedir para criar,
  implementar ou adicionar um novo nó/feature de processamento de imagem
  (ex.: histograma, convolução, detecção de bordas, blur, threshold).
---

# Adicionar um novo nó de processamento

Fluxo completo para implementar um nó novo do zero, seguindo TDD e mantendo a documentação sincronizada.

## Passo a passo

1. **Especificar o nó**: confirme com o usuário (ou infira de `architecture-docs/mvp.md`) entrada, saída e parâmetros configuráveis do nó. Verifique se já existe uma entrada correspondente em `architecture-docs/mvp.md` (Fase 2) — se não existir, é uma boa oportunidade para adicionar.

2. **Teste primeiro (TDD)**: crie `backend/internal/nodes/<nome>_test.go` com uma imagem de entrada conhecida e o resultado esperado em nível de pixel. Rode `go test ./...` a partir de `backend/` e confirme que falha (arquivo de implementação ainda não existe/está incompleto).

3. **Implementar o nó**: crie `backend/internal/nodes/<nome>.go` com uma struct que implementa a interface `Node` (`GetID() string`, `Process(ctx *ProcessContext) error`, ver `backend/internal/nodes/node.go`). Leia apenas as imagens vindas dos pais diretos via `ctx.NodeOutputs`. Rode os testes de novo até passar.

4. **Registrar no backend**: adicione o novo tipo no `switch` de instanciação de nós em `backend/internal/api/websocket.go` (mapeia a string `OriginalType` vinda do frontend para a struct Go).

5. **Registrar no frontend**:
   - Em `frontend/src/core/nodeCatalog.ts`, adicione uma entrada em `NODE_CATALOG` com `label`, `categoryId` (reaproveite uma categoria existente em `CATEGORIES` sempre que fizer sentido) e `docs` (`technical` e `effect`, em pt-br — usados no painel de propriedades e em `nodes-reference.md`).
   - Em `frontend/src/App.tsx`, adicione os inputs correspondentes no `PropertiesPanel` (sliders/selects/inputs para os parâmetros do nó). `LIBRARY_NODES` e as cores já são derivados automaticamente do catálogo, não precisa duplicar.

6. **Atualizar documentação**:
   - Marque o item correspondente como `[x]` em `architecture-docs/mvp.md` e adicione a referência ao arquivo-fonte.
   - Adicione a entrada correspondente (técnico + efeito) em `architecture-docs/nodes-reference.md`, na categoria correta, espelhando o que foi colocado em `nodeCatalog.ts`.
   - Adicione uma entrada em `architecture-docs/logs.md` descrevendo o que foi feito (ver `.cursor/rules/20-logging.mdc`).
   - Se o nó for relevante para a demo da WORCAP, considere atualizar `architecture-docs/worcap.md`.

7. **Validar**: rode `go test ./...` no backend uma última vez e teste manualmente o fluxo no Electron (`npm run dev` no frontend + `go run cmd/main.go` no backend) antes de considerar concluído.

## Notas

- Prefira nomes de arquivo e struct consistentes com os nós existentes (ex.: `crop_resize.go` → `CropResizeNode`).
- Se o nó precisar de um parâmetro que se aplica a todo o lote (ex.: ângulo de rotação), o valor deve ser definido uma vez no painel de propriedades e aplicado a todas as imagens do lote — não peça o parâmetro por imagem.
- Código e comentários em inglês; textos de UI em pt-br, seguindo o padrão já usado no projeto.
