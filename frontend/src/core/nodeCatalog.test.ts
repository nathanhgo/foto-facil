import { describe, expect, it } from 'vitest';
import {
  CATEGORIES,
  LIBRARY_NODES,
  NODE_CATALOG,
  getCategoryForLabel,
  getNodeColor,
  getNodeDocs,
} from './nodeCatalog';

describe('nodeCatalog', () => {
  it('every catalog entry references a valid category id', () => {
    const validIds = new Set(CATEGORIES.map((c) => c.id));
    for (const entry of NODE_CATALOG) {
      expect(validIds.has(entry.categoryId)).toBe(true);
    }
  });

  it('every catalog entry has non-empty technical and effect docs', () => {
    for (const entry of NODE_CATALOG) {
      expect(entry.docs.technical.length).toBeGreaterThan(0);
      expect(entry.docs.effect.length).toBeGreaterThan(0);
    }
  });

  it('has no duplicate node labels', () => {
    const labels = NODE_CATALOG.map((n) => n.label);
    expect(new Set(labels).size).toBe(labels.length);
  });

  it('derives LIBRARY_NODES with matching length and category labels', () => {
    expect(LIBRARY_NODES.length).toBe(NODE_CATALOG.length);
    for (const node of LIBRARY_NODES) {
      expect(typeof node.label).toBe('string');
      expect(typeof node.category).toBe('string');
    }
  });

  it('getCategoryForLabel returns the correct category for a known node', () => {
    const cat = getCategoryForLabel('Álgebra de Bandas (NDVI)');
    expect(cat.id).toBe('remote');
  });

  it('getCategoryForLabel falls back to the first category for unknown labels', () => {
    const cat = getCategoryForLabel('Não Existe');
    expect(cat.id).toBe(CATEGORIES[0].id);
  });

  it('getNodeColor returns a dark or light color depending on the flag', () => {
    const dark = getNodeColor('Image Input', true);
    const light = getNodeColor('Image Input', false);
    expect(dark).not.toBe(light);
    expect(dark.startsWith('#')).toBe(true);
    expect(light.startsWith('#')).toBe(true);
  });

  it('getNodeDocs returns docs for a known node and undefined for unknown ones', () => {
    expect(getNodeDocs('FFT')).toBeDefined();
    expect(getNodeDocs('Não Existe')).toBeUndefined();
  });
});
