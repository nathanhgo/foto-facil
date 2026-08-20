// Default port used by the Go backend when VITE_BACKEND_WS_PORT is not set.
export const DEFAULT_BACKEND_WS_PORT = '8080';

type EnvLike = Record<string, string | undefined>;

/**
 * Builds the backend WebSocket URL, honoring VITE_BACKEND_WS_PORT (loaded by
 * Vite from the repository's root .env) with a safe fallback to the default
 * development port.
 */
export function getBackendWebSocketUrl(env?: EnvLike): string {
  const source = env ?? (import.meta.env as unknown as EnvLike);
  const port = source.VITE_BACKEND_WS_PORT || DEFAULT_BACKEND_WS_PORT;
  return `ws://localhost:${port}/ws`;
}
