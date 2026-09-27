# Backend Go

Convenções do backend (nós, DAG, storage). Escopo: `backend/**/*.go`.

- Todo nó de processamento implementa a interface `Node` (`backend/internal/nodes/node.go`): `GetID() string` e `Process(ctx *ProcessContext) error`.
- `ProcessContext.NodeOutputs` mapeia ID do nó para suas imagens de saída — cada nó deve ler estritamente as saídas de seus **pais diretos**, nunca um buffer global, para não vazar dados entre ramificações paralelas do grafo.
- Novo nó = novo arquivo (`backend/internal/nodes/<nome>.go`) + teste próprio (`<nome>_test.go`) validando entrada/saída em nível de pixel.
- Registre o nó novo no `switch` de instanciação em `backend/internal/api/websocket.go` (mapeia `OriginalType` string → struct Go).
- Cuidado com condições de corrida ao processar lotes de imagens (`[]image.Image`) — hoje a execução é sequencial, mas escreva código que não assuma isso será sempre verdade.
- Evite paths/portas hardcoded quando fizer sentido usar `.env` (mesmo que hoje `.env` ainda não esteja conectado a `backend/cmd/main.go` — ver `architecture-docs/stack.md`); se notar hardcode novo sendo introduzido, sinalize.
- SQLite: runtime usa `modernc.org/sqlite` (sem CGO); testes usam `github.com/mattn/go-sqlite3` (com CGO) — ambos devem continuar funcionando após mudanças em `backend/internal/storage/storage.go`.
- Protocolo `RUN_FLOW` correlaciona pedido/resposta via `requestId`: o payload aceita `requestId` (`ReactFlowPayload.RequestID`) e `RunFlowResponse` o ecoa junto de `action: "RUN_FLOW"` em **todos** os caminhos (sucesso, erro de nó e ciclo no grafo). Nunca remova esse eco: o frontend só conclui a rodada do usuário quando o `requestId` da resposta bate com o pedido pendente. Preview (autosave) pode omitir `requestId`. Contrato completo em `architecture-docs/websocket-protocol.md`.