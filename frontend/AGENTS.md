# Frontend React/Electron

Convenções do frontend. Escopo: `frontend/**/*.ts`, `frontend/**/*.tsx`.

- Novo nó de backend precisa de contraparte no frontend: entrada em `LIBRARY_NODES` (barra lateral, definida em `frontend/src/core/nodeCatalog.ts`) + suporte no `PropertiesPanel` (inputs específicos: sliders, selects, etc.) em `frontend/src/App.tsx`.
- `frontend/src/App.tsx` é hoje um arquivo monolítico grande. Ao adicionar funcionalidade não trivial, prefira extrair um componente novo em arquivo próprio (por exemplo sob `frontend/src/components/`) em vez de crescer ainda mais o arquivo único.
- Não hardcode cores (ex.: `#9b51e0`) — use o objeto de tema existente para suportar dark/light mode.
- Comunicação com o backend é via WebSocket, com ações `RUN_FLOW`, `SAVE_FLOW`, `LOAD_FLOWS`, `DELETE_FLOW` — não introduza um novo canal de comunicação sem necessidade. A URL é montada por `getBackendWebSocketUrl()` (`frontend/src/core/websocketUrl.ts`), que usa `VITE_BACKEND_WS_PORT` (do `.env` na raiz) com fallback para a porta de dev padrão.
- O status do botão "Rodar Fluxo" só muda com resposta correlacionada do backend (`requestId`), através de `frontend/src/core/runFlowStatus.ts`: `done` exige a resposta daquele pedido específico, erro do backend vira `error` com a mensagem real, e resposta de preview/execução anterior atualiza apenas thumbnails. **Nunca** marque `done`/`running` por timer (`setTimeout`) — foi exatamente esse o bug que fez a UI mostrar "✓ Concluído" sem ter executado nada.
- Diálogos nativos (abrir arquivo/pasta) passam pelo IPC do Electron (`frontend/electron/main.ts` + `preload.ts` com `contextBridge`), nunca acesse `fs` diretamente do processo de renderização.