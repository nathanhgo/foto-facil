# Backend Go

Convenções do backend (nós, DAG, storage). Escopo: `backend/**/*.go`.

- Todo nó de processamento implementa a interface `Node` (`backend/internal/nodes/node.go`): `GetID() string` e `Process(ctx *ProcessContext) error`.
- `ProcessContext.NodeOutputs` mapeia ID do nó para suas imagens de saída — cada nó deve ler estritamente as saídas de seus **pais diretos**, nunca um buffer global, para não vazar dados entre ramificações paralelas do grafo.
- Novo nó = novo arquivo (`backend/internal/nodes/<nome>.go`) + teste próprio (`<nome>_test.go`) validando entrada/saída em nível de pixel.
- Registre o nó novo no `switch` de instanciação em `backend/internal/api/websocket.go` (mapeia `OriginalType` string → struct Go).
- Cuidado com condições de corrida ao processar lotes de imagens (`[]image.Image`) — hoje a execução é sequencial, mas escreva código que não assuma isso será sempre verdade.
- Evite paths/portas hardcoded quando fizer sentido usar `.env` (mesmo que hoje `.env` ainda não esteja conectado a `backend/cmd/main.go` — ver `architecture-docs/stack.md`); se notar hardcode novo sendo introduzido, sinalize.
- SQLite: runtime usa `modernc.org/sqlite` (sem CGO); testes usam `github.com/mattn/go-sqlite3` (com CGO) — ambos devem continuar funcionando após mudanças em `backend/internal/storage/storage.go`.
- **Protocolo do WebSocket** (`backend/internal/api/websocket.go`): toda resposta de `RUN_FLOW` ecoa o `requestId` recebido no pedido. É o que permite ao frontend saber que a resposta pertence à execução que ele pediu — não remova o eco nem invente um id quando o pedido não trouxe nenhum (o autosave manda `RUN_FLOW` sem id e a resposta dele não pode parecer a de um pedido do botão). Nó de tipo desconhecido (fora do `switch` de instanciação) agora **falha a execução** com erro: antes era ignorado em silêncio e a execução terminava como sucesso sem ter processado o nó.