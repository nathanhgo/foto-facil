# Planejamento do MVP e Fases de Desenvolvimento

## Visão Geral
O desenvolvimento do **Foto Fácil** será dividido em fases incrementais. Isso nos permite ter uma base sólida rodando rápido, focando primeiro no motor central e na UI principal, e posteriormente expandindo o ecossistema de nós.

---

## Fase 1: Estrutura Base e MVP (Minimum Viable Product)
**Objetivo:** Construir o "chassi" da aplicação, garantindo que front e back se comuniquem perfeitamente usando WebSockets, que o estado do grafo seja salvo em SQLite e que a premissa principal de processamento de nós funcione perfeitamente usando a base de TDD.

### Core Backend (Golang)
- [x] Configurar roteamento WebSocket.
- [x] Implementar motor do DAG (`DAG Scheduler`) com resolução de dependências e bloqueio de ciclos.
- [x] Implementar banco de dados SQLite para salvar/carregar fluxos (`.json` na tabela).
- [x] Implementar sistema de geração de *thumbnails* em baixa resolução para nós intermediários.

### Core Frontend (Electron + React.js + TS)
- [x] Setup do Electron integrado ao React.js.
- [x] Implementar UI principal: Canvas infinito (zoom/pan) e sistema de abas.
- [x] Implementar lógica de conexão de fios entre nós.
- [x] Implementar painel de propriedades que varia os inputs conforme o nó selecionado.
- [x] Adicionar botão flutuante para "Iniciar" e "Interromper" fluxos.

### Nós Iniciais (MVP)
- [x] Image Input & Directory Batch
- [x] Image Output & Download
- [x] Brightness & Contrast
- [x] Color Space (Grayscale)
- [x] Crop/Resize
- [x] Rotate/Flip
- [x] IA/Machine Learning (Apenas um modelo base, ex: Remoção de Fundo via ONNX)
- [x] Comparação Visual (Slider before/after)

---

## Fase 2: Expansão de PDI Acadêmico e Otimização
**Objetivo:** Trazer o poder pesado do Processamento Digital de Imagem, permitindo análises complexas, filtros avançados e matrizes, suportando perfeitamente o lote (batching).

### Novidades no Core
- [ ] Implementar paralelismo/concorrência refinada no processamento de lote (Batch) em Golang.
- [ ] Agrupamento de Nós (Sub-grafos ou "Macros").

### Nós de Expansão
- Blur / Smoothing (Gaussiano, Média, Mediana)
- Noise (Gaussiano, Salt & Pepper)
- Edge Detection (Sobel, Laplacian, Canny)
- Morphological Operations (Erosão, Dilatação)
- Histogram & Statistics (Média, PSNR, etc)
- Matriz/Array de Pixels
- Domínio da Frequência (FFT)

---

## Fase 3: Personalização e Futuro
**Objetivo:** Deixar a ferramenta "hackeável" para usuários avançados.

### Novidades
- [ ] Nó de "Custom Node" (Scripting) via Python/Golang/Wasm para o usuário programar os próprios filtros.
- [ ] Exportação automatizada de pipelines.

---

## Estrutura de Desenvolvimento e Scripts

Para facilitar nosso workflow diário, a raiz do projeto conterá as seguintes configurações e automatizações:

### Tarefas do VSCode (`.vscode/tasks.json`)
A ideia é ter tarefas padronizadas no editor para levantar todo o ecossistema com um único clique (ou atalho).

- **`Start Frontend (Dev)`**: Roda o servidor React de desenvolvimento.
- **`Start Backend (Dev)`**: Roda `go run .` no diretório do backend.
- **`Start All`**: Tarefa composta que executa o Frontend, Backend e inicializa o Electron apontando para o servidor local.
- **`Run Backend Tests`**: Roda a suíte de testes `go test ./...` para validar o TDD.

### Dependências Iniciais (Docker)
Embora a aplicação final seja um executável Desktop nativo contendo os binários rodando embutidos, para fins de desenvolvimento podemos necessitar compilar bibliotecas C++ pesadas (OpenCV). 
- Caso o processo de instalar o wrapper C++ (`gocv`) seja muito complexo nativamente no Linux da máquina do desenvolvedor, podemos empacotar a compilação do Backend de dev usando um contêiner **Docker**, porém inicialmente tentaremos rodar o Go puro, visto que a arquitetura Electron costuma rodar melhor de forma nativa. O banco `SQLite` dispensa contêineres e será apenas um arquivo `app.db` na pasta AppData.

## Arquitetura de Diretórios Proposta

O repositório seguirá uma estrutura dividida claramente entre o código da interface (Frontend/Electron) e o motor de processamento (Backend Golang):

```text
foto-facil/
│
├── architecture-docs/       # Documentação de arquitetura e planejamento (onde estamos)
│
├── frontend/                # Electron + React.js (TypeScript)
│   ├── public/              # Assets estáticos
│   ├── src/                 # Código React e Electron
│   │   ├── components/      # Componentes UI (Nós, Fios, Canvas, Menus)
│   │   ├── core/            # Lógica de gerenciamento de estado do grafo e WebSockets
│   │   ├── styles/          # Estilização global
│   │   └── main.ts          # Processo principal do Electron (Main Process)
│   ├── package.json
│   └── vite.config.ts       # Ou equivalente (bundler)
│
├── backend/                 # Golang (Motor de PDI e DAG)
│   ├── cmd/                 # Ponto de entrada (main.go)
│   ├── internal/            # Lógica de negócios
│   │   ├── api/             # Handlers WebSocket e HTTP
│   │   ├── dag/             # Motor de execução do grafo e Scheduler (TDD)
│   │   ├── nodes/           # Implementação matemática de cada nó
│   │   └── storage/         # Persistência de fluxos no SQLite
│   ├── pkg/                 # Utilitários (ex: processamento paralelo de arrays de imagens)
│   ├── tests/               # Testes automatizados (TDD)
│   ├── go.mod
│   └── go.sum
│
├── .vscode/                 # Tarefas (tasks.json)
├── .env                     # Configurações do ambiente
├── .gitignore
└── README.md
```

## Variáveis de Ambiente (`.env`)

Para este projeto, como o processamento (incluindo IA e manipulações pesadas) será totalmente **local e offline** usando processamento da CPU/GPU da sua máquina e modelos baseados em ONNX/TFLite embutidos, **NÃO precisaremos de chaves de API pagas** (como OpenAI, AWS, etc.). 

O arquivo `.env` (na raiz do projeto) servirá puramente para configurações de portas e caminhos do sistema:

```env
# Ambiente
APP_ENV=development

# Comunicação
FRONTEND_PORT=3000
BACKEND_WS_PORT=8080

# Banco de dados e Cache
SQLITE_DB_PATH=./app.db
TEMP_IMAGE_DIR=./tmp_processed_images

# Caminho para os arquivos de IA (onde os .onnx vão ficar)
AI_MODELS_PATH=./backend/models
```
Ou seja, você não precisará caçar nenhuma chave na internet. Basta eu configurar isso automaticamente nos scripts e estará pronto para uso.
