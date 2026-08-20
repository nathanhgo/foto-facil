import { describe, expect, it } from 'vitest';
import { DEFAULT_BACKEND_WS_PORT, getBackendWebSocketUrl } from './websocketUrl';

describe('getBackendWebSocketUrl', () => {
  it('falls back to the default port when no env var is provided', () => {
    expect(getBackendWebSocketUrl({})).toBe(`ws://localhost:${DEFAULT_BACKEND_WS_PORT}/ws`);
  });

  it('uses VITE_BACKEND_WS_PORT when it is set', () => {
    expect(getBackendWebSocketUrl({ VITE_BACKEND_WS_PORT: '9090' })).toBe('ws://localhost:9090/ws');
  });

  it('ignores an empty string and falls back to the default', () => {
    expect(getBackendWebSocketUrl({ VITE_BACKEND_WS_PORT: '' })).toBe(`ws://localhost:${DEFAULT_BACKEND_WS_PORT}/ws`);
  });
});
