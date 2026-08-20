// Catálogo central de nós: categorias, cores por categoria e documentação
// (explicação técnica + efeito esperado) de cada nó disponível na biblioteca.
// Ao adicionar um novo nó, atualize este arquivo (ver .cursor/skills/add-processing-node/SKILL.md).

export interface NodeCategory {
  id: string;
  label: string;
  /** Cor de fundo do nó/card no modo escuro. */
  bgDark: string;
  /** Cor de fundo do nó/card no modo claro. */
  bgLight: string;
  /** Cor de destaque (usada em bordas/ícones) da categoria. */
  accent: string;
}

export interface NodeDoc {
  /** Explicação técnica do algoritmo/operação por trás do nó. */
  technical: string;
  /** Efeito esperado no resultado, em linguagem simples. */
  effect: string;
}

export interface NodeCatalogEntry {
  label: string;
  categoryId: string;
  docs: NodeDoc;
}

export const CATEGORIES: NodeCategory[] = [
  { id: 'io',        label: 'Entrada/Saída',          bgDark: '#0d2233', bgLight: '#E0F2FE', accent: '#38bdf8' },
  { id: 'adjust',     label: 'Ajustes Básicos',        bgDark: '#1f1a00', bgLight: '#FEF9C3', accent: '#eab308' },
  { id: 'transform',  label: 'Transformação',          bgDark: '#1f0d1f', bgLight: '#F3E8FF', accent: '#a855f7' },
  { id: 'filters',    label: 'Filtros & Realce',       bgDark: '#0d1f1f', bgLight: '#CFFAFE', accent: '#22d3ee' },
  { id: 'analysis',   label: 'Visualização & Análise', bgDark: '#221708', bgLight: '#FFEDD5', accent: '#fb923c' },
  { id: 'remote',     label: 'Sensoriamento Remoto',   bgDark: '#0d220d', bgLight: '#DCFCE7', accent: '#22c55e' },
  { id: 'ai',         label: 'Inteligência Artificial', bgDark: '#2e0d1a', bgLight: '#FCE7F3', accent: '#ec4899' },
  { id: 'dataset',    label: 'Dataset / ML',           bgDark: '#221a0d', bgLight: '#FEF3C7', accent: '#f59e0b' },
];

const CATEGORY_BY_ID = new Map(CATEGORIES.map((c) => [c.id, c]));

export const NODE_CATALOG: NodeCatalogEntry[] = [
  // ---------- Entrada/Saída ----------
  {
    label: 'Image Input',
    categoryId: 'io',
    docs: {
      technical: 'Carrega um ou mais arquivos de imagem do disco (via diálogo do Electron) e os injeta como imagens de origem no grafo de processamento.',
      effect: 'Define de onde vêm as imagens que serão processadas pelo restante do fluxo.',
    },
  },
  {
    label: 'Image Output',
    categoryId: 'io',
    docs: {
      technical: 'Recebe a(s) imagem(ns) resultante(s) do fluxo e as grava em disco, em uma pasta de saída escolhida pelo usuário.',
      effect: 'Salva o resultado final do processamento como arquivo de imagem.',
    },
  },
  // ---------- Ajustes Básicos ----------
  {
    label: 'Brightness & Contrast',
    categoryId: 'adjust',
    docs: {
      technical: 'Soma um deslocamento (offset) ao valor de cada canal de cor por pixel, com clamp em [0, 255].',
      effect: 'Torna a imagem mais clara (valores positivos) ou mais escura (valores negativos).',
    },
  },
  {
    label: 'Color Space',
    categoryId: 'adjust',
    docs: {
      technical: 'Converte a imagem RGB para outra representação de cor: escala de cinza (luminância), HSV, Lab ou YCbCr, mapeando os componentes resultantes de volta para canais RGB para visualização.',
      effect: 'Muda a forma como as cores/informação da imagem são representadas — útil para isolar brilho de cor, ou preparar dados para outros algoritmos.',
    },
  },
  // ---------- Transformação ----------
  {
    label: 'Crop/Resize',
    categoryId: 'transform',
    docs: {
      technical: 'Recorta uma sub-região retangular (X, Y, largura, altura) e/ou redimensiona a imagem para novas dimensões usando amostragem de vizinho mais próximo/interpolação.',
      effect: 'Corta partes indesejadas da imagem e/ou muda seu tamanho final.',
    },
  },
  {
    label: 'Rotate/Flip',
    categoryId: 'transform',
    docs: {
      technical: 'Aplica rotação em múltiplos de 90° e/ou espelhamento (flip) horizontal, vertical ou em ambos os eixos.',
      effect: 'Gira e/ou espelha a imagem, útil para corrigir orientação ou gerar variações.',
    },
  },
  // ---------- Filtros & Realce ----------
  {
    label: 'Convolução',
    categoryId: 'filters',
    docs: {
      technical: 'Aplica um kernel N×N definido pelo usuário a cada pixel, somando o produto do kernel pelos vizinhos (correlação/convolução discreta), com clamp nas bordas.',
      effect: 'Nó didático genérico: dependendo do kernel, pode desfocar, realçar bordas, ou simular outros filtros espaciais (base conceitual de filtros em CNNs).',
    },
  },
  {
    label: 'Equalização de Histograma',
    categoryId: 'filters',
    docs: {
      technical: 'Calcula a função de distribuição acumulada (CDF) do histograma de cada canal e usa uma LUT (look-up table) para redistribuir as intensidades de forma mais uniforme.',
      effect: 'Aumenta o contraste global da imagem, revelando detalhes em áreas muito claras ou muito escuras.',
    },
  },
  {
    label: 'Detecção de Bordas',
    categoryId: 'filters',
    docs: {
      technical: 'Aplica o operador de Sobel (magnitude do gradiente Gx/Gy) ou o operador Laplaciano (derivada de segunda ordem) sobre a luminância da imagem.',
      effect: 'Realça contornos e transições abruptas de intensidade, apagando áreas de cor uniforme.',
    },
  },
  {
    label: 'Operações Morfológicas',
    categoryId: 'filters',
    docs: {
      technical: 'Aplica erosão (mínimo local), dilatação (máximo local), abertura (erosão + dilatação) ou fechamento (dilatação + erosão) com um elemento estruturante quadrado NxN.',
      effect: 'Remove ruído fino, separa objetos conectados (erosão), preenche pequenos buracos ou conecta objetos próximos (dilatação).',
    },
  },
  {
    label: 'Blur/Suavização',
    categoryId: 'filters',
    docs: {
      technical: 'Suaviza a imagem usando filtro Gaussiano (pesos por distância), média (box filter) ou mediana (ordena vizinhança e pega o valor central), com kernel NxN configurável.',
      effect: 'Reduz ruído e detalhes finos, deixando a imagem mais "suave"; a mediana é especialmente boa contra ruído sal-e-pimenta.',
    },
  },
  {
    label: 'Ruído',
    categoryId: 'filters',
    docs: {
      technical: 'Injeta ruído Gaussiano (aditivo, distribuição normal), sal-e-pimenta (pixels extremos aleatórios) ou speckle (multiplicativo), com semente determinística para reprodutibilidade.',
      effect: 'Degrada propositalmente a imagem — usado para testar robustez de outros filtros ou gerar dados sintéticos de treinamento.',
    },
  },
  {
    label: 'Limiarização',
    categoryId: 'filters',
    docs: {
      technical: 'Converte a imagem para preto e branco usando um corte fixo definido manualmente ou o método de Otsu (que escolhe o limiar que minimiza a variância intra-classe do histograma).',
      effect: 'Segmenta a imagem em duas regiões (primeiro plano/fundo), útil para isolar objetos de interesse.',
    },
  },
  {
    label: 'FFT',
    categoryId: 'filters',
    docs: {
      technical: 'Calcula a Transformada Rápida de Fourier 2D (Cooley-Tukey radix-2) da luminância da imagem. Modo "Espectro" renderiza a magnitude log-escalada; modo "Filtro" aplica um filtro ideal passa-baixa/passa-alta no domínio da frequência e reconstrói a imagem via FFT inversa.',
      effect: 'Revela os componentes de frequência (padrões repetitivos/texturas) da imagem, ou permite suavizar (passa-baixa) ou realçar detalhes/bordas (passa-alta) filtrando frequências.',
    },
  },
  // ---------- Visualização & Análise ----------
  {
    label: 'Comparação Visual',
    categoryId: 'analysis',
    docs: {
      technical: 'Renderiza um slider interativo que sobrepõe a imagem antes/depois do processamento, dividido por uma linha arrastável.',
      effect: 'Permite comparar visualmente o resultado de um nó com a imagem original de forma interativa.',
    },
  },
  {
    label: 'Histograma',
    categoryId: 'analysis',
    docs: {
      technical: 'Conta a frequência de cada intensidade (0-255) por canal de cor e renderiza um gráfico de barras RGB sobreposto.',
      effect: 'Mostra a distribuição de tons/cores da imagem — picos concentrados indicam baixo contraste, distribuição espalhada indica alto contraste.',
    },
  },
  {
    label: 'Matriz de Pixels',
    categoryId: 'analysis',
    docs: {
      technical: 'Extrai uma sub-região pequena (X, Y, tamanho) e renderiza uma grade com o valor numérico de luminância de cada pixel, usando fonte bitmap.',
      effect: 'Nó didático que mostra "por dentro" como uma imagem é apenas uma matriz de números — útil para explicar processamento digital de imagens.',
    },
  },
  {
    label: 'Estatísticas da Imagem',
    categoryId: 'analysis',
    docs: {
      technical: 'Calcula média, mediana e desvio padrão da luminância dos pixels. Quando há uma imagem original disponível no contexto, também calcula o Erro Quadrático Médio (MSE) e a Relação Sinal-Ruído de Pico (PSNR) em relação a ela.',
      effect: 'Fornece métricas quantitativas sobre brilho/contraste da imagem e, quando aplicável, o quanto um filtro alterou/degradou a imagem original.',
    },
  },
  // ---------- Sensoriamento Remoto ----------
  {
    label: 'Detecção de Mudanças',
    categoryId: 'remote',
    docs: {
      technical: 'Compara duas imagens pixel a pixel (mesma cena, momentos diferentes), calcula a diferença de luminância e destaca em vermelho as regiões cuja diferença excede um limiar configurável.',
      effect: 'Evidencia o que mudou entre duas capturas — simula, por exemplo, monitoramento de desmatamento por imagens de satélite.',
    },
  },
  {
    label: 'Empilhamento de Imagens',
    categoryId: 'remote',
    docs: {
      technical: 'Recebe 2 ou mais imagens da mesma cena e calcula a média aritmética por pixel entre todas elas.',
      effect: 'Reduz ruído aleatório da captura, técnica clássica em astrofotografia e calibração de sensores de satélite.',
    },
  },
  {
    label: 'Registro de Imagens',
    categoryId: 'remote',
    docs: {
      technical: 'Busca exaustivamente, dentro de um raio máximo configurável, a translação (dx, dy) que minimiza a Soma das Diferenças Absolutas (SAD) entre a luminância da imagem "móvel" e a "referência".',
      effect: 'Alinha espacialmente duas imagens levemente deslocadas — essencial antes de comparar ou combinar passagens diferentes de um mesmo satélite/sensor.',
    },
  },
  {
    label: 'Composição de Bandas',
    categoryId: 'remote',
    docs: {
      technical: 'Recebe 3 imagens em escala de cinza (bandas espectrais) e as combina, uma em cada canal R, G e B, formando uma única imagem colorida.',
      effect: 'Gera composições "falsa-cor", técnica usada em sensoriamento remoto para visualizar bandas fora do espectro visível (ex.: infravermelho).',
    },
  },
  {
    label: 'Álgebra de Bandas (NDVI)',
    categoryId: 'remote',
    docs: {
      technical: 'Calcula uma operação pixel a pixel entre duas bandas: diferença normalizada ((A-B)/(A+B), no estilo do índice NDVI), diferença simples ou soma.',
      effect: 'Realça índices espectrais como vegetação (NDVI) — áreas com alta diferença normalizada tendem a indicar vegetação saudável em imagens de satélite reais.',
    },
  },
  {
    label: 'Pan-sharpening',
    categoryId: 'remote',
    docs: {
      technical: 'Funde uma imagem colorida de baixa resolução com uma imagem pancromática de alta resolução usando a técnica de razão de Brovey, transferindo o detalhe espacial da banda pancromática para as cores.',
      effect: 'Produz uma imagem colorida com resolução espacial mais alta — técnica clássica de satélites de observação da Terra (ex.: Landsat, CBERS).',
    },
  },
  // ---------- Inteligência Artificial ----------
  {
    label: 'IA/Machine Learning',
    categoryId: 'ai',
    docs: {
      technical: 'Aplica heurísticas de remoção de fundo (baseada em tolerância de cor) ou melhoria de contraste inteligente sobre a imagem.',
      effect: 'Remove o fundo da imagem automaticamente ou melhora seu contraste de forma adaptativa, sem ajuste manual de parâmetros complexos.',
    },
  },
  // ---------- Dataset / ML ----------
  {
    label: 'Normalização',
    categoryId: 'dataset',
    docs: {
      technical: 'Reescala as intensidades de cada canal por min-max (estica para o intervalo 0-255) ou z-score (centraliza na média, mapeando ±3 desvios padrão para 0-255).',
      effect: 'Padroniza o brilho/contraste entre imagens de um dataset — etapa comum de pré-processamento antes de treinar modelos de IA.',
    },
  },
  {
    label: 'Augmentação de Dados',
    categoryId: 'dataset',
    docs: {
      technical: 'Aplica transformações aleatórias (mas reprodutíveis via seed) de rotação, flips e ruído Gaussiano a cada imagem do lote.',
      effect: 'Gera variações sintéticas do dataset original, ajudando modelos de IA a generalizar melhor com menos dados.',
    },
  },
  {
    label: 'Exportação Rotulada',
    categoryId: 'dataset',
    docs: {
      technical: 'Distribui as imagens do lote em subpastas train/val/test, com proporções configuráveis, dividindo deterministicamente pelo índice da imagem.',
      effect: 'Organiza automaticamente um dataset já pronto para ser consumido por um pipeline de treinamento de IA.',
    },
  },
];

export const LIBRARY_NODES = NODE_CATALOG.map((n) => ({
  label: n.label,
  category: CATEGORY_BY_ID.get(n.categoryId)!.label,
}));

export function getCategory(categoryId: string): NodeCategory {
  return CATEGORY_BY_ID.get(categoryId) ?? CATEGORIES[0];
}

export function getCategoryForLabel(label: string): NodeCategory {
  const entry = NODE_CATALOG.find((n) => n.label === label);
  return entry ? getCategory(entry.categoryId) : CATEGORIES[0];
}

export function getNodeColor(label: string, darkMode: boolean): string {
  const cat = getCategoryForLabel(label);
  return darkMode ? cat.bgDark : cat.bgLight;
}

export function getNodeDocs(label: string): NodeDoc | undefined {
  return NODE_CATALOG.find((n) => n.label === label)?.docs;
}
