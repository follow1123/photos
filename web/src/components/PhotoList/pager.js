/**
 * @typedef {Object} PagerOptions
 * @property {number} pageSize 每页大小
 * @property {number} total 总共要显示的元素
 */

/**
 * 分页器
 * @class
 */
export default class Pager {
  /** @type {number} */
  total;
  /** @type {number} */
  num = 0;
  /** @type {number} */
  size;
  /** @type {number} */
  maxNumber;

  /**
   * @constructor
   * @param {PagerOptions} options
   */
  constructor(options) {
    let opt = options || {};
    if (!opt.total) throw new Error("total not be null");
    if (!opt.pageSize) throw new Error("page size not be null");
    if (opt.total < opt.pageSize)
      throw new Error("'total' not less then 'page size'");

    this.total = opt.total;
    this.size = opt.pageSize;
    this.maxNumber = Math.floor((this.total + (this.size - 1)) / this.size);
  }

  next() {
    this.set(this.getNext());
  }

  previous() {
    this.set(this.getPrevious());
  }

  /**
   * @param {number} [pageNum]
   * @returns {boolean}
   */
  hasNext(pageNum) {
    return pageNum ? pageNum < this.maxNumber : this.num < this.maxNumber;
  }

  /**
   * @param {number} [pageNum]
   * @returns {boolean}
   */
  hasPrevious(pageNum) {
    return pageNum ? pageNum > 1 : this.num > 1;
  }

  /**
   * @returns {number}
   */
  getNext() {
    return this.num + 1;
  }

  /**
   * @returns {number}
   */
  getPrevious() {
    return this.num - 1;
  }

  /**
   * @param {number} num
   */
  set(num) {
    if (num <= 0 || num > this.maxNumber)
      throw new Error(`error page number: ${num}`);
    this.num = num;
  }
}
