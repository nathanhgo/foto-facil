# TDD é obrigatório

Toda feature nova (nó, endpoint, lógica de negócio) segue red-green-refactor:

1. Escreva o teste primeiro (`*_test.go` ao lado do arquivo que será criado/alterado) e confirme que ele falha.
2. Implemente o mínimo necessário para o teste passar.
3. Refatore mantendo os testes verdes.

## Backend (Go)

- Rode `go test ./...` (a partir de `backend/`) antes de considerar qualquer tarefa concluída.
- Novos nós de processamento: teste deve validar a matriz de pixels de saída contra valores esperados, não apenas "não retornou erro".
- Use testes table-driven quando houver múltiplos casos (ex.: diferentes ângulos de rotação, diferentes kernels de convolução).

## Frontend (TypeScript/React)

- A suíte é o Vitest (`npm test` em `frontend/`); a cobertura ainda é pequena e não há testes de componentes React.
- Se a tarefa envolver lógica testável isolável (parsing, validação, transformação de dados do grafo), escreva o teste antes. Se faltar infraestrutura para o que precisa ser testado, sinalize a lacuna e proponha o que configurar antes de prosseguir, em vez de simplesmente pular a etapa de teste.
- Estado atual dos testes e da cobertura: `architecture-docs/stack.md`.

## Antes de finalizar

Sempre rode a suíte de testes relevante e reporte o resultado. Não marque uma tarefa como concluída com testes quebrando.
