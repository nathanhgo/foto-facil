---
name: sync-project-docs
description: >-
  Revisa o código do FotoFácil e compara com architecture-docs/stack.md e
  architecture-docs/mvp.md, corrigindo divergências entre o que os documentos
  descrevem e o que está de fato implementado. Use quando o usuário pedir para
  sincronizar, auditar ou atualizar a documentação de arquitetura, ou quando
  suspeitar que os docs estão desatualizados em relação ao código.
disable-model-invocation: true
---

# Sincronizar documentação com o código real

Skill sob demanda para auditar `architecture-docs/` contra o estado real do código, evitando que a IA (ou você) confie em informação desatualizada.

## Passo a passo

1. **Levantar o estado real**:
   - Liste os nós implementados em `backend/internal/nodes/*.go` (ignorando `*_test.go`) e o que cada um realmente faz.
   - Confira `backend/go.mod` e `frontend/package.json` para versões de dependências.
   - Confira `backend/internal/api/websocket.go` para as ações WebSocket suportadas.
   - Rode `go test ./...` para confirmar quais testes existem e passam.

2. **Comparar com os documentos**:
   - `architecture-docs/stack.md`: as tecnologias listadas (e a seção "Aspiracional vs Real") ainda batem com o código?
   - `architecture-docs/mvp.md`: os checkboxes `[x]`/`[ ]`/`[~]` refletem o estado real? Novo nó implementado sem estar marcado? Nó marcado como feito mas removido/alterado?
   - `architecture-docs/idea.md`: alguma seção descreve uma decisão que já mudou?

3. **Corrigir divergências**: edite os arquivos diretamente, mantendo o formato e o idioma (pt-br) já usados. Preserve a estrutura de seções existente — não reescreva do zero.

4. **Registrar**: adicione uma entrada em `architecture-docs/logs.md` resumindo o que foi corrigido e por quê.

## Exemplo de divergência real já encontrada

O `mvp.md` marcava "IA/Machine Learning" como `[x]` citando ONNX, mas `backend/internal/nodes/ai.go` implementa apenas heurísticas (distância cromática de cantos + sigmoide) — sem modelo real carregado. Esse tipo de gap é exatamente o que esta skill deve pegar.
