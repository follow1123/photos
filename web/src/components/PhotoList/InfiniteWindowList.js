/**
 * @template E
 * @callback QueryElementsFn
 * @param {number} pageNum
 * @param {number} pageSize
 * @param {() => E} next
 * @returns {Promise<number>}
 */

/**
 * @template E
 * @typedef {Object} ElementManager
 * @property {QueryElementsFn<E>} queryElement
 * @property {() => E} createElement
 * @property {(element: E) => void} load
 * @property {(element: E) => void} unload
 */

/**
 * @template E
 * @typedef {Object} Options
 * @property {Element} root 窗口元素
 * @property {number | undefined} [pageSize] 分页大小
 * @property {ElementManager<E>} manager 元素管理器
 */

/**
 * 无限窗口列表
 * @template E
 * @class
 */
export default class InfiniteWindowList {
  /** @type {Element} */
  #root;

  /** @type {number} */
  #counter = 0;
  /** @type {number} */
  #pageSize;
  /** @type {number} */
  #pageNum = 0;
  /** @type {number | null} */
  #total = null;

  /** @type {IntersectionObserver} */
  #observer;
  /** @type {Map<string, HTMLElement>} */
  #observedMap;

  /** @type {ElementManager<E>} */
  #manager;

  /**
   * @constructor
   * @param {Options<E>} options
   */
  constructor(options) {
    // 参数校验
    let opt = options ?? {};
    if (!opt.root) throw new Error("window root cannot be null");
    if (!opt.manager) throw new Error("element manager cannot be null");
    this.#root = opt.root;
    this.#pageSize = opt.pageSize ?? 20;
    this.#manager = opt.manager;
    this.#observedMap = new Map();
    this.#observer = new IntersectionObserver(
      this.#handleIntersectionObserver.bind(this),
      { root: this.#root, rootMargin: "50%" },
    );
  }

  /** @type {IntersectionObserverCallback} */
  #handleIntersectionObserver(observedEntries) {
    observedEntries.forEach((entry) => {
      let element = /** @type {E} */ (entry.target);
      if (entry.isIntersecting) {
        this.#manager.load(element);
        if (this.#observedMap.get("end") === element) {
          this.#load();
        }
      } else {
        this.#manager.unload(element);
      }
    });
  }

  /** @returns {number} */
  #calcPageSize() {
    let rootRect = this.#root.getBoundingClientRect();
    let firstChild = this.#root.children.item(0);
    if (!firstChild) throw new Error("child elements not initialized");
    let childRect = firstChild.getBoundingClientRect();
    let numOfElementsPerRow = Math.floor(rootRect.width / childRect.width);
    return numOfElementsPerRow * 3;
  }

  init() {
    this.#load();
  }

  async #load() {
    if (!this.#pageSize) throw new Error("loading size cannot be null");
    if (this.#total && this.#counter >= this.#total) return;

    let total = await this.#manager.queryElement(
      this.#pageNum,
      this.#pageSize,
      this.#newElement.bind(this),
    );
    if (!this.#total) {
      this.#total = total;
    } else if (this.#total != total) {
      console.warn("data changed in this condition");
      this.#total = total;
    }
    this.#observedMap.set(
      "end",
      /** @type {HTMLElement} */ (this.#root.lastElementChild),
    );
  }

  /** @returns {E} */
  #newElement() {
    if (!this.#pageSize) throw new Error("loading size cannot be null");
    let element = this.#manager.createElement();
    this.#observer.observe(/** @type {HTMLElement} */ (element));
    this.#root.append(/** @type {HTMLElement} */ (element));
    this.#counter += 1;
    if (this.#counter % this.#pageSize == 0) {
      this.#pageNum += 1;
    }
    return element;
  }

  /**
   * @param {number} pageSize
   */
  setPageSize(pageSize) {
    this.#pageSize = pageSize;
  }
}
