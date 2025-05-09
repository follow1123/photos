/**
 * @callback ClearRangeFn
 * @param {number} start
 * @param {number} end
 * @returns {void}
 */

/**
 * @typedef {Object} Range
 * @property {number} start 开始下标
 * @property {number} end 结束下标
 */

/**
 * @class
 */
export default class PageRange {
  /** @type {Map<number, Range>} */
  #map;

  /** @type {number | null} */
  #maxPage = null;
  /** @type {number | null} */
  #minPage = null;

  /** @type {number} */
  #maxPageRangeSize = 3;

  /** @type {ClearRangeFn} */
  #clearRange;

  /**
   * @constructor
   * @param {ClearRangeFn} clearRangeFn
   */
  constructor(clearRangeFn) {
    if (!clearRangeFn) throw new Error("clear range function cannot be null");
    this.#map = new Map();
    this.#clearRange = clearRangeFn;
  }

  /**
   * @param {number} num
   * @returns {Range}
   */
  get(num) {
    let range = this.#map.get(num);
    if (!range) throw new Error(`invalid page: ${num}`);
    return range;
  }

  /**
   * @param {number} num
   * @param {Range} range
   */
  set(num, range) {
    //1. 判断是否已经存在
    if (this.#map.has(num)) {
      console.warn(`number: ${num} already set`);
      return;
    }
    //2. 判断是否和最大最小之间有间隔
    if (
      this.#minPage !== null &&
      this.#maxPage !== null &&
      (num > this.#maxPage + 1 || num < this.#minPage - 1)
    ) {
      this.#clearRange(
        this.get(this.#minPage).start,
        this.get(this.#maxPage).end,
      );
      for (let i = this.#minPage; i <= this.#maxPage; i++) {
        this.remove(i);
      }
    }

    //3. 判断已经是最大size
    if (this.#map.size >= this.#maxPageRangeSize) {
      const min = this.minPage;
      const max = this.maxPage;
      if (num < min) {
        this.remove(max);
      } else if (num > max) {
        this.remove(min);
      }
    }

    if (this.#minPage === null || num < this.#minPage) this.#minPage = num;
    if (this.#maxPage === null || num > this.#maxPage) this.#maxPage = num;
    this.#map.set(num, range);
  }

  /**
   * @param {number} num
   */
  remove(num) {
    if (this.#map.size === 0) throw new Error("no element");
    if (!this.#map.has(num)) return;
    if (num === this.minPage) {
      let minRange = this.get(num);
      this.#clearRange(minRange.start, minRange.end);
      this.#minPage = num + 1;
    } else if (num === this.maxPage) {
      let maxRange = this.get(num);
      this.#clearRange(maxRange.start, maxRange.end);
      this.#maxPage = num - 1;
    } else {
      // TODO: 是否需要可以移除中间的页面
      throw new Error("not implemented");
    }

    this.#map.delete(num);
    if (this.#map.size === 0) {
      this.#minPage = null;
      this.#maxPage = null;
    }
  }

  /**
   * @param {number} offset
   */
  shift(offset) {
    this.#map.forEach((r) => {
      r.start += offset;
      r.end += offset;
    });
  }

  /**
   * @param {number} pageSize
   * @param {(start: number, end: number) => void} unload
   * @param {boolean} direction 收缩方向，true: 向前，false: 向后
   */
  resizeDown(pageSize, unload, direction) {
    this.#maxPageRangeSize = pageSize < 67 ? 3 : 2;
    if (direction) {
      let originMaxRangeEnd = this.maxRange.end;
      for (
        let i = this.minPage, s = this.minRange.start;
        i <= this.maxPage;
        i++, s += pageSize
      ) {
        let range = this.get(i);
        range.start = s;
        range.end = s + pageSize - 1;
      }
      unload(this.maxRange.end + 1, originMaxRangeEnd);
    } else {
      let originMinRangeStart = this.minRange.start;
      for (
        let i = this.maxPage, e = this.maxRange.end;
        i >= this.minPage;
        i--, e -= pageSize
      ) {
        let range = this.get(i);
        range.end = e;
        range.start = e - pageSize + 1;
      }
      unload(originMinRangeStart, this.minRange.start - 1);
    }
  }

  /**
   * 由于扩容可能会遇到
   * @param {number} pageSize
   * @param {(pageNum: number, pageSize: number, start: number, end: number) => Promise<number>} load
   * @param {boolean} direction 扩容方向，true: 向后，false: 向前
   * @returns {Promise<void>}
   */
  async resizeUp(pageSize, load, direction) {
    this.#maxPageRangeSize = pageSize < 67 ? 3 : 2;
    let originMinPage = this.minPage;
    let originMaxPage = this.maxPage;
    let originMinRangeStart = this.minRange.start;
    let originMaxRangeEnd = this.maxRange.end;
    this.clear();
    if (direction) {
      let range = {
        start: originMinRangeStart,
        end: originMinRangeStart + pageSize - 1,
      };
      this.set(originMinPage, range);
      return load(originMinPage, pageSize, range.start, range.end).then();
    } else {
      let range = {
        start: originMaxRangeEnd - pageSize + 1,
        end: originMaxRangeEnd,
      };
      this.set(originMaxPage, range);
      return load(originMaxPage, pageSize, range.start, range.end).then();
    }
  }

  clear() {
    this.#map.clear();
    this.#minPage = null;
    this.#maxPage = null;
  }

  get minPage() {
    if (this.#minPage === null) throw new Error("no min page");
    return this.#minPage;
  }

  get maxPage() {
    if (this.#maxPage === null) throw new Error("no max page");
    return this.#maxPage;
  }

  get minRange() {
    return this.get(this.minPage);
  }

  get maxRange() {
    return this.get(this.maxPage);
  }

  get allPagesRange() {
    return /** @type {Range} */ ({
      start: this.minRange.start,
      end: this.maxRange.end,
    });
  }

  get size() {
    return this.#map.size;
  }

  /**
   * @returns {string}
   */
  toString() {
    return JSON.stringify(Array.from(this.#map.entries()));
  }
}
