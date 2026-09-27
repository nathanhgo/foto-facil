/**
 * Modos do nó "Color Space".
 *
 * A lista precisa espelhar exatamente o que o backend suporta em
 * `backend/internal/nodes/color_space.go` (ColorSpaceNode): "grayscale"
 * (padrão), "hsv", "lab" e "ycbcr". Não invente opções aqui sem implementá-las
 * no backend.
 */

export interface ColorSpaceMode {
  value: string;
  label: string;
}

export const COLOR_SPACE_MODES: ColorSpaceMode[] = [
  { value: 'grayscale', label: 'Grayscale' },
  { value: 'hsv', label: 'HSV' },
  { value: 'lab', label: 'Lab' },
  { value: 'ycbcr', label: 'YCbCr' },
];

export const DEFAULT_COLOR_SPACE_MODE = 'grayscale';

export function isColorSpaceMode(value: unknown): value is string {
  return typeof value === 'string' && COLOR_SPACE_MODES.some((m) => m.value === value);
}

/** Valor a exibir no select, refletindo `data.colorMode` do nó selecionado. */
export function resolveColorSpaceMode(data: unknown): string {
  if (data && typeof data === 'object') {
    const mode = (data as Record<string, unknown>).colorMode;
    if (isColorSpaceMode(mode)) return mode;
  }
  return DEFAULT_COLOR_SPACE_MODE;
}

/** Propaga a escolha para os dados do nó (mesmo mecanismo das outras props). */
export function applyColorSpaceMode(
  data: unknown,
  mode: string,
): Record<string, unknown> {
  const base = data && typeof data === 'object' ? (data as Record<string, unknown>) : {};
  return {
    ...base,
    colorMode: isColorSpaceMode(mode) ? mode : DEFAULT_COLOR_SPACE_MODE,
  };
}
