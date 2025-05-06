/**
 * @template E
 * @typedef {Object} ViewerOptions
 * @property {number} numOfElements 元素数量
 * @property {number} total 总共要显示的数量
 * @property {() => E} createElementFn 创建元素方法
 * @property {(element: E) => void} addElementFn 添加元素方法
 * @property {(index: number) => E} getElementFn 获取元素方法
 * @property {(start: number, end: number) => void} moveToTopFn 指定元素移动到顶部方法
 * @property {(start: number, end: number) => void} moveToBottomFn 指定元素移动到底部方法
 */

/**
 * @template E
 * @class
 */
export default class FixedSizeViewer {
  /** @type {number} */
  total;
  /** @type {number} */
  counter;
  /** @type {number} */
  numOfElements;
  /** @type {number} */
  numOfMovedElements;

  /** @type {number} */
  nextIdx;
  /** @type {number} */
  prevIdx;

  /** @type {() => E} */
  createElement;

  /** @type {(element: E) => void} */
  addElement;

  /** @type {(index: number) => E} */
  getElement;

  /** @type {(start: number, end: number) => void} */
  moveToTop;

  /** @type {(start: number, end: number) => void} */
  moveToBottom;

  /**
   * @constructor
   * @param {ViewerOptions<E>} options
   */
  constructor(options) {
    let opt = options ?? {};
    if (!opt.numOfElements) throw new Error("number of elements not be null");
    if (!opt.total) throw new Error("total not be null");
    if (!opt.createElementFn)
      throw new Error("create element function not be null");
    if (!opt.addElementFn) throw new Error("add element function not be null");
    if (!opt.getElementFn) throw new Error("get element function not be null");
    if (!opt.moveToTopFn) throw new Error("move to top function not be null");
    if (!opt.moveToBottomFn)
      throw new Error("move to bottom function not be null");
    this.createElement = opt.createElementFn;
    this.addElement = opt.addElementFn;
    this.getElement = opt.getElementFn;
    this.moveToTop = opt.moveToTopFn;
    this.moveToBottom = opt.moveToBottomFn;

    this.total = opt.total;
    this.counter = opt.numOfElements;
    this.numOfElements = opt.numOfElements;
    this.prevIdx = Math.floor(this.numOfElements * 0.25);
    this.nextIdx = Math.floor(this.numOfElements * 0.75);
    this.numOfMovedElements = this.prevIdx;
  }

  init() {
    for (let i = 0; i < this.numOfElements; i++) {
      this.addElement(this.createElement());
    }
  }
  /**
   * @returns {E}
   */
  next() {
    this.moveToBottom(0, this.numOfMovedElements);
    let e = this.getNext();
    this.counter += this.numOfMovedElements;
    return e;
  }
  /**
   * @returns {E}
   */
  previous() {
    this.moveToTop(
      this.numOfElements - this.numOfMovedElements,
      this.numOfElements,
    );
    let e = this.getPrevious();
    this.counter -= this.numOfMovedElements;
    return e;
  }

  /**
   * @returns {boolean}
   */
  hasNext() {
    return this.counter < this.total;
  }

  /**
   * @returns {boolean}
   */
  hasPrevious() {
    return this.counter > this.numOfElements;
  }

  /**
   * @returns {E}
   */
  getNext() {
    return this.getElement(this.nextIdx);
  }

  /**
   * @returns {E}
   */
  getPrevious() {
    return this.getElement(this.prevIdx);
  }
}
