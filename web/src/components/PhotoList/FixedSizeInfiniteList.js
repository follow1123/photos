/**
 * @template E
 * @typedef {Object} ElementManager
 * @property {() => E} createElement 创建元素方法
 * @property {(element: E) => void} addElement 添加元素方法
 * @property {(index: number) => E} getElement 获取元素方法
 * @property {(start: number, end: number) => void} moveToTop 指定元素移动到顶部方法
 * @property {(start: number, end: number) => void} moveToBottom 指定元素移动到底部方法
 * @exports
 */

/**
 * @typedef {Object} ViewerOptions
 * @property {number} len 元素数量
 * @property {ElementManager<any>} manager 元素管理器
 */

/**
 * @template E
 * @class
 */
export default class FixedSizeInfiniteList {
  /** @type {number} */
  length;

  /** @type {ElementManager<E>} */
  manager;

  /**
   * @constructor
   * @param {ViewerOptions} options
   */
  constructor(options) {
    let opt = options ?? {};
    if (!opt.len) throw new Error("number of elements not be null");
    if (!opt.manager) throw new Error("element manager cannot be null");

    this.manager = opt.manager;
    this.length = opt.len;
  }

  init() {
    for (let i = 0; i < this.length; i++) {
      this.manager.addElement(this.manager.createElement());
    }
  }

  /**
   * @param {number} len
   */
  reset(len) {
    this.length = len;
  }

  /**
   * @param {number} movedSize
   */
  scrollDown(movedSize) {
    this.manager.moveToBottom(0, movedSize);
  }

  /**
   * @param {number} movedSize
   */
  scrollUp(movedSize) {
    this.manager.moveToTop(this.length - movedSize, this.length);
  }
}
