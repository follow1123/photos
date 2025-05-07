import { expect, describe, it, beforeEach } from "vitest";
import FixedSizeInfiniteList from "@components/PhotoList/FixedSizeInfiniteList";

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

describe("new FixSizeViewer", () => {
  it("create success", () => {
    expect(
      () =>
        new FixedSizeInfiniteList({
          len: 128,
          manager: ArrayManager(new Array()),
        }),
    ).not.toThrowError();
  });
});

describe("scroll", () => {
  /** @type {Array<number>} */
  let arr;
  /** @type {FixedSizeInfiniteList} */
  let list;

  beforeEach(() => {
    arr = new Array();
    list = new FixedSizeInfiniteList({
      len: 10,
      manager: ArrayManager(arr),
    });
    list.init();
  });

  it("scroll down", () => {
    list.scrollDown(2);
    expect(arr).toEqual([3, 4, 5, 6, 7, 8, 9, 10, 1, 2]);
    list.scrollDown(4);
    expect(arr).toEqual([7, 8, 9, 10, 1, 2, 3, 4, 5, 6]);
    list.scrollDown(8);
    expect(arr).toEqual([5, 6, 7, 8, 9, 10, 1, 2, 3, 4]);
  });

  it("move to top", () => {
    list.scrollUp(2);
    expect(arr).toEqual([9, 10, 1, 2, 3, 4, 5, 6, 7, 8]);
    list.scrollUp(4);
    expect(arr).toEqual([5, 6, 7, 8, 9, 10, 1, 2, 3, 4]);
    list.scrollUp(8);
    expect(arr).toEqual([7, 8, 9, 10, 1, 2, 3, 4, 5, 6]);
    list.scrollUp(1);
    expect(arr).toEqual([6, 7, 8, 9, 10, 1, 2, 3, 4, 5]);
  });
});
