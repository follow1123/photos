import { describe, it, expect, beforeEach } from "vitest";
import Pager from "@components/PhotoList/pager";

describe("new pager", () => {
  it("create success", () => {
    expect(
      () =>
        new Pager({
          pageSize: 10,
          total: 100,
        }),
    ).not.toThrowError();
  });
  it("create failure", () => {
    expect(() => new Pager(null)).toThrowError();
    expect(() => new Pager({ total: 10 })).toThrowError();
    expect(() => new Pager({ pageSize: 10 })).toThrowError();
    expect(() => new Pager({ pageSize: 100, total: 10 })).toThrowError();
  });
});

describe("set page number", () => {
  /** @type {Pager} */
  let pager;
  beforeEach(
    () =>
      (pager = new Pager({
        pageSize: 10,
        total: 100,
      })),
  );

  it("set success", () => {
    expect(() => pager.set(1)).not.toThrowError();
    expect(() => pager.set(10)).not.toThrowError();
  });

  it("set failure", () => {
    expect(() => pager.set(0)).toThrowError();
    expect(() => pager.set(11)).toThrowError();
    expect(() => pager.set(-9)).toThrowError();
  });
});

describe("check number", () => {
  /** @type {Pager} */
  let pager;
  beforeEach(
    () =>
      (pager = new Pager({
        pageSize: 10,
        total: 100,
      })),
  );
  it("has next page", () => {
    expect(pager.hasNext()).toBe(true);
    expect(pager.hasNext(9)).toBe(true);
    expect(pager.hasNext(10)).toBe(false);
  });
  it("has previous page", () => {
    expect(pager.hasPrevious()).toBe(false);
    pager.set(8);
    expect(pager.hasPrevious()).toBe(true);
    expect(pager.hasPrevious(1)).toBe(false);
  });
});
