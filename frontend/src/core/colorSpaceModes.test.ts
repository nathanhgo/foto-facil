import { describe, expect, it } from 'vitest';
import {
  COLOR_SPACE_MODES,
  DEFAULT_COLOR_SPACE_MODE,
  applyColorSpaceMode,
  resolveColorSpaceMode,
} from './colorSpaceModes';

describe('colorSpaceModes', () => {
  it('lista exatamente os modos suportados pelo backend (ColorSpaceNode)', () => {
    expect(COLOR_SPACE_MODES.map((m) => m.value)).toEqual(['grayscale', 'hsv', 'lab', 'ycbcr']);
  });

  it('o modo padrão é grayscale', () => {
    expect(DEFAULT_COLOR_SPACE_MODE).toBe('grayscale');
  });

  it('resolve o valor a partir de data.colorMode do nó', () => {
    expect(resolveColorSpaceMode({ colorMode: 'hsv' })).toBe('hsv');
    expect(resolveColorSpaceMode({ colorMode: 'lab' })).toBe('lab');
  });

  it('cai no padrão quando data.colorMode está ausente ou é inválido', () => {
    expect(resolveColorSpaceMode(undefined)).toBe('grayscale');
    expect(resolveColorSpaceMode({})).toBe('grayscale');
    expect(resolveColorSpaceMode({ colorMode: 'cmyk' })).toBe('grayscale');
    expect(resolveColorSpaceMode({ colorMode: 42 })).toBe('grayscale');
  });

  it('propaga a escolha atualizando data.colorMode sem perder os demais campos', () => {
    const data = { originalType: 'Color Space', label: 'Cor', colorMode: 'grayscale' };
    const updated = applyColorSpaceMode(data, 'ycbcr');
    expect(updated).toEqual({ originalType: 'Color Space', label: 'Cor', colorMode: 'ycbcr' });
    expect(resolveColorSpaceMode(updated)).toBe('ycbcr');
  });
});
