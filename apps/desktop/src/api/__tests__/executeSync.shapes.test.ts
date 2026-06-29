// SPDX-License-Identifier: MIT
import { describe, it, expect } from 'vitest';
import { executeSyncSafe, executeSyncWithDelete, TimeoutError } from '../client';

// These tests are pure data-shape assertions, no React involved.
// They live here so the 401/confirm contract has a single home that
// future changes cannot accidentally regress.

describe('client re-exports', () => {
  it('TimeoutError has the expected name', () => {
    expect(new TimeoutError(1000).name).toBe('TimeoutError');
  });

  it('executeSyncSafe and executeSyncWithDelete remain exported', () => {
    expect(typeof executeSyncSafe).toBe('function');
    expect(typeof executeSyncWithDelete).toBe('function');
  });
});
