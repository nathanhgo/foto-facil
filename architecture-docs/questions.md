# Questões e Decisões

## Questões em Aberto
1. **Engine do Frontend:** Qual framework utilizaremos em conjunto com o Electron? (Sugestões: React.js para ter boa reatividade, ou Vue.js para facilidade de integração).
Resposta: Vamos usar React.JS, com typescript se for possível

2. **Comunicação Front/Back:** A comunicação entre o Electron (Frontend) e o Golang (Backend) será feita via gRPC, WebSockets ou IPC padrão com requisições HTTP locais? (WebSocket é bom para atualizações em tempo real de progresso de processamento em lote).
Resposta: Vamos usar majoritariamente websockets, em situações mais simples podemos usar o HTTP, mas priorize o websocket

3. **Gerenciamento de Estado do DAG:** Onde o estado principal do grafo ficará? No frontend (e enviado no momento do processamento) ou sincronizado no backend?
Resposta: Aproveitamento a capacidade de processamento da linguagem e websocket, vamos manter sincronizado com o back para, caso o usuário tome um erro e o programa feche, ao voltar tenha salvo o estado o grafo (lembre-se de regular essa sincronização para não consumir processamento demais, mas manter "checkpoints" do estado do grafo de forma usável)

4. **Armazenamento Local:** Como armazenaremos os fluxos salvos pelo usuário? Arquivos `.json` soltos no sistema de arquivos ou num banco local leve como SQLite/BoltDB?
Resposta: Use um banco de dados local leve nessa primeira versão, como SQLite

5. **Pré-visualização de Imagem:** Como lidar com a exibição de imagens muito pesadas no frontend sem travar o Electron? (Sugestão: o backend gera thumbnails ou imagens em baixa resolução sob demanda para os nós intermediários).
Resposta: Perfeito, vamos usar a sugestão, usando o backend para gerar imagens que tenham uma qualidade minima para visualizar, mas não consuma muito (lembrando que essa regra não se aplica quando o usuário for fazer download do resultado do fluxo, ai as imagens precisam vir na resolução original ou conveniente)

6. **Processamento Assíncrono:** Como o usuário poderá cancelar um lote de imagens no meio do processamento?
Resposta: De forma semelhante ai n8n, deve ter um botão flutuante na parde de baixo centralizado  da tela com opções de iniciar o fluxo, interrompe-lo, etc.

## Backlog / Ideias Futuras
- **Nó de IA/Machine Learning:** Integração com modelos ONNX ou TFLite para nós de remoção de fundo, super-resolução, etc.
Resposta: Gostei muito da ideia, adicione esse nó entre os listados no idea.md, vamos já adiciona-lo na primeira versão do projeto

- **Nó de Script Personalizado:** Expandir para suportar WebAssembly (Wasm) ou Lua, além de Python/Golang para não depender de ambientes externos.
Resposta: Gostei também, mas vamos deixar para o futuro por enquanto

- **Nó de Comparação Visual:** Nó que recebe duas imagens e exibe um slider (before/after) para comparação direta.
Resposta: Gostei muito da ideia tbm, adicione o nó no idea.md e vamos adiciona-lo na primeira versão

- **Agrupamento de Nós (Sub-grafos):** Capacidade de selecionar um grupo de nós e transformá-los em um único nó consolidado ("Macro").
Resposta: Também acho essencial, mas vamos deixar isso para a segunda versão da aplicação por enquanto