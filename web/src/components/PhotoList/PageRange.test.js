import { expect, describe, it, beforeEach } from "vitest";
import PageRange from "@components/PhotoList/PageRange";

describe("set method", () => {
  /** @type {PageRange} */
  let pr;

  beforeEach(() => {
    pr = new PageRange(() => {});
  });

  it("set exists number", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.set(1, { start: 0, end: 1 });
    pr.set(1, { start: 0, end: 1 });
    expect(pr.size).toBe(1);
    expect(pr.minPage).toBe(1);
    expect(pr.maxPage).toBe(1);
  });

  it("set gap number", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.set(3, { start: 0, end: 1 });
    expect(() => pr.get(1)).toThrowError();
    expect(pr.minPage).toBe(3);
    expect(pr.maxPage).toBe(3);
    pr.set(9, { start: 0, end: 1 });
    expect(() => pr.get(1)).toThrowError();
    expect(() => pr.get(3)).toThrowError();
    expect(pr.minPage).toBe(9);
    expect(pr.maxPage).toBe(9);
    expect(pr.size).toBe(1);
  });

  it("set ordered number", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.set(2, { start: 0, end: 1 });
    pr.set(3, { start: 0, end: 1 });
    expect(pr.minPage).toBe(1);
    expect(pr.maxPage).toBe(3);
    expect(pr.size).toBe(3);
  });

  it("set ordered number to max size", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.set(2, { start: 0, end: 1 });
    pr.set(3, { start: 0, end: 1 });
    pr.set(4, { start: 0, end: 1 });
    expect(pr.minPage).toBe(2);
    expect(pr.maxPage).toBe(4);
    expect(pr.size).toBe(3);
    pr.set(5, { start: 0, end: 1 });
    expect(pr.minPage).toBe(3);
    expect(pr.maxPage).toBe(5);
    expect(pr.size).toBe(3);
    pr.set(2, { start: 0, end: 1 });
    expect(pr.minPage).toBe(2);
    expect(pr.maxPage).toBe(4);
    expect(pr.size).toBe(3);
  });
});

describe("remove method", () => {
  /** @type {PageRange} */
  let pr;

  beforeEach(() => {
    pr = new PageRange(() => {});
  });

  it("remove when empty", () => {
    expect(() => pr.remove(1)).toThrowError();
  });

  it("remove not exists number", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.remove(2);
    expect(pr.size).toBe(1);
  });

  it("remove min", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.set(2, { start: 0, end: 1 });
    pr.set(3, { start: 0, end: 1 });
    pr.remove(1);
    expect(pr.size).toBe(2);
    expect(pr.minPage).toBe(2);
    expect(pr.maxPage).toBe(3);
  });

  it("remove max", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.set(2, { start: 0, end: 1 });
    pr.set(3, { start: 0, end: 1 });
    pr.remove(3);
    expect(pr.size).toBe(2);
    expect(pr.minPage).toBe(1);
    expect(pr.maxPage).toBe(2);
  });

  it("remove when only one element", () => {
    pr.set(1, { start: 0, end: 1 });
    pr.remove(1);
    expect(pr.size).toBe(0);
    expect(() => pr.minPage).toThrowError();
    expect(() => pr.maxPage).toThrowError();
  });
});

describe("shift method", () => {
  /** @type {PageRange} */
  let pr;

  beforeEach(() => {
    pr = new PageRange(() => {});
  });

  it("shift left", () => {
    pr.set(2, { start: 10, end: 19 });
    pr.set(3, { start: 20, end: 29 });
    pr.shift(-5);
    expect(pr.get(2)).toEqual({ start: 5, end: 14 });
    expect(pr.get(3)).toEqual({ start: 15, end: 24 });
  });

  it("shift right", () => {
    pr.set(2, { start: 10, end: 19 });
    pr.set(3, { start: 20, end: 29 });
    pr.shift(7);
    expect(pr.get(2)).toEqual({ start: 17, end: 26 });
    expect(pr.get(3)).toEqual({ start: 27, end: 36 });
  });
});

describe("resize down method", () => {
  /** @type {PageRange} */
  let pr;

  beforeEach(() => {
    pr = new PageRange(() => {});
  });

  it("discard end range 1", () => {
    pr.set(1, { start: 0, end: 9 });
    pr.resizeDown(4, () => {}, true);
    expect(pr.size).toBe(1);
    expect(pr.get(1)).toEqual({ start: 0, end: 3 });
  });

  it("discard end range 2", () => {
    pr.set(1, { start: 0, end: 9 });
    pr.set(2, { start: 10, end: 19 });
    pr.resizeDown(4, () => {}, true);
    expect(pr.size).toBe(2);
    expect(pr.get(1)).toEqual({ start: 0, end: 3 });
    expect(pr.get(2)).toEqual({ start: 4, end: 7 });
  });

  it("discard end range 3", () => {
    pr.set(1, { start: 0, end: 9 });
    pr.set(2, { start: 10, end: 19 });
    pr.set(3, { start: 20, end: 29 });
    pr.resizeDown(4, () => {}, true);
    expect(pr.size).toBe(3);
    expect(pr.get(1)).toEqual({ start: 0, end: 3 });
    expect(pr.get(2)).toEqual({ start: 4, end: 7 });
    expect(pr.get(3)).toEqual({ start: 8, end: 11 });
  });

  it("discard start range 1", () => {
    pr.set(1, { start: 0, end: 9 });
    pr.resizeDown(4, () => {}, false);
    expect(pr.size).toBe(1);
    expect(pr.get(1)).toEqual({ start: 6, end: 9 });
  });

  it("discard start range 2", () => {
    pr.set(1, { start: 0, end: 9 });
    pr.set(2, { start: 10, end: 19 });
    pr.resizeDown(4, () => {}, false);
    expect(pr.size).toBe(2);
    expect(pr.get(1)).toEqual({ start: 12, end: 15 });
    expect(pr.get(2)).toEqual({ start: 16, end: 19 });
  });

  it("discard start range 3", () => {
    pr.set(1, { start: 0, end: 9 });
    pr.set(2, { start: 10, end: 19 });
    pr.set(3, { start: 20, end: 29 });
    pr.resizeDown(4, () => {}, false);
    expect(pr.size).toBe(3);
    expect(pr.get(1)).toEqual({ start: 18, end: 21 });
    expect(pr.get(2)).toEqual({ start: 22, end: 25 });
    expect(pr.get(3)).toEqual({ start: 26, end: 29 });
  });
});

describe("resize up method", () => {
  /** @type {PageRange} */
  let pr;

  beforeEach(() => {
    pr = new PageRange(() => {});
  });

  it("expand end range 1", async () => {
    pr.set(1, { start: 0, end: 3 });
    await pr.resizeUp(10, () => Promise.resolve(1), true);
    expect(pr.size).toBe(1);
    expect(pr.get(1)).toEqual({ start: 0, end: 9 });
  });

  it("expand end range 2", async () => {
    pr.set(1, { start: 0, end: 3 });
    pr.set(2, { start: 4, end: 7 });
    await pr.resizeUp(10, () => Promise.resolve(1), true);
    expect(pr.size).toBe(1);
    expect(pr.get(1)).toEqual({ start: 0, end: 9 });
  });

  it("expand end range 3", async () => {
    pr.set(1, { start: 0, end: 3 });
    pr.set(2, { start: 4, end: 7 });
    pr.set(3, { start: 8, end: 11 });
    await pr.resizeUp(10, () => Promise.resolve(1), true);
    expect(pr.size).toBe(1);
    expect(pr.get(1)).toEqual({ start: 0, end: 9 });
  });

  it("expand start range 1", async () => {
    pr.set(3, { start: 8, end: 11 });
    await pr.resizeUp(10, () => Promise.resolve(1), false);
    expect(pr.size).toBe(1);
    expect(pr.get(3)).toEqual({ start: 2, end: 11 });
  });

  it("expand start range 2", async () => {
    pr.set(2, { start: 4, end: 7 });
    pr.set(3, { start: 8, end: 11 });
    await pr.resizeUp(10, () => Promise.resolve(1), false);
    expect(pr.size).toBe(1);
    expect(pr.get(3)).toEqual({ start: 2, end: 11 });
  });

  it("expand start range 3", async () => {
    pr.set(1, { start: 0, end: 3 });
    pr.set(2, { start: 4, end: 7 });
    pr.set(3, { start: 8, end: 11 });
    await pr.resizeUp(10, () => Promise.resolve(1), false);
    expect(pr.size).toBe(1);
    expect(pr.get(3)).toEqual({ start: 2, end: 11 });
  });
});
