import { describe, expect, it } from 'vitest';
import type { Edge, Node } from '@xyflow/react';
import { hasProcessingChanges } from './flowProcessing';

const img = (x: number, y: number, data: Record<string, unknown> = {}): Node => ({
  id: 'n1',
  type: 'default',
  position: { x, y },
  data: { originalType: 'Image Input', filePaths: 'a.png', ...data },
});

const edge = (id: string, source: string, target: string): Edge => ({ id, source, target });

describe('hasProcessingChanges', () => {
  it('não considera mudança apenas de posição (arrastar nó) como processamento', () => {
    const before = [img(0, 0)];
    const after = [img(120, 80)];
    expect(hasProcessingChanges(before, [], after, [])).toBe(false);
  });

  it('não considera mudança de seleção como processamento', () => {
    const before = [{ ...img(0, 0), selected: false }];
    const after = [{ ...img(0, 0), selected: true }];
    expect(hasProcessingChanges(before, [], after, [])).toBe(false);
  });

  it('não considera atualização de thumbnail (preview do backend) como processamento', () => {
    const before = [img(0, 0)];
    const after = [img(0, 0, { thumbnail: 'data:image/png;base64,xxx' })];
    expect(hasProcessingChanges(before, [], after, [])).toBe(false);
  });

  it('não considera renomear o nó (label) como processamento', () => {
    const before = [img(0, 0, { label: 'Entrada' })];
    const after = [img(0, 0, { label: 'Minha Entrada' })];
    expect(hasProcessingChanges(before, [], after, [])).toBe(false);
  });

  it('detecta troca de parâmetro de processamento', () => {
    const before = [{ ...img(0, 0), id: 'b1', data: { originalType: 'Brightness & Contrast', brightness: 0 } }];
    const after = [{ ...img(0, 0), id: 'b1', data: { originalType: 'Brightness & Contrast', brightness: 40 } }];
    expect(hasProcessingChanges(before, [], after, [])).toBe(true);
  });

  it('detecta troca da imagem de entrada (filePaths)', () => {
    const before = [img(0, 0, { filePaths: 'a.png' })];
    const after = [img(0, 0, { filePaths: 'b.png' })];
    expect(hasProcessingChanges(before, [], after, [])).toBe(true);
  });

  it('detecta nó adicionado ou removido', () => {
    const a = img(0, 0);
    const b = { ...img(10, 10), id: 'n2' };
    expect(hasProcessingChanges([a], [], [a, b], [])).toBe(true);
    expect(hasProcessingChanges([a, b], [], [a], [])).toBe(true);
  });

  it('detecta conectar e desconectar aresta', () => {
    const nodes = [img(0, 0)];
    const e = edge('e1', 'n1', 'n2');
    expect(hasProcessingChanges(nodes, [], nodes, [e])).toBe(true);
    expect(hasProcessingChanges(nodes, [e], nodes, [])).toBe(true);
  });

  it('detecta troca do alvo de uma aresta', () => {
    const nodes = [img(0, 0)];
    const before = [edge('e1', 'n1', 'n2')];
    const after = [edge('e1', 'n1', 'n3')];
    expect(hasProcessingChanges(nodes, before, nodes, after)).toBe(true);
  });

  it('retorna false quando nada muda', () => {
    const nodes = [img(0, 0)];
    const edges = [edge('e1', 'n1', 'n2')];
    expect(hasProcessingChanges(nodes, edges, nodes, edges)).toBe(false);
  });
});
