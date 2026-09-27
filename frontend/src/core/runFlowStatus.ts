/**
 * Correlação entre o pedido de execução (RUN_FLOW) e a resposta do backend.
 *
 * O botão "Rodar Fluxo" manda um `requestId` junto do RUN_FLOW; o backend ecoa
 * esse id em toda resposta daquela execução (sucesso, erro de nó, ciclo no
 * grafo). Só a resposta com o id que está pendente pode mudar o status da UI —
 * resposta de preview (autosave, sem id) ou execução anterior (id diferente)
 * atualiza no máximo os thumbnails, nunca o status.
 *
 * Isto existe porque a versão anterior marcava `done` com um `setTimeout` fixo,
 * sem ouvir o backend: a UI dizia "✓ Concluído" sem nada ter rodado, e o estado
 * `error` era código morto.
 */

export type RunStatus = 'idle' | 'running' | 'done' | 'error';

/** Resposta vinda do WebSocket — só os campos que interessam à correlação. */
export interface RunFlowResponseLike {
  status?: string;
  error?: string;
  message?: string;
  requestId?: string;
  thumbnails?: Record<string, string>;
}

export interface RunOutcome {
  status: Extract<RunStatus, 'done' | 'error'>;
  /** Mensagem para mostrar ao usuário quando o resultado é `error`. */
  message: string;
}

/** Gera o id de uma execução. `crypto.randomUUID` existe no Electron e nos
 *  navegadores atuais; o fallback cobre ambiente de teste antigo. */
export function createRunRequestId(): string {
  const webCrypto = typeof globalThis !== 'undefined' ? globalThis.crypto : undefined;
  if (webCrypto && typeof webCrypto.randomUUID === 'function') {
    return webCrypto.randomUUID();
  }
  return `run-${Date.now()}-${Math.floor(Math.random() * 1e6)}`;
}

/**
 * A resposta pertence ao pedido pendente?
 *
 * `false` quando não há pedido pendente, quando a resposta não traz `requestId`
 * (é de preview) ou quando o id é de outra execução.
 */
export function isRunResponseFor(
  response: RunFlowResponseLike | null | undefined,
  pendingRequestId: string | null,
): boolean {
  if (!response || !pendingRequestId) return false;
  const id = response.requestId;
  return typeof id === 'string' && id !== '' && id === pendingRequestId;
}

/**
 * Traduz a resposta correlacionada no próximo status da UI.
 *
 * Devolve `null` quando a resposta não é da execução pendente — nesse caso o
 * chamador não deve tocar no status. Um erro do backend traz a mensagem real
 * (ex.: "Erro no nó X: ..."); qualquer `status` diferente de `success` conta
 * como falha, para não engolir resposta inesperada do servidor.
 */
export function resolveRunOutcome(
  response: RunFlowResponseLike | null | undefined,
  pendingRequestId: string | null,
): RunOutcome | null {
  if (!isRunResponseFor(response, pendingRequestId)) return null;

  if (response!.status === 'success') {
    return { status: 'done', message: response!.message || '' };
  }
  return {
    status: 'error',
    message: response!.error || 'Falha ao executar o fluxo.',
  };
}

/** Mensagem para quando a conexão cai com uma execução ainda pendente. */
export const RUN_CONNECTION_LOST = 'Conexão com o backend foi perdida durante a execução.';
