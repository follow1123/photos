import { expect, describe, it, beforeEach } from "vitest";
import FixedSizeViewer from "@components/PhotoList/fixedSizeViewer";

describe("new FixSizeViewer", () => {
  it("create success", () => {
    expect(
      () =>
        new FixedSizeViewer({
          numOfElements: 128,
          total: 500,
          createElementFn: () => {},
          addElementFn: () => {},
          getElementFn: () => {},
          moveToTopFn: () => {},
          moveToBottomFn: () => {},
        }),
    ).not.toThrowError();
  });
});

//let arr = new Array();
//let arrCounter = 1;
//let viewer = new FixedSizeViewer({
//  numOfElements: 10,
//  max: 30,
//  addElementFn: (element) => arr.push(element),
//  getElementFn: (index) => arr[index],
//  createElementFn: () => arrCounter++,
//  moveToTopFn: (start, end) => arr.unshift(...arr.splice(start, end - start)),
//  moveToBottomFn: (start, end) => arr.push(...arr.splice(start, end - start)),
//});
//
//viewer.init();

describe("move elements", () => {
  let arr = new Array();
  /** @type {number} */
  let arrCounter;
  /** @type {FixedSizeViewer<number>} */
  let viewer;

  beforeEach(() => {
    arr = new Array();
    arrCounter = 1;
    viewer = new FixedSizeViewer({
      numOfElements: 10,
      total: 30,
      addElementFn: (element) => arr.push(element),
      getElementFn: (index) => arr[index],
      createElementFn: () => arrCounter++,
      moveToTopFn: (start, end) =>
        arr.unshift(...arr.splice(start, end - start)),
      moveToBottomFn: (start, end) =>
        arr.push(...arr.splice(start, end - start)),
    });
    viewer.init();
  });

  it("move to bottom", () => {
    viewer.next();
    expect(arr).toEqual([3, 4, 5, 6, 7, 8, 9, 10, 1, 2]);
    viewer.next();
    viewer.next();
    expect(arr).toEqual([7, 8, 9, 10, 1, 2, 3, 4, 5, 6]);
    viewer.next();
    viewer.next();
    viewer.next();
    viewer.next();
    expect(arr).toEqual([5, 6, 7, 8, 9, 10, 1, 2, 3, 4]);
  });

  it("move to top", () => {
    viewer.previous();
    expect(arr).toEqual([9, 10, 1, 2, 3, 4, 5, 6, 7, 8]);
    viewer.previous();
    viewer.previous();
    expect(arr).toEqual([5, 6, 7, 8, 9, 10, 1, 2, 3, 4]);
    viewer.previous();
    viewer.previous();
    viewer.previous();
    viewer.previous();
    expect(arr).toEqual([7, 8, 9, 10, 1, 2, 3, 4, 5, 6]);
  });
});

describe("check view", () => {
  let arr = new Array();
  /** @type {number} */
  let arrCounter;
  /** @type {FixedSizeViewer<number>} */
  let viewer;

  beforeEach(() => {
    arr = new Array();
    arrCounter = 1;
    viewer = new FixedSizeViewer({
      numOfElements: 10,
      total: 30,
      addElementFn: (element) => arr.push(element),
      getElementFn: (index) => arr[index],
      createElementFn: () => arrCounter++,
      moveToTopFn: (start, end) =>
        arr.unshift(...arr.splice(start, end - start)),
      moveToBottomFn: (start, end) =>
        arr.push(...arr.splice(start, end - start)),
    });
    viewer.init();
  });

  it("has next", () => {
    expect(viewer.hasNext()).toBe(true);
    for (let i = 0; i < 10; i++) {
      viewer.next();
    }
    expect(viewer.hasNext()).toBe(false);
  });
  it("has previous", () => {
    expect(viewer.hasPrevious()).toBe(false);
    for (let i = 0; i < 10; i++) {
      viewer.next();
    }
    expect(viewer.hasPrevious()).toBe(true);
  });
});
