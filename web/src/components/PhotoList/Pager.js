import FixedSizeInfiniteList from "@components/PhotoList/FixedSizeInfiniteList";
import PageRange from "@components/PhotoList/PageRange";

/**
 * @typedef {import("@components/PhotoList/PageRange").Range} Range
 */

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
export default class Pager {
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

  /** @type {PageRange} */
  #range;

  /**
   * @constructor
   * @param {Options<E>} options
   */
  constructor(options) {
    let opt = options || {};
    if (!opt.elementMgr) throw new Error("element manager cannot be null");
    if (!opt.pageMgr) throw new Error("page manager cannot be null");
    let len = opt.elementLength ?? 200;

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

    this.#range = new PageRange(this.#pageMgr.unload.bind(this.#pageMgr));
  }

  init() {
    this.#list.init();
  }

  /**
   * @param {number} total
   */
  reset(total) {
    this.total = total;
    if (this.#range.size > 0) {
      let allRange = this.#range.allPagesRange;
      this.#pageMgr.unload(allRange.start, allRange.end);
    }
    this.#range.clear();
  }

  /**
   * @param {number} pageSize
   * @returns {Promise<Range>}
   */
  async resize(pageSize) {
    if (this.pageSize === pageSize) throw new Error("same page");

    let direction = this.#range.maxPage !== this.#maxPageNum;
    if (this.pageSize < pageSize) {
      this.#range.resizeDown(
        pageSize,
        this.#pageMgr.unload.bind(this.#pageMgr),
        direction,
      );
    } else {
      await this.#range.resizeUp(
        pageSize,
        (pageNum, pageSize, start, end) => {
          if (start < 0) {
            let movedSize = -start;
            this.#list.scrollUp(movedSize);
            this.#range.shift(-movedSize);
          } else if (end >= this.#list.length) {
            let movedSize = this.#list.length - end;
            this.#list.scrollDown(movedSize);
            this.#range.shift(-movedSize);
          }
          return this.#pageMgr.load(pageNum, pageSize, start, end);
        },
        direction,
      );
    }
    this.pageSize = pageSize;
    console.log("page resize: ", pageSize);
    return this.#range.allPagesRange;
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
   * @returns {Promise<Range>}
   */
  async next() {
    if (this.#range.size === 0) {
      let range = { start: 0, end: this.pageSize - 1 };
      this.#range.set(1, range);
      this.total = await this.#pageMgr.load(
        1,
        this.pageSize,
        range.start,
        range.end,
      );
      return range;
    }

    let maxPageNum = this.#range.maxPage;
    let maxPageRange = this.#range.get(maxPageNum);
    if (maxPageNum === this.#maxPageNum)
      return Promise.reject(new Error("already on last page"));

    let nextNum = maxPageNum + 1;
    let range = {
      start: maxPageRange.end + 1,
      end: maxPageRange.end + this.pageSize,
    };
    this.#range.set(nextNum, range);
    let minPageRange = this.#range.minRange;

    if (range.end >= this.#list.length) {
      let movedSize = Math.min(this.pageSize * 2, minPageRange.start - 1);
      this.#list.scrollDown(movedSize);
      this.#range.shift(-movedSize);
    }

    const total = await this.#pageMgr.load(
      nextNum,
      this.pageSize,
      range.start,
      range.end,
    );
    if (this.total !== total)
      throw new Error("some data changed in this query");
    console.log("next: ", this.#range.toString());
    return { start: minPageRange.start, end: range.end };
  }

  /**
   * @async
   * @returns {Promise<Range>}
   */
  async previous() {
    if (this.#range.size < 2)
      return Promise.reject(new Error("already on last page"));

    let pageSize = this.pageSize;

    let minPageNum = this.#range.minPage;
    let minPageRange = this.#range.get(minPageNum);
    if (minPageNum === 1)
      return Promise.reject(new Error("already on last page"));

    let prevNum = minPageNum - 1;
    let range = {
      start: minPageRange.start - pageSize,
      end: minPageRange.start - 1,
    };
    this.#range.set(prevNum, range);
    let maxPageRange = this.#range.maxRange;

    if (range.start < 0) {
      let movedSize = Math.min(
        pageSize * 2,
        this.#list.length - maxPageRange.end,
      );
      this.#list.scrollUp(movedSize);
      this.#range.shift(movedSize);
    }

    const total = await this.#pageMgr.load(
      prevNum,
      pageSize,
      range.start,
      range.end,
    );
    if (this.total !== total)
      throw new Error("some data changed in this query");
    console.log("previous: ", this.#range.toString());
    return { start: range.start, end: maxPageRange.end };
  }

  /**
   * @returns {boolean}
   */
  hasNext() {
    return this.#maxPageNum ? this.#range.maxPage < this.#maxPageNum : true;
  }

  /**
   * @returns {boolean}
   */
  hasPrevious() {
    return this.#range.minPage > 1;
  }

  get maxRange() {
    return this.#range.maxRange;
  }

  get minRange() {
    return this.#range.minRange;
  }
}
