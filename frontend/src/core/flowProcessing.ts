import type { Edge, Node } from '@xyflow/react';

/**
 * Assinatura do fluxo que importa para o resultado do processamento.
 *
 * Mover/arrastar um nó, selecionar, renomear (label) ou receber thumbnail do
 * backend não alteram o resultado do grafo — só a posição/layout. Comparar as
 * assinaturas antes/depois permite salvar o layout sem re-executar o pipeline.
 *
 * Campos de processamento (ex.: parâmetros do nó, filePaths, tipo do nó) e a
 * topologia das arestas continuam sendo considerados.
 */

const NON_PROCESSING_DATA_KEYS = new Set(['thumbnail', 'label']);

function stableData(data: unknown): Record<string, unknown> {
  const out: Record<string, unknown> = {};
  if (!data || typeof data !== 'object') return out;
  const record = data as Record<string, unknown>;
  for (const key of Object.keys(record).sort()) {
    if (NON_PROCESSING_DATA_KEYS.has(key)) continue;
    out[key] = record[key];
  }
  return out;
}

function nodesSignature(nodes: Node[]): string {
  return JSON.stringify(
    nodes.map((n) => ({
      id: n.id,
      type: n.type ?? null,
      data: stableData(n.data),
    })),
  );
}

function edgesSignature(edges: Edge[]): string {
  return JSON.stringify(
    edges.map((e) => ({
      id: e.id,
      source: e.source,
      target: e.target,
      sourceHandle: e.sourceHandle ?? null,
      targetHandle: e.targetHandle ?? null,
    })),
  );
}

export function hasProcessingChanges(
  prevNodes: Node[],
  prevEdges: Edge[],
  nextNodes: Node[],
  nextEdges: Edge[],
): boolean {
  return (
    nodesSignature(prevNodes) !== nodesSignature(nextNodes) ||
    edgesSignature(prevEdges) !== edgesSignature(nextEdges)
  );
}
