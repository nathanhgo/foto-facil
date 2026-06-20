# Arquitetura Geral

## Visão Geral
O **Foto Fácil** é uma aplicação desktop de processamento digital de imagens baseada em grafos (DAG - Directed Acyclic Graph). O sistema é dividido entre um frontend reativo (Electron) e um backend de alto desempenho (Golang) encarregado do processamento bruto das imagens usando OpenCV.

## Tecnologias
- **Frontend:** Electron + (Framework React)
- **Backend:** Golang + biblioteca de manipulação de imagem (ex: `gocv` que é wrapper do OpenCV)
- **Comunicação:** WebSockets ou IPC para comunicação rápida e atualizações de progresso em tempo real.
- **Armazenamento:** JSON para serialização/deserialização dos fluxos.

## Fluxo de Execução
1. **Interface do Usuário (Frontend):** O usuário interage com a tela montando nós. Cada nó tem um ID e propriedades.
2. **Serialização:** Quando o usuário clica em "Rodar" ou "Preview", o grafo é convertido para um JSON descrevendo as dependências (qual nó consome a saída de qual nó).
3. **Processamento (Backend - Golang):**
   - O backend recebe o JSON.
   - Um **Scheduler** analisa o DAG e identifica as dependências.
   - Nós independentes são disparados em Goroutines separadas (concorrência).
   - Para nós que processam lotes (várias imagens), as imagens são divididas em workers (paralelismo).
4. **Feedback de Progresso:** O backend emite eventos (ex: `node_started`, `node_progress`, `node_finished`) de volta para o frontend atualizar a UI.
5. **Resultado:** As imagens resultantes são salvas no disco ou enviadas em base64/buffer para visualização nos nós intermediários/finais.

## Tratamento de Múltiplas Imagens
- Quando um nó como `Directory Batch` injeta um array de imagens, a "aresta" do grafo passa a transportar um tipo `[]Image`.
- Nós subsequentes devem ser capazes de iterar sobre o input `[]Image` aplicando de forma unificada seus processamentos em todo o lote de uma vez, otimizando as rotinas com concorrência.
