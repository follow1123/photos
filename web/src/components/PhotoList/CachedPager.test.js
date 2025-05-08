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
 * @returns {import("@components/PhotoList/CachedPager").PageManager<number>}
 */
function NumPageManager() {
  /** @type {import("@components/PhotoList/CachedPager").PageManager<number>} */
  let pl = {
    load: (pageNum, pageSize, next) => Promise.resolve(100),
    unload: (e) => {},
    show: (e) => {},
    hide: (e) => {},
  };
  return pl;
}

describe("next page", () => {
  /** @type {Array<number>} */
  let arr;
  /** @type {CachedPager<number>} */
  let pager;

  beforeEach(() => {
    arr = new Array();
    pager = new CachedPager({
      elementLength: 62,
      elementMgr: ArrayManager(arr),
      pageMgr: NumPageManager(),
    });
    pager.init();
    pager.pageSize = 10;
  });

  it("next", async () => {
    await expect(pager.next()).resolves.toEqual({ start: 0, end: 9 });
    await expect(pager.next()).resolves.toEqual({ start: 0, end: 19 });
    await expect(pager.next()).resolves.toEqual({ start: 0, end: 29 });
    await expect(pager.next()).resolves.toEqual({ start: 10, end: 39 });
    await expect(pager.next()).resolves.toEqual({ start: 20, end: 49 });
  });

  it("scroll triggered", async () => {
    for (let i = 0; i < 6; i++) {
      await pager.next();
    }
    await expect(pager.next()).resolves.toEqual({ start: 20, end: 49 });
    await expect(pager.next()).resolves.toEqual({ start: 30, end: 59 });
  });

  it("to bottom", async () => {
    for (let i = 0; i < 10; i++) {
      await pager.next();
    }
    await expect(() => pager.next()).rejects.toThrowError();
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
      elementLength: 62,
      elementMgr: ArrayManager(arr),
      pageMgr: NumPageManager(),
    });
    pager.init();
    pager.pageSize = 10;
  });

  it("previous", async () => {
    for (let i = 0; i < 10; i++) {
      await pager.next();
    }

    await expect(pager.previous()).resolves.toEqual({ start: 20, end: 49 });
    await expect(pager.previous()).resolves.toEqual({ start: 10, end: 39 });
    await expect(pager.previous()).resolves.toEqual({ start: 0, end: 29 });
  });

  it("scroll triggered", async () => {
    for (let i = 0; i < 10; i++) {
      await pager.next();
    }
    await expect(pager.previous()).resolves.toEqual({ start: 20, end: 49 });
    await expect(pager.previous()).resolves.toEqual({ start: 10, end: 39 });
    await expect(pager.previous()).resolves.toEqual({ start: 0, end: 29 });
    await expect(pager.previous()).resolves.toEqual({ start: 10, end: 39 });
  });

  it("to top", async () => {
    await expect(() => pager.previous()).rejects.toThrowError();
    for (let i = 0; i < 10; i++) {
      await pager.next();
    }

    for (let i = 0; i < 7; i++) {
      await pager.previous();
    }
    await expect(() => pager.previous()).rejects.toThrowError();
  });
});
