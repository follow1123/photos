import FixedSizeInfiniteList from "@components/PhotoList/FixedSizeInfiniteList";

/**
 * @template E
 * @callback LoadPageFn
 * @param {number} pageNum
 * @param {number} pageSize
 * @param {() => E | null} next
 * @returns {Promise<number>}
 */

/**
 * @template E
 * @callback ElementConsumer
 * @param {E} element
 * @returns {void}
 */

/**
 * @template E
 * @typedef {Object} PageManager
 * @property {LoadPageFn<E>} load
 * @property {ElementConsumer<E>} unload
 * @property {ElementConsumer<E>} show 测试
 * @property {ElementConsumer<E>} hide
 */

/**
 * @template E
 * @typedef {Object} Options
 * @property {number} [elementLength] 总共要显示的元素
 * @property {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<E>} elementMgr
 * @property {PageManager<E>} pageMgr
 */

/**
 * @typedef {Object} PageRange
 * @property {number} start 开始下标
 * @property {number} end 结束下标
 */

class Cache {
  /** @type {Map<number, PageRange>} */
  map;

  /** @type {number} */
  #maxPage = 0;
  /** @type {number} */
  #minPage = 0;

  constructor() {
    this.map = new Map();
  }

  /**
   * @returns {number}
   */
  getMinPage() {
    if (this.#minPage === 0) throw new Error("no min page");
    return this.#minPage;
  }

  /**
   * @returns {number}
   */
  getMaxPage() {
    if (this.#maxPage === 0) throw new Error("no max page");
    return this.#maxPage;
  }

  /**
   * @param {number} num
   */
  remove(num) {
    if (this.map.size === 0) throw new Error("no element");
    this.map.delete(num);
    this.#minPage = Math.min(...this.map.keys());
    this.#maxPage = Math.max(...this.map.keys());
  }

  /**
   * @returns {number}
   */
  size() {
    return this.map.size;
  }

  clear() {
    this.map.clear();
    this.#minPage = 0;
    this.#maxPage = 0;
  }

  /**
   * @param {number} num
   * @returns {PageRange}
   */
  get(num) {
    let range = this.map.get(num);
    if (!range) throw new Error(`invalid page: ${this.#maxPage}`);
    return range;
  }

  /**
   * @param {number} num
   * @param {PageRange} range
   */
  set(num, range) {
    if (this.#minPage === 0 || num < this.#minPage) this.#minPage = num;
    if (this.#maxPage === 0 || num > this.#maxPage) this.#maxPage = num;
    this.map.set(num, range);
  }
}

/**
 * @template E
 * @class
 */
class PageMgrWrapper {
  /** @type {number} */
  #elementLength;
  /** @type {() => number} */
  #getMaxPageNum;

  /** @type {PageManager<E>} */
  #pageMgr;
  /** @type {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<E>} */
  #elementMgr;

  /**
   * @constructor
   * @param {PageManager<E>} pageMgr
   * @param {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<E>} elementMgr
   * @param {number} elementLen
   * @param {() => number} maxPageNumFn
   */
  constructor(pageMgr, elementMgr, maxPageNumFn, elementLen) {
    this.#pageMgr = pageMgr;
    this.#elementMgr = elementMgr;
    this.#elementLength = elementLen;
    this.#getMaxPageNum = maxPageNumFn;
  }

  /**
   * @async
   * @param {number} pageNum
   * @param {number} pageSize
   * @param {number} start
   * @param {number} end
   * @returns {Promise<number>}
   */
  async load(pageNum, pageSize, start, end) {
    let idx = start;
    let next = () => {
      if (idx > end) return null;
      let e = this.#elementMgr.getElement(idx++);
      this.#pageMgr.show(e);
      return e;
    };

    //this.#pageMgr.load(
    const total = await this.#pageMgr.load(pageNum, pageSize, next);
    if (idx <= end || pageNum === this.#getMaxPageNum()) {
      for (let i = idx; i < this.#elementLength; i++) {
        this.#pageMgr.hide(this.#elementMgr.getElement(i));
      }
    }
    return total;
  }

  /**
   * @param {number} start
   * @param {number} end
   */
  unload(start, end) {
    for (let i = start; i <= end; i++) {
      let e = this.#elementMgr.getElement(i);
      this.#pageMgr.show(e);
      this.#pageMgr.unload(e);
    }
  }

  /**
   * @param {number} start
   * @param {number} end
   */
  hide(start, end) {
    for (let i = start; i <= end; i++) {
      this.#pageMgr.hide(this.#elementMgr.getElement(i));
    }
  }
}

/**
 * 分页器
 * @template E
 * @class
 */
export default class CachedPager {
  /** @type {number | null} */
  #total = null;
  /** @type {number | null} */
  #pageSize = null;
  /** @type {number | null} */
  #maxPageNum = null;

  /** @type {FixedSizeInfiniteList<E>} */
  #list;

  /** @type {PageMgrWrapper<E>} */
  #pageMgr;

  /** @type {Cache} */
  #cache;

  /** @type {number} */
  #defaultMaxCacheSize = 3;
  /** @type {number} */
  #maxCacheSize;

  /**
   * @constructor
   * @param {Options<E>} options
   */
  constructor(options) {
    let opt = options || {};
    if (!opt.elementMgr) throw new Error("element manager cannot be null");
    if (!opt.pageMgr) throw new Error("page manager cannot be null");
    let len = opt.elementLength ?? 128;

    this.#cache = new Cache();
    this.#maxCacheSize = this.#defaultMaxCacheSize;

    this.#list = new FixedSizeInfiniteList({
      len: len,
      manager: opt.elementMgr,
    });

    this.#pageMgr = new PageMgrWrapper(
      opt.pageMgr,
      opt.elementMgr,
      () => this.#maxPageNum ?? -1,
      len,
    );
  }

  init() {
    this.#list.init();
  }

  /**
   * @param {number} total
   */
  reset(total) {
    this.total = total;
    if (this.#cache.size() > 0) {
      this.#pageMgr.unload(this.minPageRange.start, this.maxPageRange.end);
    }
    this.#cache.clear();
  }

  /**
   * @param {number} pageSize
   */
  resize(pageSize) {
    if (this.pageSize === pageSize) {
      return;
    }
    // TODO: 重置逻辑
    console.log("page resize: ", pageSize);
  }

  get total() {
    if (!this.#total) throw new Error("total not initialized");
    return this.#total;
  }

  set total(value) {
    this.#maxPageNum = Math.floor(
      (value + (this.pageSize - 1)) / this.pageSize,
    );
    this.#total = value;
  }

  get pageSize() {
    if (!this.#pageSize) throw new Error("page size not initialized");
    return this.#pageSize;
  }

  set pageSize(value) {
    if (this.#total) {
      this.#maxPageNum = Math.floor((this.#total + (value - 1)) / value);
    }
    this.#pageSize = value;
  }

  /**
   * @async
   * @returns {Promise<PageRange>}
   */
  async next() {
    if (this.#cache.map.size === 0) {
      let range = { start: 0, end: this.pageSize - 1 };
      this.#cache.set(1, range);
      this.total = await this.#pageMgr.load(
        1,
        this.pageSize,
        range.start,
        range.end,
      );
      return range;
    }

    if (this.#cache.size() === this.#maxCacheSize) {
      let unloadPageNum = this.#cache.getMinPage();
      let range = this.#cache.get(unloadPageNum);
      this.#pageMgr.unload(range.start, range.end);

      this.#cache.remove(unloadPageNum);
    }

    let maxPageNum = this.#cache.getMaxPage();
    let maxPageRange = this.#cache.get(maxPageNum);
    if (maxPageNum === this.#maxPageNum)
      return Promise.reject(new Error("already on last page"));

    let minPageRange = this.minPageRange;

    let nextNum = maxPageNum + 1;
    let range = {
      start: maxPageRange.end + 1,
      end: maxPageRange.end + this.pageSize,
    };
    this.#cache.set(nextNum, range);
    if (range.end >= this.#list.length) {
      let movedSize = Math.min(this.pageSize * 2, minPageRange.start - 1);
      this.#list.scrollDown(movedSize);
      this.#cache.map.forEach((r) => {
        r.start -= movedSize;
        r.end -= movedSize;
      });
    }

    const total = await this.#pageMgr.load(
      nextNum,
      this.pageSize,
      range.start,
      range.end,
    );
    if (this.total !== total)
      throw new Error("some data changed in this query");
    console.log(
      "next: ",
      JSON.stringify(Array.from(this.#cache.map.entries())),
    );
    return { start: minPageRange.start, end: range.end };
  }

  /**
   * @async
   * @returns {Promise<PageRange>}
   */
  async previous() {
    if (this.#cache.size() < 2)
      return Promise.reject(new Error("already on last page"));

    let pageSize = this.pageSize;
    if (this.#cache.size() === this.#maxCacheSize) {
      let unloadPageNum = this.#cache.getMaxPage();
      let range = this.#cache.get(unloadPageNum);
      this.#pageMgr.unload(range.start, range.end);

      this.#cache.remove(unloadPageNum);
    }

    let minPageNum = this.#cache.getMinPage();
    let minPageRange = this.#cache.get(minPageNum);
    if (minPageNum === 1)
      return Promise.reject(new Error("already on last page"));

    let maxPageNum = this.#cache.getMaxPage();
    let maxPageRange = this.#cache.get(maxPageNum);

    let prevNum = minPageNum - 1;
    let range = {
      start: minPageRange.start - pageSize,
      end: minPageRange.start - 1,
    };
    this.#cache.set(prevNum, range);

    if (range.start < 0) {
      let movedSize = Math.min(
        pageSize * 2,
        this.#list.length - maxPageRange.end,
      );
      this.#list.scrollUp(movedSize);
      this.#cache.map.forEach((r) => {
        r.start += movedSize;
        r.end += movedSize;
      });
    }

    const total = await this.#pageMgr.load(
      prevNum,
      pageSize,
      range.start,
      range.end,
    );
    if (this.total !== total)
      throw new Error("some data changed in this query");
    console.log(
      "next: ",
      JSON.stringify(Array.from(this.#cache.map.entries())),
    );
    return { start: range.start, end: maxPageRange.end };
  }

  /**
   * @returns {boolean}
   */
  hasNext() {
    return this.#maxPageNum
      ? this.#cache.getMaxPage() < this.#maxPageNum
      : true;
  }

  /**
   * @returns {boolean}
   */
  hasPrevious() {
    return this.#cache.getMinPage() > 1;
  }

  /**
   * @returns {PageRange}
   */
  get maxPageRange() {
    return this.#cache.get(this.#cache.getMaxPage());
  }

  /**
   * @returns {PageRange}
   */
  get minPageRange() {
    return this.#cache.get(this.#cache.getMinPage());
  }
}
