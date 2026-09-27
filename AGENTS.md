# FotoFácil

Editor de imagens baseado em nós (DAG) em Electron/React/TypeScript (frontend) + Go (backend), com foco em processamento digital de imagens (PDI) para uso leigo e acadêmico. O projeto será apresentado na WORCAP (INPE).

## Quando ler qual arquivo

| Se você vai... | Leia primeiro |
| --- | --- |
| Entender o produto, público-alvo e requisitos | `architecture-docs/idea.md` |
| Saber a stack e o estado **real** de implementação | `architecture-docs/stack.md` |
| Ver fases, features e status de cada nó | `architecture-docs/mvp.md` |
| Saber como testar / o que já tem cobertura | `architecture-docs/tests.md` |
| Ver decisões em aberto e dúvidas do projeto | `architecture-docs/questions.md` |
| Trabalhar em UI, tema ou design | `architecture-docs/visual.md` |
| Ver o histórico do que já foi feito (log de sessão) | `architecture-docs/logs.md` |
| Adicionar/alterar nós de processamento | `architecture-docs/nodes-reference.md` |
| Seguir o workflow de TDD | `.agents/rules/10-tdd-workflow.md` |
| Saber que log/entrada adicionar após editar código | `.agents/rules/20-logging.md` |
| Padrão de mensagens de commit | `.agents/rules/90-commits.md` |
| Ver convenções do backend Go (nós, DAG, storage) | `backend/AGENTS.md` |
| Ver convenções do frontend React/Electron | `frontend/AGENTS.md` |

> Não confie em `general_architecture.md`/`idea.md` para o **estado atual** de implementação — sempre confira `stack.md`.

## Idioma

- Código-fonte e comentários: **inglês**.
- Documentos em `architecture-docs/` (incluindo `logs.md`): **pt-br**.
- Textos visíveis na UI: pt-br (padrão já usado no projeto).

## TDD (obrigatório)

Toda feature nova segue red-green-refactor: escreva o teste primeiro e confirme que ele falha, implemente o mínimo para passar, então refatore mantendo verde. Rode a suíte relevante (backend: `go test ./...` em `backend/`) antes de concluir e reporte o resultado.
Detalhes completos em `.agents/rules/10-tdd-workflow.md`.

## Convenções específicas

- **Backend Go** (nós, DAG, storage, WebSocket): ver `backend/AGENTS.md`.
- **Frontend React/Electron** (nós na UI, tema, IPC): ver `frontend/AGENTS.md`.

## Regras gerais

- Evite cores e valores hardcoded na UI — use o sistema de tema existente (roxo como accent color, dark/light mode).
- Este é um projeto pessoal com apresentação acadêmica agendada (WORCAP); ao adicionar nós novos, priorize os listados na Fase 2 de `mvp.md` (PDI acadêmico) sobre features "consumer".
- Nenhuma chave de API paga é necessária — processamento é 100% local/offline.