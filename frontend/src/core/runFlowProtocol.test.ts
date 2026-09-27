import { describe, expect, it } from 'vitest';
import {
  buildRunFlowMessage,
  createRunRequestId,
  isRunFlowResponseFor,
  type BackendRunFlowResponse,
} from './runFlowProtocol';

describe('buildRunFlowMessage', () => {
  it('builds an explicit RUN_FLOW action carrying the requestId and flow', () => {
    const nodes = [{ id: 'n1' }];
    const edges = [{ source: 'n1', target: 'n2' }];

    const message = buildRunFlowMessage('req-1', nodes, edges);

    expect(message.action).toBe('RUN_FLOW');
    expect(message.requestId).toBe('req-1');
    expect(message.flow).toEqual({ nodes, edges });
  });
});

describe('isRunFlowResponseFor', () => {
  const matching: BackendRunFlowResponse = {
    action: 'RUN_FLOW',
    requestId: 'req-1',
    status: 'success',
  };

  it('accepts the response that echoes the pending requestId', () => {
    expect(isRunFlowResponseFor(matching, 'req-1')).toBe(true);
  });

  it('rejects a stale response from another request', () => {
    expect(isRunFlowResponseFor({ ...matching, requestId: 'req-old' }, 'req-1')).toBe(false);
  });

  it('rejects responses without the RUN_FLOW action', () => {
    expect(isRunFlowResponseFor({ ...matching, action: undefined }, 'req-1')).toBe(false);
  });

  it('rejects preview responses that carry no requestId', () => {
    expect(isRunFlowResponseFor({ action: 'RUN_FLOW', status: 'success' }, 'req-1')).toBe(false);
  });

  it('rejects when there is no pending request', () => {
    expect(isRunFlowResponseFor(matching, null)).toBe(false);
  });

  it('rejects empty payloads', () => {
    expect(isRunFlowResponseFor(undefined, 'req-1')).toBe(false);
    expect(isRunFlowResponseFor(null, 'req-1')).toBe(false);
  });
});

describe('createRunRequestId', () => {
  it('returns a non-empty unique identifier', () => {
    const first = createRunRequestId();
    const second = createRunRequestId();

    expect(first).toBeTruthy();
    expect(second).toBeTruthy();
    expect(first).not.toBe(second);
  });
});
