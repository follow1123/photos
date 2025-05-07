import { expect, describe, it, beforeEach } from "vitest";
import CachedPager from "@components/PhotoList/CachedPager";

/**
 * @param {Array<number>} arr
 * @returns {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<number>}
 */
function ArrayManager(arr) {
  let arrCounter = 1;
  /** @type {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<number>} */
  let m = {
    addElement: (element) => arr.push(element),
    getElement: (index) => arr[index],
    createElement: () => arrCounter++,
    moveToTop: (start, end) => arr.unshift(...arr.splice(start, end - start)),
    moveToBottom: (start, end) => arr.push(...arr.splice(start, end - start)),
  };
  return m;
}

/**
 * @returns {import("@components/PhotoList/CachedPager").PageLoader<number>}
 */
function NumPageLoader() {
  /** @type {import("@components/PhotoList/CachedPager").PageLoader<number>} */
  let pl = {
    load: (pageNum, pageSize, next) => {},
    unload: (elements) => {},
  };
  return pl;
}

//describe("new FixSizeViewer", () => {
//  it("create success", () => {
//    expect(
//      () =>
//        new FixedSizeInfiniteList({
//          len: 128,
//          manager: ArrayManager(new Array()),
//        }),
//    ).not.toThrowError();
//  });
//});

describe("next page", () => {
  /** @type {Array<number>} */
  let arr;
  /** @type {CachedPager<number>} */
  let pager;

  beforeEach(() => {
    arr = new Array();
    pager = new CachedPager({
      total: 300,
      pageSize: 10,
      manager: ArrayManager(arr),
      loader: NumPageLoader(),
    });
    pager.init();
  });

  it("next", () => {
    expect(pager.next()).toEqual({ start: 0, end: 9 });
    expect(pager.next()).toEqual({ start: 0, end: 19 });
    expect(pager.next()).toEqual({ start: 0, end: 29 });
    expect(pager.next()).toEqual({ start: 10, end: 39 });
    expect(pager.next()).toEqual({ start: 20, end: 49 });
    expect(pager.next()).toEqual({ start: 30, end: 59 });
  });

  it("scroll triggered", () => {
    for (let i = 0; i < 11; i++) {
      pager.next();
    }
    expect(pager.next()).toEqual({ start: 90, end: 119 });
    expect(pager.next()).toEqual({ start: 80, end: 109 });
  });

  it("to bottom", () => {
    for (let i = 0; i < 30; i++) {
      pager.next();
    }
    expect(() => pager.next()).toThrowError();
  });
});

describe("previous page", () => {
  /** @type {Array<number>} */
  let arr;
  /** @type {CachedPager<number>} */
  let pager;

  beforeEach(() => {
    arr = new Array();
    pager = new CachedPager({
      total: 300,
      pageSize: 10,
      manager: ArrayManager(arr),
      loader: NumPageLoader(),
    });
    pager.init();
  });

  it("previous", () => {
    for (let i = 0; i < 10; i++) {
      pager.next();
    }

    expect(pager.previous()).toEqual({ start: 60, end: 89 });
    expect(pager.previous()).toEqual({ start: 50, end: 79 });
    expect(pager.previous()).toEqual({ start: 40, end: 69 });
  });

  it("scroll triggered", () => {
    for (let i = 0; i < 12; i++) {
      pager.next();
    }
    expect(pager.previous()).toEqual({ start: 80, end: 109 });
  });

  it("to top", () => {
    expect(() => pager.previous()).toThrowError();
    for (let i = 0; i < 30; i++) {
      pager.next();
    }

    for (let i = 0; i < 27; i++) {
      pager.previous();
    }
    expect(() => pager.previous()).toThrowError();
  });
});
