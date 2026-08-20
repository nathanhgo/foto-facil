import { useCallback, useRef, useState, useEffect } from 'react';
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  BackgroundVariant,
  ReactFlowProvider,
  SelectionMode,
} from '@xyflow/react';
import type { Connection, Edge, Node, OnSelectionChangeParams } from '@xyflow/react';
import { getBackendWebSocketUrl } from './core/websocketUrl';
import { CATEGORIES, LIBRARY_NODES, getNodeColor, getNodeDocs } from './core/nodeCatalog';

declare global {
  interface Window {
    electronAPI?: {
      openFileDialog: () => Promise<string[] | null>;
      openDirectoryDialog: () => Promise<string | null>;
    };
  }
}

let idCounter = 0;
const getId = () => `dndnode_${idCounter++}_${Date.now()}`;

// Cores, categorias e documentação de cada nó agora vivem em `src/core/nodeCatalog.ts`
// (import LIBRARY_NODES, getNodeColor, getCategory, getNodeDocs acima).

// ---------- CompareSlider ----------
function CompareSlider({ thumbnail }: { thumbnail: string }) {
  const [pos, setPos] = useState(50);
  const [width, setWidth] = useState(250);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!containerRef.current) return;
    const obs = new ResizeObserver((entries) => {
      for (const entry of entries) {
        if (entry.contentRect.width > 0) {
          setWidth(entry.contentRect.width);
        }
      }
    });
    obs.observe(containerRef.current);
    return () => obs.disconnect();
  }, []);

  const onMouseMove = (e: React.MouseEvent) => {
    if (e.buttons !== 1 || !containerRef.current) return;
    const rect = containerRef.current.getBoundingClientRect();
    const pct = Math.round(((e.clientX - rect.left) / rect.width) * 100);
    setPos(Math.max(0, Math.min(100, pct)));
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
      <label style={{ fontSize: '11px', color: '#888', textTransform: 'uppercase', letterSpacing: '0.5px' }}>
        Antes / Depois ({pos}%)
      </label>
      <div
        ref={containerRef}
        onMouseMove={onMouseMove}
        style={{
          position: 'relative',
          width: '100%',
          aspectRatio: '1',
          cursor: 'col-resize',
          borderRadius: '4px',
          overflow: 'hidden',
          border: '1px solid #444',
          userSelect: 'none',
        }}
      >
        {/* After image (Right half) */}
        <div style={{ position: 'absolute', inset: 0, overflow: 'hidden' }}>
          <img
            src={thumbnail}
            alt="after"
            draggable={false}
            onDragStart={(e) => e.preventDefault()}
            style={{
              position: 'absolute',
              left: -width,
              width: width * 2,
              height: '100%',
              objectFit: 'cover',
              userSelect: 'none',
              pointerEvents: 'none',
            }}
          />
        </div>

        {/* Before image (Left half, clipped by pos%) */}
        <div
          style={{
            position: 'absolute',
            left: 0,
            top: 0,
            bottom: 0,
            width: `${pos}%`,
            overflow: 'hidden',
          }}
        >
          <img
            src={thumbnail}
            alt="before"
            draggable={false}
            onDragStart={(e) => e.preventDefault()}
            style={{
              position: 'absolute',
              left: 0,
              width: width * 2,
              height: '100%',
              objectFit: 'cover',
              userSelect: 'none',
              pointerEvents: 'none',
            }}
          />
        </div>

        {/* Divider line */}
        <div
          style={{
            position: 'absolute',
            top: 0,
            bottom: 0,
            left: `${pos}%`,
            transform: 'translateX(-50%)',
            width: '2px',
            backgroundColor: '#fff',
            boxShadow: '0 0 4px rgba(0,0,0,0.5)',
            pointerEvents: 'none',
          }}
        >
          <div
            style={{
              position: 'absolute',
              top: '50%',
              left: '50%',
              transform: 'translate(-50%,-50%)',
              width: '20px',
              height: '20px',
              borderRadius: '50%',
              backgroundColor: '#fff',
              boxShadow: '0 0 4px rgba(0,0,0,0.5)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '10px',
              color: '#333',
            }}
          >
            ⟷
          </div>
        </div>
      </div>
    </div>
  );
}

// ---------- PropertiesPanel ----------
function PropertiesPanel({
  selectedNode,
  selectedCount,
  setNodes,
  edges,
  triggerPreview,
  theme,
}: {
  selectedNode: Node | null;
  selectedCount: number;
  setNodes: any;
  edges: Edge[];
  triggerPreview: (nodes: Node[], edges: Edge[]) => void;
  theme: any;
}) {
  if (selectedCount > 1) {
    return (
      <div style={{ padding: '16px' }}>
        <h3 style={{ color: '#666', fontSize: '11px', fontWeight: 700, letterSpacing: '1px', margin: 0 }}>PROPERTIES</h3>
        <div style={{ backgroundColor: theme.cardBg, padding: '14px', borderRadius: '6px', border: `1px solid ${theme.cardBorder}`, marginTop: '12px' }}>
          <p style={{ fontSize: '13px', fontWeight: 600, margin: '0 0 8px 0', color: theme.textColor }}>{selectedCount} nós selecionados</p>
          <p style={{ fontSize: '11px', color: '#888', margin: 0, lineHeight: 1.5 }}>
            Ctrl+C / Ctrl+V para copiar, Delete/Backspace para excluir todos os nós selecionados.
          </p>
        </div>
      </div>
    );
  }

  if (!selectedNode) {
    return (
      <div style={{ padding: '16px' }}>
        <h3 style={{ color: '#666', fontSize: '11px', fontWeight: 700, letterSpacing: '1px', margin: 0 }}>PROPERTIES</h3>
        <p style={{ fontSize: '12px', color: '#888', marginTop: '12px' }}>Selecione um nó no canvas.</p>
      </div>
    );
  }

  const updateNodeData = (updates: Record<string, any>) => {
    setNodes((nds: Node[]) => {
      const updatedNodes = nds.map((n) => (n.id === selectedNode.id ? { ...n, data: { ...n.data, ...updates } } : n));
      triggerPreview(updatedNodes, edges);
      return updatedNodes;
    });
    // `selectedNode` é derivado de `nodes` + `selectedIds` no componente pai,
    // então já reflete a atualização automaticamente no próximo render.
  };

  const nodeType = selectedNode.data?.originalType as string;
  const thumb = selectedNode.data?.thumbnail as string | undefined;

  const selectStyle: React.CSSProperties = {
    padding: '6px',
    backgroundColor: theme.inputBg,
    color: theme.inputColor,
    border: `1px solid ${theme.inputBorder}`,
    borderRadius: '4px',
    width: '100%',
  };

  return (
    <div style={{ padding: '16px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
      <h3 style={{ color: '#666', fontSize: '11px', fontWeight: 700, letterSpacing: '1px', margin: 0 }}>PROPERTIES</h3>
      <div style={{ backgroundColor: theme.cardBg, padding: '14px', borderRadius: '6px', border: `1px solid ${theme.cardBorder}` }}>
        <p style={{ fontSize: '14px', fontWeight: 600, margin: '0 0 12px 0', color: theme.textColor }}>{nodeType}</p>

        {nodeType === 'Image Input' && (
          <PropSection label="Imagem(ns) de Entrada">
            {(selectedNode.data?.filePaths as string) && (
              <p style={{ fontSize: '11px', color: '#9b51e0', margin: '0 0 6px 0', wordBreak: 'break-all' }}>
                {selectedNode.data.filePaths as string}
              </p>
            )}
            <ActionButton label="Selecionar Arquivo(s)..." onClick={async () => {
              if (window.electronAPI) {
                const paths = await window.electronAPI.openFileDialog();
                if (paths) updateNodeData({ filePaths: paths.join(',') });
              }
            }} />
          </PropSection>
        )}

        {nodeType === 'Image Output' && (
          <PropSection label="Pasta de Saída">
            {(selectedNode.data?.outputDir as string) && (
              <p style={{ fontSize: '11px', color: '#9b51e0', margin: '0 0 6px 0', wordBreak: 'break-all' }}>
                {selectedNode.data.outputDir as string}
              </p>
            )}
            <ActionButton label="Selecionar Pasta..." onClick={async () => {
              if (window.electronAPI) {
                const dir = await window.electronAPI.openDirectoryDialog();
                if (dir) updateNodeData({ outputDir: dir });
              }
            }} />
          </PropSection>
        )}

        {nodeType === 'Brightness & Contrast' && (
          <PropSection label={`Brightness (${(selectedNode.data?.brightness as number) ?? 0})`}>
            <input type="range" min="-255" max="255"
              value={(selectedNode.data?.brightness as number) ?? 0}
              onChange={(e) => updateNodeData({ brightness: parseInt(e.target.value) })}
              style={{ width: '100%', accentColor: '#9b51e0' }}
            />
          </PropSection>
        )}

        {nodeType === 'Color Space' && (
          <PropSection label="Modo">
            <select value="grayscale" style={selectStyle} onChange={() => {}}>
              <option value="grayscale">Grayscale</option>
            </select>
          </PropSection>
        )}

        {nodeType === 'Crop/Resize' && (
          <>
            <PropSection label="Redimensionar para">
              <div style={{ display: 'flex', gap: '8px' }}>
                <NumberInput label="W" value={(selectedNode.data?.resizeW as number) ?? 0} onChange={(v) => updateNodeData({ resizeW: v })} theme={theme} />
                <NumberInput label="H" value={(selectedNode.data?.resizeH as number) ?? 0} onChange={(v) => updateNodeData({ resizeH: v })} theme={theme} />
              </div>
            </PropSection>
            <PropSection label="Recorte (Crop)">
              <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
                <NumberInput label="X" value={(selectedNode.data?.cropX as number) ?? 0} onChange={(v) => updateNodeData({ cropX: v })} theme={theme} />
                <NumberInput label="Y" value={(selectedNode.data?.cropY as number) ?? 0} onChange={(v) => updateNodeData({ cropY: v })} theme={theme} />
                <NumberInput label="W" value={(selectedNode.data?.cropW as number) ?? 0} onChange={(v) => updateNodeData({ cropW: v })} theme={theme} />
                <NumberInput label="H" value={(selectedNode.data?.cropH as number) ?? 0} onChange={(v) => updateNodeData({ cropH: v })} theme={theme} />
              </div>
            </PropSection>
          </>
        )}

        {nodeType === 'Rotate/Flip' && (
          <>
            <PropSection label="Rotação">
              <select value={(selectedNode.data?.rotateAngle as number) ?? 0}
                onChange={(e) => updateNodeData({ rotateAngle: parseInt(e.target.value) })}
                style={selectStyle}>
                <option value={0}>Sem rotação</option>
                <option value={90}>90°</option>
                <option value={180}>180°</option>
                <option value={270}>270°</option>
              </select>
            </PropSection>
            <PropSection label="Flip">
              <select value={(selectedNode.data?.flipAxis as number) ?? -1}
                onChange={(e) => {
                  const v = parseInt(e.target.value);
                  updateNodeData({ flipAxis: v, doFlip: v !== -1 });
                }}
                style={selectStyle}>
                <option value={-1}>Sem flip</option>
                <option value={0}>Horizontal</option>
                <option value={1}>Vertical</option>
                <option value={2}>Ambos</option>
              </select>
            </PropSection>
          </>
        )}

        {nodeType === 'IA/Machine Learning' && (
          <>
            <PropSection label="Ferramenta de IA">
              <select
                value={(selectedNode.data?.aiTool as string) ?? 'background_removal'}
                onChange={(e) => updateNodeData({ aiTool: e.target.value })}
                style={selectStyle}
              >
                <option value="background_removal">Remoção de Fundo (AI Background Removal)</option>
                <option value="contrast_boost">Melhoria de Contraste Inteligente</option>
              </select>
            </PropSection>
            {((selectedNode.data?.aiTool as string) ?? 'background_removal') === 'background_removal' && (
              <PropSection label={`Tolerância (${(selectedNode.data?.aiTolerance as number) ?? 80})`}>
                <input type="range" min="10" max="250"
                  value={(selectedNode.data?.aiTolerance as number) ?? 80}
                  onChange={(e) => updateNodeData({ aiTolerance: parseInt(e.target.value) })}
                  style={{ width: '100%', accentColor: '#9b51e0' }}
                />
              </PropSection>
            )}
          </>
        )}

        {nodeType === 'Detecção de Mudanças' && (
          <PropSection label={`Limiar de Mudança (${(selectedNode.data?.changeThreshold as number) ?? 30})`}>
            <input type="range" min="1" max="255"
              value={(selectedNode.data?.changeThreshold as number) ?? 30}
              onChange={(e) => updateNodeData({ changeThreshold: parseInt(e.target.value) })}
              style={{ width: '100%', accentColor: '#9b51e0' }}
            />
            <p style={{ fontSize: '11px', color: '#888', margin: '4px 0 0 0' }}>
              Conecte duas entradas (antes/depois) — ex.: duas passagens de satélite sobre a mesma área.
            </p>
          </PropSection>
        )}

        {nodeType === 'Empilhamento de Imagens' && (
          <PropSection label="Empilhamento (Stacking)">
            <p style={{ fontSize: '11px', color: '#888', margin: 0 }}>
              Conecte 2+ entradas do mesmo cenário; a saída é a média por pixel (reduz ruído).
            </p>
          </PropSection>
        )}

        {nodeType === 'Registro de Imagens' && (
          <PropSection label={`Busca máx. de deslocamento (${(selectedNode.data?.maxShift as number) ?? 10}px)`}>
            <input type="range" min="1" max="30"
              value={(selectedNode.data?.maxShift as number) ?? 10}
              onChange={(e) => updateNodeData({ maxShift: parseInt(e.target.value) })}
              style={{ width: '100%', accentColor: '#9b51e0' }}
            />
            <p style={{ fontSize: '11px', color: '#888', margin: '4px 0 0 0' }}>
              Entrada 1 = referência, Entrada 2 = imagem a alinhar.
            </p>
          </PropSection>
        )}

        {nodeType === 'Composição de Bandas' && (
          <PropSection label="Composição Falso-Cor">
            <p style={{ fontSize: '11px', color: '#888', margin: 0 }}>
              Conecte 3 entradas (bandas) → mapeadas para R, G e B.
            </p>
          </PropSection>
        )}

        {nodeType === 'Álgebra de Bandas (NDVI)' && (
          <PropSection label="Operação">
            <select
              value={(selectedNode.data?.bandMathOperation as string) ?? 'normalized_difference'}
              onChange={(e) => updateNodeData({ bandMathOperation: e.target.value })}
              style={selectStyle}
            >
              <option value="normalized_difference">Diferença Normalizada (estilo NDVI)</option>
              <option value="difference">Diferença Simples</option>
              <option value="sum">Soma</option>
            </select>
            <p style={{ fontSize: '11px', color: '#888', margin: '4px 0 0 0' }}>
              Entrada 1 = banda A, Entrada 2 = banda B.
            </p>
          </PropSection>
        )}

        {nodeType === 'Pan-sharpening' && (
          <PropSection label="Fusão Pancromática">
            <p style={{ fontSize: '11px', color: '#888', margin: 0 }}>
              Entrada 1 = imagem colorida (baixa resolução), Entrada 2 = pancromática (alta resolução).
            </p>
          </PropSection>
        )}

        {nodeType === 'Convolução' && (
          <>
            <PropSection label="Tamanho do Kernel">
              <NumberInput label="N (NxN)" value={(selectedNode.data?.kernelSize as number) ?? 3} onChange={(v) => updateNodeData({ kernelSize: v })} theme={theme} />
            </PropSection>
            <PropSection label="Kernel (valores separados por vírgula)">
              <input
                type="text"
                value={(selectedNode.data?.kernel as string) ?? '0,0,0,0,1,0,0,0,0'}
                onChange={(e) => updateNodeData({ kernel: e.target.value })}
                placeholder="0,0,0,0,1,0,0,0,0"
                style={selectStyle}
              />
            </PropSection>
          </>
        )}

        {nodeType === 'Matriz de Pixels' && (
          <PropSection label="Região a Exibir">
            <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
              <NumberInput label="X" value={(selectedNode.data?.regionX as number) ?? 0} onChange={(v) => updateNodeData({ regionX: v })} theme={theme} />
              <NumberInput label="Y" value={(selectedNode.data?.regionY as number) ?? 0} onChange={(v) => updateNodeData({ regionY: v })} theme={theme} />
              <NumberInput label="Tamanho" value={(selectedNode.data?.regionSize as number) ?? 8} onChange={(v) => updateNodeData({ regionSize: v })} theme={theme} />
            </div>
          </PropSection>
        )}

        {nodeType === 'Detecção de Bordas' && (
          <PropSection label="Método">
            <select
              value={(selectedNode.data?.edgeMethod as string) ?? 'sobel'}
              onChange={(e) => updateNodeData({ edgeMethod: e.target.value })}
              style={selectStyle}
            >
              <option value="sobel">Sobel (gradiente)</option>
              <option value="laplacian">Laplaciano</option>
            </select>
          </PropSection>
        )}

        {nodeType === 'Operações Morfológicas' && (
          <>
            <PropSection label="Operação">
              <select
                value={(selectedNode.data?.morphOperation as string) ?? 'erosion'}
                onChange={(e) => updateNodeData({ morphOperation: e.target.value })}
                style={selectStyle}
              >
                <option value="erosion">Erosão</option>
                <option value="dilation">Dilatação</option>
                <option value="opening">Abertura</option>
                <option value="closing">Fechamento</option>
              </select>
            </PropSection>
            <PropSection label="Elemento Estruturante">
              <NumberInput label="N (NxN)" value={(selectedNode.data?.kernelSize as number) ?? 3} onChange={(v) => updateNodeData({ kernelSize: v })} theme={theme} />
            </PropSection>
          </>
        )}

        {nodeType === 'Blur/Suavização' && (
          <>
            <PropSection label="Método">
              <select
                value={(selectedNode.data?.blurMethod as string) ?? 'gaussian'}
                onChange={(e) => updateNodeData({ blurMethod: e.target.value })}
                style={selectStyle}
              >
                <option value="gaussian">Gaussiano</option>
                <option value="mean">Média</option>
                <option value="median">Mediana</option>
              </select>
            </PropSection>
            <PropSection label="Tamanho do Kernel">
              <NumberInput label="N (NxN)" value={(selectedNode.data?.kernelSize as number) ?? 3} onChange={(v) => updateNodeData({ kernelSize: v })} theme={theme} />
            </PropSection>
          </>
        )}

        {nodeType === 'Ruído' && (
          <>
            <PropSection label="Tipo">
              <select
                value={(selectedNode.data?.noiseType as string) ?? 'gaussian'}
                onChange={(e) => updateNodeData({ noiseType: e.target.value })}
                style={selectStyle}
              >
                <option value="gaussian">Gaussiano</option>
                <option value="salt_pepper">Sal e Pimenta</option>
                <option value="speckle">Speckle</option>
              </select>
            </PropSection>
            <PropSection label={`Intensidade (${(selectedNode.data?.noiseAmount as number) ?? 20})`}>
              <input type="range" min="1" max="100"
                value={(selectedNode.data?.noiseAmount as number) ?? 20}
                onChange={(e) => updateNodeData({ noiseAmount: parseInt(e.target.value) })}
                style={{ width: '100%', accentColor: '#9b51e0' }}
              />
            </PropSection>
          </>
        )}

        {nodeType === 'Limiarização' && (
          <>
            <PropSection label="Método">
              <select
                value={(selectedNode.data?.thresholdMethod as string) ?? 'global'}
                onChange={(e) => updateNodeData({ thresholdMethod: e.target.value })}
                style={selectStyle}
              >
                <option value="global">Global (manual)</option>
                <option value="otsu">Otsu (automático)</option>
              </select>
            </PropSection>
            {((selectedNode.data?.thresholdMethod as string) ?? 'global') === 'global' && (
              <PropSection label={`Corte (${(selectedNode.data?.thresholdValue as number) ?? 128})`}>
                <input type="range" min="1" max="255"
                  value={(selectedNode.data?.thresholdValue as number) ?? 128}
                  onChange={(e) => updateNodeData({ thresholdValue: parseInt(e.target.value) })}
                  style={{ width: '100%', accentColor: '#9b51e0' }}
                />
              </PropSection>
            )}
          </>
        )}

        {nodeType === 'FFT' && (
          <>
            <PropSection label="Modo">
              <select
                value={(selectedNode.data?.fftMode as string) ?? 'spectrum'}
                onChange={(e) => updateNodeData({ fftMode: e.target.value })}
                style={selectStyle}
              >
                <option value="spectrum">Espectro de Magnitude</option>
                <option value="filter">Filtro em Frequência</option>
              </select>
            </PropSection>
            {((selectedNode.data?.fftMode as string) ?? 'spectrum') === 'filter' && (
              <>
                <PropSection label="Filtro">
                  <select
                    value={(selectedNode.data?.fftFilter as string) ?? 'lowpass'}
                    onChange={(e) => updateNodeData({ fftFilter: e.target.value })}
                    style={selectStyle}
                  >
                    <option value="lowpass">Passa-Baixa</option>
                    <option value="highpass">Passa-Alta</option>
                  </select>
                </PropSection>
                <PropSection label={`Corte (${((selectedNode.data?.cutoffRatio as number) ?? 0.2).toFixed(2)})`}>
                  <input type="range" min="0.05" max="0.95" step="0.05"
                    value={(selectedNode.data?.cutoffRatio as number) ?? 0.2}
                    onChange={(e) => updateNodeData({ cutoffRatio: parseFloat(e.target.value) })}
                    style={{ width: '100%', accentColor: '#9b51e0' }}
                  />
                </PropSection>
              </>
            )}
          </>
        )}

        {nodeType === 'Normalização' && (
          <PropSection label="Método">
            <select
              value={(selectedNode.data?.normalizationMethod as string) ?? 'minmax'}
              onChange={(e) => updateNodeData({ normalizationMethod: e.target.value })}
              style={selectStyle}
            >
              <option value="minmax">Min-Max (0-255)</option>
              <option value="zscore">Z-Score (±3σ)</option>
            </select>
          </PropSection>
        )}

        {nodeType === 'Augmentação de Dados' && (
          <PropSection label="Seed (reprodutibilidade)">
            <NumberInput label="Seed" value={(selectedNode.data?.augmentationSeed as number) ?? 42} onChange={(v) => updateNodeData({ augmentationSeed: v })} theme={theme} />
          </PropSection>
        )}

        {nodeType === 'Exportação Rotulada' && (
          <>
            <PropSection label="Pasta de Saída">
              {(selectedNode.data?.outputDir as string) && (
                <p style={{ fontSize: '11px', color: '#9b51e0', margin: '0 0 6px 0', wordBreak: 'break-all' }}>
                  {selectedNode.data.outputDir as string}
                </p>
              )}
              <ActionButton label="Selecionar Pasta..." onClick={async () => {
                if (window.electronAPI) {
                  const dir = await window.electronAPI.openDirectoryDialog();
                  if (dir) updateNodeData({ outputDir: dir });
                }
              }} />
            </PropSection>
            <PropSection label="Divisão Train/Val/Test">
              <div style={{ display: 'flex', gap: '8px' }}>
                <NumberInput label="Train %" value={Math.round(((selectedNode.data?.trainRatio as number) ?? 0.7) * 100)} onChange={(v) => updateNodeData({ trainRatio: v / 100 })} theme={theme} />
                <NumberInput label="Val %" value={Math.round(((selectedNode.data?.valRatio as number) ?? 0.15) * 100)} onChange={(v) => updateNodeData({ valRatio: v / 100 })} theme={theme} />
              </div>
            </PropSection>
          </>
        )}

        {(nodeType === 'Histograma' || nodeType === 'Equalização de Histograma' || nodeType === 'Estatísticas da Imagem') && (
          <PropSection label="Sem parâmetros">
            <p style={{ fontSize: '11px', color: '#888', margin: 0 }}>Este nó processa a imagem automaticamente, sem configuração adicional.</p>
          </PropSection>
        )}

        {nodeType === 'Comparação Visual' && thumb && (
          <CompareSlider thumbnail={thumb} />
        )}

        {nodeType !== 'Comparação Visual' && thumb && (
          <PropSection label="Preview">
            <img src={thumb} alt="preview"
              style={{ width: '100%', borderRadius: '4px', border: '1px solid #444', display: 'block', pointerEvents: 'none', userSelect: 'none' }} />
          </PropSection>
        )}

        <NodeDocsSection nodeType={nodeType} theme={theme} />
      </div>
    </div>
  );
}

// ---------- NodeDocsSection ----------
function NodeDocsSection({ nodeType, theme }: { nodeType: string; theme: any }) {
  const docs = getNodeDocs(nodeType);
  const [open, setOpen] = useState(false);
  if (!docs) return null;

  return (
    <div style={{ marginTop: '4px', borderTop: `1px solid ${theme.cardBorder}`, paddingTop: '10px' }}>
      <button
        onClick={() => setOpen((v) => !v)}
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          width: '100%',
          background: 'transparent',
          border: 'none',
          cursor: 'pointer',
          padding: 0,
          color: '#888',
          fontSize: '11px',
          fontWeight: 700,
          letterSpacing: '0.5px',
          textTransform: 'uppercase',
        }}
      >
        <span>Sobre este nó</span>
        <span style={{ fontSize: '12px' }}>{open ? '▾' : '▸'}</span>
      </button>
      {open && (
        <div style={{ marginTop: '8px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
          <div>
            <p style={{ fontSize: '10px', color: '#888', textTransform: 'uppercase', letterSpacing: '0.5px', margin: '0 0 3px 0' }}>Explicação técnica</p>
            <p style={{ fontSize: '12px', color: theme.textColor, margin: 0, lineHeight: 1.5 }}>{docs.technical}</p>
          </div>
          <div>
            <p style={{ fontSize: '10px', color: '#888', textTransform: 'uppercase', letterSpacing: '0.5px', margin: '0 0 3px 0' }}>Efeito esperado</p>
            <p style={{ fontSize: '12px', color: theme.textColor, margin: 0, lineHeight: 1.5 }}>{docs.effect}</p>
          </div>
        </div>
      )}
    </div>
  );
}

function PropSection({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '5px', marginBottom: '10px' }}>
      <label style={{ fontSize: '11px', color: '#888', textTransform: 'uppercase', letterSpacing: '0.5px' }}>{label}</label>
      {children}
    </div>
  );
}

function ActionButton({ label, onClick }: { label: string; onClick: () => void }) {
  return (
    <button onClick={onClick} style={{ padding: '7px 12px', backgroundColor: '#9b51e0', color: '#fff', border: 'none', borderRadius: '4px', cursor: 'pointer', fontSize: '13px', fontWeight: 600 }}>
      {label}
    </button>
  );
}

function NumberInput({ label, value, onChange, theme }: { label: string; value: number; onChange: (v: number) => void; theme: any }) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2px', flex: 1 }}>
      <label style={{ fontSize: '10px', color: '#888' }}>{label}</label>
      <input type="number" value={value} min={0}
        onChange={(e) => onChange(parseInt(e.target.value) || 0)}
        style={{
          width: '100%',
          padding: '4px 6px',
          backgroundColor: theme.inputBg,
          color: theme.inputColor,
          border: `1px solid ${theme.inputBorder}`,
          borderRadius: '4px',
          fontSize: '13px',
        }}
      />
    </div>
  );
}

// ---------- Main Flow ----------
function Flow() {
  const reactFlowWrapper = useRef<HTMLDivElement>(null);
  
  // Theme state
  const [darkMode, setDarkMode] = useState<boolean>(true);

  // Flows / Tabs state (Excel Sheets Style)
  const [flows, setFlows] = useState<Array<{ id: string; name: string; nodes: Node[]; edges: Edge[] }>>([
    { id: 'flow-1', name: 'Flow 1', nodes: [], edges: [] }
  ]);
  const [activeFlowId, setActiveFlowId] = useState<string>('flow-1');
  const [editingFlowId, setEditingFlowId] = useState<string | null>(null);

  // React Flow states
  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);
  const [reactFlowInstance, setReactFlowInstance] = useState<any>(null);
  const [runStatus, setRunStatus] = useState<'idle' | 'running' | 'done' | 'error'>('idle');

  // Selection state — guardamos apenas os IDs selecionados e derivamos os nós
  // sempre a partir do estado bruto `nodes` (nunca de `displayNodes`, que envolve
  // o label em JSX para permitir renomear inline — gravar isso de volta no estado
  // via copiar/colar corrompia o nó e quebrava a renderização da tela toda).
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const selectedNodes = nodes.filter((n) => selectedIds.has(n.id));
  const selectedNode = selectedNodes.length === 1 ? selectedNodes[0] : null;
  const clipboardRef = useRef<Node[]>([]);

  // Resizable side panels (sidebar de nós / propriedades)
  const [leftWidth, setLeftWidth] = useState<number>(200);
  const [rightWidth, setRightWidth] = useState<number>(270);
  const [resizingPanel, setResizingPanel] = useState<'left' | 'right' | null>(null);

  // ---------- Undo/Redo (Ctrl+Z / Ctrl+Y) ----------
  // Histórico simples de snapshots de nodes/edges, "assentado" após 500ms de
  // inatividade (evita empilhar um passo de undo por frame durante um arraste).
  const historyRef = useRef<{ past: Array<{ nodes: Node[]; edges: Edge[] }>; future: Array<{ nodes: Node[]; edges: Edge[] }> }>({ past: [], future: [] });
  const committedStateRef = useRef<{ nodes: Node[]; edges: Edge[] }>({ nodes: [], edges: [] });
  const skipHistoryRef = useRef(false);
  const historyDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  // Contador que só existe para forçar um re-render quando o histórico muda,
  // já que mutações em `historyRef.current` (um ref) não disparam re-render sozinhas.
  const [, setHistoryVersion] = useState(0);
  const bumpHistoryVersion = useCallback(() => setHistoryVersion((v) => v + 1), []);

  // Node editing state
  const [editingNodeId, setEditingNodeId] = useState<string | null>(null);

  const handleRenameNode = (id: string, newLabel: string) => {
    setEditingNodeId(null);
    if (!newLabel.trim()) return;
    setNodes((nds) =>
      nds.map((n) => (n.id === id ? { ...n, data: { ...n.data, label: newLabel } } : n))
    );
  };

  const displayNodes = nodes.map((n) => {
    const isEditing = editingNodeId === n.id;
    return {
      ...n,
      data: {
        ...n.data,
        label: isEditing ? (
          <input
            type="text"
            defaultValue={n.data.label as string}
            autoFocus
            onBlur={(e) => handleRenameNode(n.id, e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') handleRenameNode(n.id, e.currentTarget.value);
            }}
            onClick={(e) => e.stopPropagation()}
            onMouseDown={(e) => e.stopPropagation()}
            onPointerDown={(e) => e.stopPropagation()}
            style={{
              border: 'none',
              background: 'transparent',
              color: darkMode ? '#fff' : '#000',
              width: '100%',
              outline: 'none',
              textAlign: 'center',
              fontSize: '13px',
              fontWeight: 500,
            }}
          />
        ) : (
          <span>{n.data.label as string}</span>
        )
      }
    };
  });

  // Context Menu state
  const [menu, setMenu] = useState<{ id: string; type: 'node' | 'edge'; x: number; y: number } | null>(null);

  const wsRef = useRef<WebSocket | null>(null);
  const debounceSavePreviewRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Theme object
  const theme = {
    bgHeader: darkMode ? '#161616' : '#F1F5F9',
    bgSidebar: darkMode ? '#141414' : '#E2E8F0',
    bgCanvas: darkMode ? '#0E0E0E' : '#F8FAFC',
    textColor: darkMode ? '#E0E0E0' : '#0F172A',
    borderColor: darkMode ? '#252525' : '#CBD5E1',
    cardBg: darkMode ? '#1E1E1E' : '#FFFFFF',
    cardBorder: darkMode ? '#2A2A2A' : '#E2E8F0',
    inputBg: darkMode ? '#111' : '#FFF',
    inputColor: darkMode ? '#EEE' : '#0F172A',
    inputBorder: darkMode ? '#444' : '#94A3B8',
  };

  // ---------- Painéis redimensionáveis ----------
  // Arraste nos divisores entre sidebar/canvas e canvas/propriedades para
  // aumentar ou diminuir a largura de cada área.
  const LEFT_MIN = 160, LEFT_MAX = 420, RIGHT_MIN = 220, RIGHT_MAX = 480;

  const startResizing = useCallback((panel: 'left' | 'right') => (e: React.MouseEvent) => {
    e.preventDefault();
    setResizingPanel(panel);
  }, []);

  useEffect(() => {
    if (!resizingPanel) return;

    const handleMouseMove = (e: MouseEvent) => {
      if (resizingPanel === 'left') {
        setLeftWidth(Math.min(LEFT_MAX, Math.max(LEFT_MIN, e.clientX)));
      } else {
        setRightWidth(Math.min(RIGHT_MAX, Math.max(RIGHT_MIN, window.innerWidth - e.clientX)));
      }
    };
    const handleMouseUp = () => setResizingPanel(null);

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [resizingPanel]);

  // ---------- Desfazer/Refazer (Ctrl+Z / Ctrl+Y) ----------
  // Se houver um commit de histórico pendente (dentro da janela de debounce de
  // 500ms), assenta ele imediatamente antes de desfazer/refazer — senão o undo
  // "pularia" o passo mais recente, que ainda não tinha sido empilhado.
  const flushPendingHistory = useCallback(() => {
    if (!historyDebounceRef.current) return;
    clearTimeout(historyDebounceRef.current);
    historyDebounceRef.current = null;
    historyRef.current.past.push(committedStateRef.current);
    if (historyRef.current.past.length > 100) historyRef.current.past.shift();
    historyRef.current.future = [];
    committedStateRef.current = { nodes, edges };
  }, [nodes, edges]);

  const handleUndo = useCallback(() => {
    flushPendingHistory();
    const { past } = historyRef.current;
    if (past.length === 0) return;
    const previous = past.pop()!;
    historyRef.current.future.push({ nodes, edges });
    skipHistoryRef.current = true;
    setNodes(previous.nodes);
    setEdges(previous.edges);
    setSelectedIds(new Set());
    bumpHistoryVersion();
  }, [nodes, edges, setNodes, setEdges, bumpHistoryVersion, flushPendingHistory]);

  const handleRedo = useCallback(() => {
    flushPendingHistory();
    const { future } = historyRef.current;
    if (future.length === 0) return;
    const next = future.pop()!;
    historyRef.current.past.push({ nodes, edges });
    skipHistoryRef.current = true;
    setNodes(next.nodes);
    setEdges(next.edges);
    setSelectedIds(new Set());
    bumpHistoryVersion();
  }, [nodes, edges, setNodes, setEdges, bumpHistoryVersion, flushPendingHistory]);

  // Assenta um snapshot no histórico 500ms após a última mudança de nodes/edges
  // (evita um passo de undo por frame durante um arraste ou por tecla digitada).
  useEffect(() => {
    if (skipHistoryRef.current) {
      skipHistoryRef.current = false;
      committedStateRef.current = { nodes, edges };
      return;
    }
    const baseline = committedStateRef.current;
    if (historyDebounceRef.current) clearTimeout(historyDebounceRef.current);
    historyDebounceRef.current = setTimeout(() => {
      historyRef.current.past.push(baseline);
      if (historyRef.current.past.length > 100) historyRef.current.past.shift();
      historyRef.current.future = [];
      committedStateRef.current = { nodes, edges };
      bumpHistoryVersion();
    }, 500);
  }, [nodes, edges, bumpHistoryVersion]);

  // ---------- Copiar/colar nós (Ctrl+C / Ctrl+V) e Desfazer/Refazer (Ctrl+Z / Ctrl+Y) ----------
  useEffect(() => {
    const isEditableTarget = (el: EventTarget | null) => {
      if (!(el instanceof HTMLElement)) return false;
      return el.tagName === 'INPUT' || el.tagName === 'TEXTAREA' || el.isContentEditable;
    };

    const handleKeyDown = (e: KeyboardEvent) => {
      if (!(e.ctrlKey || e.metaKey) || isEditableTarget(e.target)) return;

      if (e.key === 'c' || e.key === 'C') {
        if (selectedNodes.length > 0) {
          clipboardRef.current = selectedNodes;
        }
      } else if (e.key === 'v' || e.key === 'V') {
        if (clipboardRef.current.length === 0) return;
        e.preventDefault();
        const offset = 40;
        const pasted = clipboardRef.current.map((n) => ({
          ...n,
          id: getId(),
          selected: true,
          position: { x: n.position.x + offset, y: n.position.y + offset },
          data: { ...n.data, thumbnail: undefined },
        }));
        setNodes((nds) => nds.map((n) => ({ ...n, selected: false })).concat(pasted));
        setSelectedIds(new Set(pasted.map((n) => n.id)));
        clipboardRef.current = pasted;
      } else if (e.key === 'z' || e.key === 'Z') {
        e.preventDefault();
        if (e.shiftKey) handleRedo();
        else handleUndo();
      } else if (e.key === 'y' || e.key === 'Y') {
        e.preventDefault();
        handleRedo();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [selectedNodes, setNodes, handleUndo, handleRedo]);

  // Helper to send flow state to save and preview
  const triggerSaveAndPreview = useCallback((latestNodes: Node[], latestEdges: Edge[]) => {
    if (debounceSavePreviewRef.current) clearTimeout(debounceSavePreviewRef.current);
    debounceSavePreviewRef.current = setTimeout(() => {
      const ws = wsRef.current;
      if (!ws || ws.readyState !== WebSocket.OPEN) return;

      const activeName = flows.find(f => f.id === activeFlowId)?.name || 'Flow';

      // 1. SAVE_FLOW
      ws.send(JSON.stringify({
        action: 'SAVE_FLOW',
        id: activeFlowId,
        name: activeName,
        flow: { nodes: latestNodes, edges: latestEdges }
      }));

      // 2. RUN_FLOW (Preview)
      const hasInput = latestNodes.some(
        (n) => n.data?.originalType === 'Image Input' && (n.data?.filePaths as string)
      );
      if (hasInput) {
        ws.send(JSON.stringify({ action: 'RUN_FLOW', flow: { nodes: latestNodes, edges: latestEdges } }));
      }
    }, 400);
  }, [activeFlowId, flows]);

  // WebSocket connection & initial loading
  useEffect(() => {
    function connect() {
      const ws = new WebSocket(getBackendWebSocketUrl());
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('[WS] Connected');
        // Initial fetch from SQLite database
        ws.send(JSON.stringify({ action: 'LOAD_FLOWS' }));
      };

      ws.onmessage = (event) => {
        try {
          const response = JSON.parse(event.data);
          
          if (response.action === 'LOAD_FLOWS' && response.status === 'success') {
            const rawFlows = response.flows || [];
            const loaded = rawFlows.map((f: any) => {
              try {
                const parsed = JSON.parse(f.data);
                const validatedNodes = (parsed.nodes || []).map((n: any) => {
                  const pos = n.position && typeof n.position.x === 'number' && typeof n.position.y === 'number'
                    ? n.position
                    : { x: 0, y: 0 };
                  return { ...n, position: pos };
                });
                return {
                  id: f.id,
                  name: f.name,
                  nodes: validatedNodes,
                  edges: parsed.edges || [],
                };
              } catch {
                return { id: f.id, name: f.name, nodes: [], edges: [] };
              }
            });
            if (loaded.length > 0) {
              setFlows(loaded);
              setActiveFlowId(loaded[0].id);
            } else {
              // Create default flow
              const defaultFlow = { id: 'flow-1', name: 'Flow 1', nodes: [], edges: [] };
              setFlows([defaultFlow]);
              setActiveFlowId('flow-1');
              ws.send(JSON.stringify({
                action: 'SAVE_FLOW',
                id: defaultFlow.id,
                name: defaultFlow.name,
                flow: { nodes: [], edges: [] }
              }));
            }
          } else if (response.thumbnails) {
            // Atualização de preview vinda do backend — não é uma edição do
            // usuário, não deve virar um passo de undo.
            skipHistoryRef.current = true;
            setNodes((nds) =>
              nds.map((n) => {
                const thumb = response.thumbnails[n.id];
                return thumb ? { ...n, data: { ...n.data, thumbnail: thumb } } : n;
              })
            );
          }
        } catch (err) {
          console.error('[WS] Error processing message:', err);
        }
      };

      ws.onclose = () => {
        setTimeout(connect, 2000);
      };
      ws.onerror = () => {
        ws.close();
      };
    }
    connect();
    return () => {
      wsRef.current?.close();
    };
  }, [setNodes]);

  // Load active flow nodes/edges when activeFlowId changes
  useEffect(() => {
    const active = flows.find((f) => f.id === activeFlowId);
    if (active) {
      // Troca de fluxo não deve virar um passo de "desfazer" nem herdar o
      // histórico do fluxo anterior.
      skipHistoryRef.current = true;
      historyRef.current = { past: [], future: [] };
      setNodes(active.nodes || []);
      setEdges(active.edges || []);
      setSelectedIds(new Set());
    }
  }, [activeFlowId]);

  // Update node styles dynamically when darkMode changes
  useEffect(() => {
    // Alternar tema é só cosmético — não deve virar um passo de undo.
    skipHistoryRef.current = true;
    setNodes((nds) =>
      nds.map((n) => {
        const label = (n.data?.originalType as string) ?? (n.data?.label as string);
        if (!label) return n;
        const bgColor = getNodeColor(label, darkMode);
        const textColor = darkMode ? '#E0E0E0' : '#0F172A';
        const borderColor = darkMode ? '#444' : '#CBD5E1';
        return {
          ...n,
          style: {
            ...n.style,
            backgroundColor: bgColor,
            color: textColor,
            border: `1px solid ${borderColor}`,
          },
        };
      })
    );
  }, [darkMode, setNodes]);

  // Save changes to database and flows array when nodes/edges are mutated
  useEffect(() => {
    const active = flows.find((f) => f.id === activeFlowId);
    if (!active) return;

    const nodesChanged = JSON.stringify(active.nodes) !== JSON.stringify(nodes);
    const edgesChanged = JSON.stringify(active.edges) !== JSON.stringify(edges);

    if (nodesChanged || edgesChanged) {
      setFlows((prev) =>
        prev.map((f) => (f.id === activeFlowId ? { ...f, nodes, edges } : f))
      );
      triggerSaveAndPreview(nodes, edges);
    }
  }, [nodes, edges, activeFlowId, triggerSaveAndPreview]);

  const onSelectionChange = useCallback(({ nodes: sel }: OnSelectionChangeParams) => {
    setSelectedIds(new Set(sel.map((n) => n.id)));
  }, []);

  const onConnect = useCallback(
    (params: Connection | Edge) => {
      setEdges((eds) => {
        const newEdges = addEdge({ ...params, animated: true }, eds);
        return newEdges;
      });
    },
    [setEdges]
  );

  const onDragStart = (e: React.DragEvent, label: string) => {
    e.dataTransfer.setData('application/reactflow', 'default');
    e.dataTransfer.setData('application/reactflow/label', label);
    e.dataTransfer.effectAllowed = 'move';
  };

  const onDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    const label = e.dataTransfer.getData('application/reactflow/label');
    if (!label || !reactFlowInstance) return;
    const position = reactFlowInstance.screenToFlowPosition({ x: e.clientX, y: e.clientY });
    
    const bgColor = getNodeColor(label, darkMode);
    const textColor = darkMode ? '#E0E0E0' : '#0F172A';
    const borderColor = darkMode ? '#444' : '#CBD5E1';

    const newNode: Node = {
      id: getId(),
      type: 'default',
      position,
      data: { label, originalType: label },
      style: {
        backgroundColor: bgColor,
        border: `1px solid ${borderColor}`,
        borderRadius: '6px',
        color: textColor,
        fontSize: '13px'
      },
    };
    setNodes((nds) => nds.concat(newNode));
  }, [reactFlowInstance, setNodes, darkMode]);

  // Context Menu handlers
  const onNodeContextMenu = useCallback(
    (event: React.MouseEvent, node: Node) => {
      event.preventDefault();
      setMenu({
        id: node.id,
        type: 'node',
        x: event.clientX,
        y: event.clientY,
      });
    },
    [setMenu]
  );

  const onEdgeContextMenu = useCallback(
    (event: React.MouseEvent, edge: Edge) => {
      event.preventDefault();
      setMenu({
        id: edge.id,
        type: 'edge',
        x: event.clientX,
        y: event.clientY,
      });
    },
    [setMenu]
  );

  // Close context menu on click anywhere
  useEffect(() => {
    const handleClose = () => setMenu(null);
    document.addEventListener('click', handleClose);
    return () => document.removeEventListener('click', handleClose);
  }, []);

  // Delete node or edge
  const handleDeleteElement = useCallback(() => {
    if (!menu) return;
    if (menu.type === 'node') {
      setNodes((nds) => nds.filter((n) => n.id !== menu.id));
      setEdges((eds) => eds.filter((e) => e.source !== menu.id && e.target !== menu.id));
      setSelectedIds((prev) => {
        if (!prev.has(menu.id)) return prev;
        const next = new Set(prev);
        next.delete(menu.id);
        return next;
      });
    } else {
      setEdges((eds) => eds.filter((e) => e.id !== menu.id));
    }
    setMenu(null);
  }, [menu, setNodes, setEdges]);

  const handleRunFlow = useCallback(() => {
    setRunStatus('running');
    triggerSaveAndPreview(nodes, edges);
    setTimeout(() => setRunStatus('done'), 1500);
  }, [nodes, edges, triggerSaveAndPreview]);

  // Tab operations
  const handleAddFlow = () => {
    const newId = `flow-${Date.now()}`;
    const newName = `Flow ${flows.length + 1}`;
    const newFlow = { id: newId, name: newName, nodes: [], edges: [] };

    setFlows((prev) => [...prev, newFlow]);
    setActiveFlowId(newId);

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        action: 'SAVE_FLOW',
        id: newId,
        name: newName,
        flow: { nodes: [], edges: [] }
      }));
    }
  };

  const handleDeleteFlow = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (flows.length <= 1) return;

    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        action: 'DELETE_FLOW',
        id
      }));
    }

    const nextFlows = flows.filter((f) => f.id !== id);
    setFlows(nextFlows);

    if (activeFlowId === id) {
      setActiveFlowId(nextFlows[0].id);
    }
  };

  const handleRenameFlowSubmit = (id: string, newName: string) => {
    setEditingFlowId(null);
    if (!newName.trim()) return;
    setFlows((prev) =>
      prev.map((f) => (f.id === id ? { ...f, name: newName } : f))
    );

    const targetFlow = flows.find((f) => f.id === id);
    if (targetFlow && wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        action: 'SAVE_FLOW',
        id,
        name: newName,
        flow: { nodes: targetFlow.nodes, edges: targetFlow.edges }
      }));
    }
  };

  const statusColor = { idle: '#9b51e0', running: '#f5a623', done: '#27ae60', error: '#e74c3c' };
  const statusLabel = { idle: 'Rodar Fluxo', running: 'Processando…', done: '✓ Concluído', error: '✗ Erro' };
  // Ordem fixa das categorias (ver src/core/nodeCatalog.ts) — uma cor por categoria.
  const categories = CATEGORIES.filter((c) => LIBRARY_NODES.some((n) => n.category === c.label));

  return (
    <div
      className={darkMode ? '' : 'theme-light'}
      style={{ width: '100vw', height: '100vh', display: 'flex', flexDirection: 'column', backgroundColor: theme.bgCanvas, color: theme.textColor, fontFamily: 'Inter, system-ui, sans-serif', cursor: resizingPanel ? 'col-resize' : undefined }}
    >
      
      {/* Header */}
      <div style={{ height: '52px', backgroundColor: theme.bgHeader, borderBottom: `1px solid ${theme.borderColor}`, display: 'flex', alignItems: 'center', padding: '0 16px', justifyContent: 'space-between', flexShrink: 0 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <span style={{ fontSize: '15px', fontWeight: 700, color: '#a855f7', letterSpacing: '-0.3px' }}>Foto Fácil</span>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          {/* Undo/Redo */}
          <div style={{ display: 'flex', gap: '2px' }}>
            <button
              onClick={handleUndo}
              disabled={historyRef.current.past.length === 0}
              title="Desfazer (Ctrl+Z)"
              style={{
                backgroundColor: darkMode ? '#222' : '#E2E8F0',
                border: `1px solid ${theme.borderColor}`,
                borderRadius: '6px 0 0 6px',
                padding: '6px 10px',
                cursor: 'pointer',
                fontSize: '14px',
                color: theme.textColor,
                opacity: historyRef.current.past.length === 0 ? 0.4 : 1,
              }}
            >
              ↶
            </button>
            <button
              onClick={handleRedo}
              disabled={historyRef.current.future.length === 0}
              title="Refazer (Ctrl+Y)"
              style={{
                backgroundColor: darkMode ? '#222' : '#E2E8F0',
                border: `1px solid ${theme.borderColor}`,
                borderLeft: 'none',
                borderRadius: '0 6px 6px 0',
                padding: '6px 10px',
                cursor: 'pointer',
                fontSize: '14px',
                color: theme.textColor,
                opacity: historyRef.current.future.length === 0 ? 0.4 : 1,
              }}
            >
              ↷
            </button>
          </div>

          {/* Dark Mode toggle */}
          <button
            onClick={() => setDarkMode(!darkMode)}
            title="Alternar Tema (Dark/Light)"
            style={{
              backgroundColor: darkMode ? '#222' : '#E2E8F0',
              border: `1px solid ${theme.borderColor}`,
              borderRadius: '6px',
              padding: '6px 12px',
              cursor: 'pointer',
              fontSize: '13px',
              color: theme.textColor,
              display: 'flex',
              alignItems: 'center',
              gap: '6px',
              fontWeight: 500,
            }}
          >
            {darkMode ? '☀️ Modo Claro' : '🌙 Modo Escuro'}
          </button>

          <button onClick={handleRunFlow} disabled={runStatus === 'running'}
            style={{ backgroundColor: statusColor[runStatus], color: '#fff', padding: '7px 20px', borderRadius: '6px', fontSize: '13px', fontWeight: 600, border: 'none', cursor: runStatus === 'running' ? 'wait' : 'pointer', transition: 'background-color 0.2s' }}>
            {statusLabel[runStatus]}
          </button>
        </div>
      </div>

      {/* Body */}
      <div style={{ flex: 1, display: 'flex', overflow: 'hidden', position: 'relative' }}>
        
        {/* Left sidebar */}
        <div
          className="discreet-scrollbar"
          style={{ width: `${leftWidth}px`, backgroundColor: theme.bgSidebar, borderRight: `1px solid ${theme.borderColor}`, padding: '12px', display: 'flex', flexDirection: 'column', gap: '16px', flexShrink: 0, overflowY: 'auto' }}
        >
          {categories.map((cat) => {
            const catAccent = cat.accent;
            return (
              <div key={cat.id}>
                <p style={{ fontSize: '10px', color: catAccent, fontWeight: 700, letterSpacing: '1px', marginBottom: '6px' }}>{cat.label.toUpperCase()}</p>
                {LIBRARY_NODES.filter((n) => n.category === cat.label).map((n) => {
                  const bgColor = getNodeColor(n.label, darkMode);
                  const textColor = darkMode ? '#CCC' : '#0F172A';
                  return (
                    <div key={n.label} draggable onDragStart={(e) => onDragStart(e, n.label)}
                      style={{ backgroundColor: bgColor, padding: '7px 10px', borderRadius: '5px', cursor: 'grab', fontSize: '13px', borderLeft: `3px solid ${catAccent}`, userSelect: 'none', color: textColor, marginBottom: '4px' }}>
                      {n.label}
                    </div>
                  );
                })}
              </div>
            );
          })}
        </div>

        {/* Divisor: sidebar <-> canvas */}
        <div
          className={`resize-divider${resizingPanel === 'left' ? ' resizing' : ''}`}
          onMouseDown={startResizing('left')}
        />

        {/* Canvas Area */}
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', height: '100%', minWidth: 0 }}>
          <div style={{ flex: 1 }} ref={reactFlowWrapper}>
            <ReactFlow
              nodes={displayNodes}
              edges={edges}
              onNodesChange={onNodesChange}
              onEdgesChange={onEdgesChange}
              onConnect={onConnect}
              onSelectionChange={onSelectionChange}
              onNodeContextMenu={onNodeContextMenu}
              onEdgeContextMenu={onEdgeContextMenu}
              onInit={setReactFlowInstance}
              onDrop={onDrop}
              onDragOver={onDragOver}
              onNodeDoubleClick={(_event, node) => setEditingNodeId(node.id)}
              colorMode={darkMode ? 'dark' : 'light'}
              proOptions={{ hideAttribution: true }}
              deleteKeyCode={['Backspace', 'Delete']}
              selectionOnDrag
              selectionMode={SelectionMode.Partial}
              panOnDrag={[1, 2]}
              fitView
            >
              <Controls
                style={{
                  backgroundColor: darkMode ? '#1E1E1E' : '#FFFFFF',
                  border: `1px solid ${theme.borderColor}`,
                  color: theme.textColor,
                }}
              />
              <MiniMap
                nodeStrokeWidth={3}
                nodeColor={darkMode ? '#1A1A1A' : '#E2E8F0'}
                maskColor={darkMode ? '#0E0E0E99' : '#F1F5F999'}
                style={{
                  backgroundColor: darkMode ? '#141414' : '#FFFFFF',
                  border: `1px solid ${theme.borderColor}`,
                }}
              />
              <Background variant={BackgroundVariant.Dots} gap={20} size={1} color={darkMode ? '#222' : '#CBD5E1'} />
            </ReactFlow>
          </div>

          {/* Bottom Tabs (Excel Sheets Style) */}
          <div style={{
            height: '38px',
            backgroundColor: theme.bgHeader,
            borderTop: `1px solid ${theme.borderColor}`,
            display: 'flex',
            alignItems: 'center',
            padding: '0 12px',
            gap: '2px',
            zIndex: 10,
            overflowX: 'auto',
          }}>
            {flows.map((f) => {
              const isActive = f.id === activeFlowId;
              return (
                <div
                  key={f.id}
                  onClick={() => setActiveFlowId(f.id)}
                  onDoubleClick={() => setEditingFlowId(f.id)}
                  style={{
                    height: '100%',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '8px',
                    padding: '0 16px',
                    backgroundColor: isActive ? theme.bgCanvas : 'transparent',
                    borderTop: isActive ? '3px solid #a855f7' : '3px solid transparent',
                    borderRight: `1px solid ${theme.borderColor}`,
                    cursor: 'pointer',
                    fontSize: '13px',
                    fontWeight: isActive ? 600 : 400,
                    color: isActive ? theme.textColor : '#777',
                    userSelect: 'none',
                    transition: 'background-color 0.2s',
                  }}
                  onMouseEnter={(e) => {
                    if (!isActive) e.currentTarget.style.backgroundColor = darkMode ? '#222' : '#E2E8F0';
                  }}
                  onMouseLeave={(e) => {
                    if (!isActive) e.currentTarget.style.backgroundColor = 'transparent';
                  }}
                >
                  {editingFlowId === f.id ? (
                    <input
                      type="text"
                      defaultValue={f.name}
                      autoFocus
                      onBlur={(e) => handleRenameFlowSubmit(f.id, e.target.value)}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter') handleRenameFlowSubmit(f.id, e.currentTarget.value);
                      }}
                      style={{
                        border: 'none',
                        background: 'transparent',
                        color: theme.textColor,
                        width: '80px',
                        outline: 'none',
                        fontSize: '13px',
                      }}
                    />
                  ) : (
                    <span>{f.name}</span>
                  )}
                  {flows.length > 1 && (
                    <span
                      onClick={(e) => handleDeleteFlow(f.id, e)}
                      style={{
                        fontSize: '11px',
                        fontWeight: 'bold',
                        color: '#999',
                        marginLeft: '4px',
                        padding: '2px 5px',
                        borderRadius: '3px',
                      }}
                      onMouseEnter={(e) => e.currentTarget.style.color = '#ff4d4d'}
                      onMouseLeave={(e) => e.currentTarget.style.color = '#999'}
                    >
                      ✕
                    </span>
                  )}
                </div>
              );
            })}
            <button
              onClick={handleAddFlow}
              title="Criar novo fluxo"
              style={{
                height: '100%',
                padding: '0 14px',
                backgroundColor: 'transparent',
                border: 'none',
                cursor: 'pointer',
                color: '#a855f7',
                fontSize: '16px',
                fontWeight: 'bold',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
              onMouseEnter={(e) => e.currentTarget.style.backgroundColor = darkMode ? '#222' : '#E2E8F0'}
              onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
            >
              +
            </button>
          </div>
        </div>

        {/* Divisor: canvas <-> propriedades */}
        <div
          className={`resize-divider${resizingPanel === 'right' ? ' resizing' : ''}`}
          onMouseDown={startResizing('right')}
        />

        {/* Right Properties Panel */}
        <div
          className="discreet-scrollbar"
          style={{ width: `${rightWidth}px`, backgroundColor: theme.bgSidebar, borderLeft: `1px solid ${theme.borderColor}`, flexShrink: 0, overflowY: 'auto' }}
        >
          <PropertiesPanel
            selectedNode={selectedNode}
            selectedCount={selectedNodes.length}
            setNodes={setNodes}
            edges={edges}
            triggerPreview={triggerSaveAndPreview}
            theme={theme}
          />
        </div>
      </div>

      {/* Floating Context Menu */}
      {menu && (
        <div
          style={{
            position: 'fixed',
            top: menu.y,
            left: menu.x,
            backgroundColor: theme.cardBg,
            border: `1px solid ${theme.cardBorder}`,
            borderRadius: '4px',
            boxShadow: '0 2px 10px rgba(0,0,0,0.5)',
            zIndex: 1000,
            padding: '4px 0',
          }}
        >
          <button
            onClick={handleDeleteElement}
            style={{
              display: 'block',
              width: '100%',
              padding: '6px 20px',
              backgroundColor: 'transparent',
              color: '#ff4d4d',
              border: 'none',
              textAlign: 'left',
              cursor: 'pointer',
              fontSize: '13px',
              fontWeight: 500,
            }}
            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = darkMode ? '#222' : '#E2E8F0'}
            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
          >
            Excluir {menu.type === 'node' ? 'Nó' : 'Conexão'}
          </button>
        </div>
      )}
    </div>
  );
}

export default function App() {
  return (
    <ReactFlowProvider>
      <Flow />
    </ReactFlowProvider>
  );
}
