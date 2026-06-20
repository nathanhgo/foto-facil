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
} from '@xyflow/react';
import type { Connection, Edge, Node, OnSelectionChangeParams } from '@xyflow/react';

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

const NODE_COLORS_DARK: Record<string, string> = {
  'Image Input':           '#0d2233',
  'Brightness & Contrast': '#1f1a00',
  'Color Space':           '#0d1f0d',
  'Crop/Resize':           '#1f0d1f',
  'Rotate/Flip':           '#0d1f1f',
  'Image Output':          '#1a0d2e',
  'Comparação Visual':     '#1a1a0d',
  'IA/Machine Learning':   '#2e0d1a',
};

const NODE_COLORS_LIGHT: Record<string, string> = {
  'Image Input':           '#E0F2FE', // azul claro
  'Brightness & Contrast': '#FEF9C3', // amarelo claro
  'Color Space':           '#DCFCE7', // verde claro
  'Crop/Resize':           '#F3E8FF', // roxo claro
  'Rotate/Flip':           '#E0F7FA', // ciano claro
  'Image Output':          '#F3E8FF', // roxo claro
  'Comparação Visual':     '#FEF9C3', // amarelo claro
  'IA/Machine Learning':   '#FCE7F3', // rosa claro
};

const LIBRARY_NODES = [
  { label: 'Image Input',           category: 'Input/Output' },
  { label: 'Image Output',          category: 'Input/Output' },
  { label: 'Brightness & Contrast', category: 'Ajustes' },
  { label: 'Color Space',           category: 'Ajustes' },
  { label: 'Crop/Resize',           category: 'Transformação' },
  { label: 'Rotate/Flip',           category: 'Transformação' },
  { label: 'Comparação Visual',     category: 'Visualização' },
  { label: 'IA/Machine Learning',   category: 'Inteligência Artificial' },
];

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
  setNodes,
  setSelectedNode,
  edges,
  triggerPreview,
  theme,
}: {
  selectedNode: Node | null;
  setNodes: any;
  setSelectedNode: any;
  edges: Edge[];
  triggerPreview: (nodes: Node[], edges: Edge[]) => void;
  theme: any;
}) {
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
    setSelectedNode((n: Node | null) => n ? { ...n, data: { ...n.data, ...updates } } : null);
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

        {nodeType === 'Comparação Visual' && thumb && (
          <CompareSlider thumbnail={thumb} />
        )}

        {nodeType !== 'Comparação Visual' && thumb && (
          <PropSection label="Preview">
            <img src={thumb} alt="preview"
              style={{ width: '100%', borderRadius: '4px', border: '1px solid #444', display: 'block', pointerEvents: 'none', userSelect: 'none' }} />
          </PropSection>
        )}
      </div>
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
  const [selectedNode, setSelectedNode] = useState<Node | null>(null);
  const [runStatus, setRunStatus] = useState<'idle' | 'running' | 'done' | 'error'>('idle');

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
      const ws = new WebSocket('ws://localhost:8080/ws');
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
      setNodes(active.nodes || []);
      setEdges(active.edges || []);
      setSelectedNode(null);
    }
  }, [activeFlowId]);

  // Update node styles dynamically when darkMode changes
  useEffect(() => {
    setNodes((nds) =>
      nds.map((n) => {
        const label = n.data?.label as string;
        if (!label) return n;
        const bgColor = darkMode
          ? (NODE_COLORS_DARK[label] || '#1E1E1E')
          : (NODE_COLORS_LIGHT[label] || '#FFFFFF');
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
    setSelectedNode(sel.length > 0 ? sel[0] : null);
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
    
    const bgColor = darkMode ? (NODE_COLORS_DARK[label] || '#1E1E1E') : (NODE_COLORS_LIGHT[label] || '#FFFFFF');
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
      if (selectedNode?.id === menu.id) {
        setSelectedNode(null);
      }
    } else {
      setEdges((eds) => eds.filter((e) => e.id !== menu.id));
    }
    setMenu(null);
  }, [menu, setNodes, setEdges, selectedNode]);

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
  const categories = [...new Set(LIBRARY_NODES.map((n) => n.category))];

  return (
    <div style={{ width: '100vw', height: '100vh', display: 'flex', flexDirection: 'column', backgroundColor: theme.bgCanvas, color: theme.textColor, fontFamily: 'Inter, system-ui, sans-serif' }}>
      
      {/* Header */}
      <div style={{ height: '52px', backgroundColor: theme.bgHeader, borderBottom: `1px solid ${theme.borderColor}`, display: 'flex', alignItems: 'center', padding: '0 16px', justifyContent: 'space-between', flexShrink: 0 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <span style={{ fontSize: '15px', fontWeight: 700, color: '#a855f7', letterSpacing: '-0.3px' }}>Foto Fácil</span>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
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
        <div style={{ width: '200px', backgroundColor: theme.bgSidebar, borderRight: `1px solid ${theme.borderColor}`, padding: '12px', display: 'flex', flexDirection: 'column', gap: '16px', flexShrink: 0, overflowY: 'auto' }}>
          {categories.map((cat) => (
            <div key={cat}>
              <p style={{ fontSize: '10px', color: '#888', fontWeight: 700, letterSpacing: '1px', marginBottom: '6px' }}>{cat.toUpperCase()}</p>
              {LIBRARY_NODES.filter((n) => n.category === cat).map((n) => {
                const bgColor = darkMode ? (NODE_COLORS_DARK[n.label] || '#1E1E1E') : (NODE_COLORS_LIGHT[n.label] || '#FFFFFF');
                const textColor = darkMode ? '#CCC' : '#0F172A';
                const borderColor = darkMode ? '#2A2A2A' : '#CBD5E1';
                return (
                  <div key={n.label} draggable onDragStart={(e) => onDragStart(e, n.label)}
                    style={{ backgroundColor: bgColor, padding: '7px 10px', borderRadius: '5px', cursor: 'grab', fontSize: '13px', border: `1px solid ${borderColor}`, userSelect: 'none', color: textColor, marginBottom: '4px' }}>
                    {n.label}
                  </div>
                );
              })}
            </div>
          ))}
        </div>

        {/* Canvas Area */}
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', height: '100%' }}>
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

        {/* Right Properties Panel */}
        <div style={{ width: '270px', backgroundColor: theme.bgSidebar, borderLeft: `1px solid ${theme.borderColor}`, flexShrink: 0, overflowY: 'auto' }}>
          <PropertiesPanel
            selectedNode={selectedNode}
            setNodes={setNodes}
            setSelectedNode={setSelectedNode}
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
