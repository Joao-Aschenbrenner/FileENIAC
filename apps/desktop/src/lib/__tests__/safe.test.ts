import { describe, it, expect } from 'vitest';
import { asArray, hasItems, safeString, safeNumber } from '../safe';

describe('asArray', () => {
  it('returns the array when given an array', () => {
    expect(asArray([1, 2, 3])).toEqual([1, 2, 3]);
  });
  it('returns empty array for null', () => {
    expect(asArray(null)).toEqual([]);
  });
  it('returns empty array for undefined', () => {
    expect(asArray(undefined)).toEqual([]);
  });
  it('returns empty array for non-array objects', () => {
    expect(asArray({} as any)).toEqual([]);
  });
});

describe('hasItems', () => {
  it('returns true for non-empty array', () => {
    expect(hasItems([1])).toBe(true);
  });
  it('returns false for empty array', () => {
    expect(hasItems([])).toBe(false);
  });
  it('returns false for null', () => {
    expect(hasItems(null)).toBe(false);
  });
  it('returns false for undefined', () => {
    expect(hasItems(undefined)).toBe(false);
  });
});

describe('safeString', () => {
  it('returns the string when given a string', () => {
    expect(safeString('hello')).toBe('hello');
  });
  it('returns fallback for null', () => {
    expect(safeString(null)).toBe('');
  });
  it('returns fallback for undefined', () => {
    expect(safeString(undefined)).toBe('');
  });
  it('uses custom fallback', () => {
    expect(safeString(null, 'default')).toBe('default');
  });
});

describe('safeNumber', () => {
  it('returns the number when valid', () => {
    expect(safeNumber(42)).toBe(42);
  });
  it('returns fallback for null', () => {
    expect(safeNumber(null)).toBe(0);
  });
  it('returns fallback for NaN', () => {
    expect(safeNumber(NaN)).toBe(0);
  });
  it('uses custom fallback', () => {
    expect(safeNumber(null, -1)).toBe(-1);
  });
});
