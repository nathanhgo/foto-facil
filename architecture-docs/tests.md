# Plano de Testes (TDD)

Todas as features serão desenvolvidas seguindo a abordagem de Test-Driven Development (TDD).

## 1. Testes Unitários do Backend (Golang)
- **Estrutura do DAG:** Testar se o agendador consegue resolver corretamente a ordem topológica do grafo.
- **Detecção de Ciclos:** Testar se o sistema lança um erro caso o usuário tente conectar um nó criando um loop fechado (já que é um *Acyclic* Graph).
- **Processamento de Nós Específicos:**
  - `BrightnessContrastNode`: Fornecer uma matriz de imagem predefinida, aplicar o nó e validar a matriz de saída.
  - `GrayscaleNode`: Validar se uma imagem RGB de entrada retorna uma matriz com apenas 1 canal correto.
- **Concorrência:** Testar condições de corrida durante o processamento de nós de lote (`Directory Batch`).

## 2. Testes de Integração (Golang)
- Testar um fluxo completo em código: Iniciar com o nó de leitura (`Image Input`), passar por um filtro de redimensionamento (`Resize`), depois de escala de cinza (`Grayscale`), e validar se a imagem no nó de saída (`Image Output`) existe no disco temporário e está correta.

## 3. Testes do Frontend
- **Grafo:** Testar se as conexões lógicas entre os nós são criadas corretamente na interface.
- **Validação de Inputs:** Testar se o painel de parâmetros bloqueia entradas inválidas (ex: letras em campos de graus de rotação).
- **Exportação:** Validar se o botão de exportar gera um JSON correto com os dados do frontend em memória.

## Primeiros Testes a Serem Criados (Fase Inicial)
1. **Backend:** Testes para o `DAG Scheduler` (resolver dependências de execução entre nós mockados).
2. **Backend:** Teste base de uma interface `Node` genérica do sistema.
3. **Backend:** Teste de inicialização segura para carregar uma imagem no espaço de memória usando as funções encapsuladas de OpenCV.
