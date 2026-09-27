import { describe, expect, it } from 'vitest';
import {
  createRunRequestId,
  isRunResponseFor,
  resolveRunOutcome,
  RUN_CONNECTION_LOST,
} from './runFlowStatus';

describe('correlação pedido/resposta do RUN_FLOW', () => {
  const pending = 'req-123';

  it('aceita a resposta cujo requestId bate com o pedido pendente', () => {
    expect(isRunResponseFor({ status: 'success', requestId: pending }, pending)).toBe(true);
  });

  it('recusa resposta de outra execução (resposta fora de ordem / outro fluxo)', () => {
    expect(isRunResponseFor({ status: 'success', requestId: 'req-antigo' }, pending)).toBe(false);
  });

  it('recusa resposta de preview, que não traz requestId', () => {
    // O autosave manda RUN_FLOW sem id para atualizar thumbnails; ele não pode
    // concluir o pedido que o botão iniciou.
    expect(isRunResponseFor({ status: 'success', thumbnails: { n1: 'data:' } }, pending)).toBe(false);
    expect(isRunResponseFor({ status: 'success', requestId: '' }, pending)).toBe(false);
  });

  it('recusa qualquer resposta quando não há pedido pendente', () => {
    expect(isRunResponseFor({ status: 'error', requestId: pending }, null)).toBe(false);
  });

  it('sucesso do pedido pendente vira done', () => {
    const out = resolveRunOutcome({ status: 'success', requestId: pending }, pending);
    expect(out?.status).toBe('done');
  });

  it('erro do backend vira error e preserva a mensagem real', () => {
    const out = resolveRunOutcome(
      { status: 'error', requestId: pending, error: 'Erro no nó abc: arquivo não encontrado' },
      pending,
    );
    expect(out?.status).toBe('error');
    expect(out?.message).toContain('arquivo não encontrado');
  });

  it('status inesperado do servidor conta como falha, não como concluído', () => {
    const out = resolveRunOutcome({ status: 'algo-novo', requestId: pending }, pending);
    expect(out?.status).toBe('error');
  });

  it('resposta não correlacionada não mexe no status', () => {
    expect(resolveRunOutcome({ status: 'error', requestId: 'outro' }, pending)).toBeNull();
    expect(resolveRunOutcome({ status: 'success' }, pending)).toBeNull();
  });

  it('existe mensagem para execução interrompida por queda de conexão', () => {
    expect(RUN_CONNECTION_LOST).toMatch(/conex[ãa]o/i);
  });

  it('os ids gerados são distintos', () => {
    expect(createRunRequestId()).not.toBe(createRunRequestId());
  });
});
