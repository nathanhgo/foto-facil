// Protocol helpers for the RUN_FLOW request/response correlation.
//
// Contract (see architecture-docs/websocket-protocol.md):
//   request : { action: "RUN_FLOW", requestId: string, flow: { nodes, edges } }
//   response: { action: "RUN_FLOW", requestId: string, status: "success"|"error",
//               message?, error?, thumbnails? }
//
// The backend echoes `requestId` verbatim so the UI can tell which request a
// response belongs to. Preview runs triggered by autosave omit `requestId` and
// therefore never resolve a user-initiated run.

export interface RunFlowMessage {
  action: 'RUN_FLOW';
  requestId: string;
  flow: { nodes: unknown[]; edges: unknown[] };
}

export interface BackendRunFlowResponse {
  action?: string;
  status?: string;
  requestId?: string;
  message?: string;
  error?: string;
  thumbnails?: Record<string, string>;
}

/**
 * Generates a unique id for a single user-initiated RUN_FLOW. Prefers the
 * Web Crypto API and falls back to a timestamp/random pair for older runtimes.
 */
export function createRunRequestId(): string {
  const cryptoApi = globalThis.crypto;
  if (cryptoApi && typeof cryptoApi.randomUUID === 'function') {
    return cryptoApi.randomUUID();
  }
  return `run-${Date.now()}-${Math.random().toString(36).slice(2)}`;
}

/**
 * Builds the explicit RUN_FLOW message the "Rodar Fluxo" button sends. A
 * requestId is always attached so the response can be correlated.
 */
export function buildRunFlowMessage(
  requestId: string,
  nodes: unknown[],
  edges: unknown[]
): RunFlowMessage {
  return {
    action: 'RUN_FLOW',
    requestId,
    flow: { nodes, edges },
  };
}

/**
 * True only when `response` is the RUN_FLOW response for the given pending
 * request. Responses from another flow, from a stale request, or from a
 * preview run (no requestId) must never be treated as the current result.
 */
export function isRunFlowResponseFor(
  response: BackendRunFlowResponse | null | undefined,
  requestId: string | null | undefined
): boolean {
  return (
    !!response &&
    response.action === 'RUN_FLOW' &&
    !!requestId &&
    response.requestId === requestId
  );
}
