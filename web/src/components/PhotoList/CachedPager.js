import FixedSizeInfiniteList from "@components/PhotoList/FixedSizeInfiniteList";

/**
 * @template E
 * @typedef {Object} PagerOptions
 * @property {number} pageSize 每页大小
 * @property {number} total 总共要显示的元素
 * @property {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<E>} manager
 * @property {PageLoader<E>} loader
 */

/**
 * @template E
 * @typedef {Object} PageLoader
 * @property {(pageNum: number, pageSize: number, next: () => E | null) => void} load
 * @property {(next: () => E | null) => void} unload
 */

/**
 * @typedef {Object} PageRange
 * @property {number} start 开始下标
 * @property {number} end 结束下标
 */

/**
 * @template E
 * 分页器
 * @class
 */
export default class CachedPager {
  /** @type {number} */
  total;
  /** @type {number} */
  size;
  /** @type {number} */
  maxNum;

  /** @type {FixedSizeInfiniteList} */
  list;

  /** @type {PageLoader<E>} */
  loader;

  /** @type {Map<number, PageRange>} */
  cache;
  /** @type {number} */
  maxCachePageSize;

  /**
   * @constructor
   * @param {PagerOptions<E>} options
   */
  constructor(options) {
    let opt = options || {};
    if (!opt.total) throw new Error("total not be null");
    if (!opt.pageSize) throw new Error("page size not be null");
    if (!opt.manager) throw new Error("element manager cannot be null");
    if (!opt.loader) throw new Error("page loader cannot be null");
    if (opt.total < opt.pageSize)
      throw new Error("'total' not less then 'page size'");

    this.cache = new Map();
    this.list = new FixedSizeInfiniteList({ len: 128, manager: opt.manager });
    this.loader = opt.loader;
    this.total = opt.total;
    this.size = opt.pageSize;
    this.maxNum = Math.floor((this.total + (this.size - 1)) / this.size);
    this.maxCachePageSize = 3;
  }

  init() {
    this.list.init();
  }

  /**
   * @param {number} pageSize
   * @param {number} total
   */
  reset(pageSize, total) {
    this.total = total;
    this.size = pageSize;
    this.maxNum = Math.floor((this.total + (this.size - 1)) / this.size);
    let maxPageNum = Math.max(...this.cache.keys());
    let maxPageRange = this.getExistsCacheRange(maxPageNum);
    let minPageNum = Math.min(...this.cache.keys());
    let minPageRange = this.getExistsCacheRange(minPageNum);
    let idx = minPageRange.start;
    this.loader.unload(() =>
      idx <= maxPageRange.end ? this.list.manager.getElement(idx) : null,
    );
    this.cache.clear();
  }

  /**
   * @param {number} pageNum
   * @returns {PageRange}
   */
  getExistsCacheRange(pageNum) {
    let range = this.cache.get(pageNum);
    if (!range)
      throw new Error(`get page cache error, page number: ${pageNum}`);
    return range;
  }

  /**
   * @returns {PageRange}
   */
  next() {
    if (this.cache.size === 0) {
      let range = { start: 0, end: this.size - 1 };
      this.cache.set(1, range);
      let idx = range.start;
      this.loader.load(1, this.size, () =>
        idx <= range.end ? this.list.manager.getElement(idx++) : null,
      );
      return range;
    }

    if (this.cache.size === this.maxCachePageSize) {
      let unloadPageNum = Math.min(...this.cache.keys());
      let range = this.getExistsCacheRange(unloadPageNum);
      let idx = range.start;
      this.loader.unload(() =>
        idx <= range.end ? this.list.manager.getElement(idx++) : null,
      );

      this.cache.delete(unloadPageNum);
    }

    let maxPageNum = Math.max(...this.cache.keys());
    let maxPageRange = this.getExistsCacheRange(maxPageNum);
    if (maxPageNum === this.maxNum) {
      throw new Error("already on last page");
    }

    let minPageNum = Math.min(...this.cache.keys());
    let minPageRange = this.getExistsCacheRange(minPageNum);

    let nextNum = maxPageNum + 1;
    let range = {
      start: maxPageRange.end + 1,
      end: maxPageRange.end + this.size,
    };
    this.cache.set(nextNum, range);

    if (range.end >= this.list.length) {
      let movedSize = Math.min(this.size * 2, minPageRange.start - 1);
      this.list.scrollDown(movedSize);
      this.cache.forEach((r) => {
        r.start -= movedSize;
        r.end -= movedSize;
      });
    }

    let idx = range.start;
    this.loader.load(nextNum, this.size, () =>
      idx <= range.end ? this.list.manager.getElement(idx++) : null,
    );

    console.log("next: ", JSON.stringify(Array.from(this.cache.entries())));
    return { start: minPageRange.start, end: range.end };
  }

  previous() {
    if (this.cache.size < 2) throw new Error("already on first page");

    if (this.cache.size === this.maxCachePageSize) {
      let unloadPageNum = Math.max(...this.cache.keys());
      let range = this.getExistsCacheRange(unloadPageNum);
      let idx = range.start;
      this.loader.unload(() =>
        idx <= range.end ? this.list.manager.getElement(idx++) : null,
      );

      this.cache.delete(unloadPageNum);
    }

    let minPageNum = Math.min(...this.cache.keys());
    let minPageRange = this.getExistsCacheRange(minPageNum);
    if (minPageNum === 1) throw new Error("already on first page");

    let maxPageNum = Math.max(...this.cache.keys());
    let maxPageRange = this.getExistsCacheRange(maxPageNum);

    let prevNum = minPageNum - 1;
    let range = {
      start: minPageRange.start - this.size,
      end: minPageRange.start - 1,
    };
    this.cache.set(prevNum, range);

    if (range.start < 0) {
      let movedSize = Math.min(
        this.size * 2,
        this.list.length - maxPageRange.end,
      );
      this.list.scrollUp(movedSize);
      this.cache.forEach((r) => {
        r.start += movedSize;
        r.end += movedSize;
      });
    }

    let idx = range.start;
    this.loader.load(prevNum, this.size, () =>
      idx <= range.end ? this.list.manager.getElement(idx++) : null,
    );

    console.log("previous: ", JSON.stringify(Array.from(this.cache.entries())));
    return { start: range.start, end: maxPageRange.end };
  }

  /**
   * @param {number} num
   * @returns {PageRange}
   */
  set(num) {
    //TODO
    return { start: 0, end: 1 };
  }
}
