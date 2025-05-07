import CachedPager from "@components/PhotoList/CachedPager";

/**
 * @typedef {Object} Options
 * @property {Element} root 窗口元素
 * @property {number} pageSize 每页大小
 * @property {number} total 总共要显示的元素
 * @property {import("@components/PhotoList/CachedPager").PageLoader<HTMLElement>} pageLoader
 * @property {() => HTMLElement} elementProvider
 */

/**
 * 分页窗口
 * @class
 */
export default class PagedWindow {
  /** @type {Element} */
  #root;
  /** @type {IntersectionObserver} */
  #observer;
  /** @type {Map<string, Element | null>} */
  #observedMap;

  /** @type {CachedPager<HTMLElement>} */
  #pager;

  /**
   * @constructor
   * @param {Options} options
   */
  constructor(options) {
    // 参数校验
    let opt = options || {};
    if (!opt.root) throw new Error("window root cannot be null");
    if (!opt.total) throw new Error("total cannot be null");
    if (!opt.pageSize) throw new Error("page size cannot be null");
    if (!opt.elementProvider)
      throw new Error("element provider cannot be null");
    if (!opt.pageLoader) throw new Error("page loader cannot be null");

    // 变量初始化
    this.#root = opt.root;
    this.#observedMap = new Map();

    /** @type {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<HTMLElement>} */
    let manager = {
      createElement: opt.elementProvider,
      addElement: (element) => this.#root.appendChild(element),
      getElement: (index) => {
        let e = this.#root.children.item(index);
        if (!(e instanceof HTMLElement))
          throw new Error(`root has not element on index ${index}`);
        return e;
      },
      moveToTop: (start, end) =>
        this.#root.prepend(
          ...Array.from(this.#root.children).slice(start, end),
        ),
      moveToBottom: (start, end) =>
        this.#root.append(...Array.from(this.#root.children).slice(start, end)),
    };

    this.#pager = new CachedPager({
      pageSize: opt.pageSize,
      total: opt.total,
      manager: manager,
      loader: opt.pageLoader,
    });

    this.#observer = new IntersectionObserver(
      this.#handleIntersectionObserver.bind(this),
      { root: opt.root, threshold: 0.5 },
    );

    //this.#root.addEventListener("wheel", this.#handleRootWheelEvent.bind(this));
  }

  init() {
    this.#pager.init();

    let range = this.#pager.next();
    this.#observe("pages-end", this.#getRootItem(range.end));
  }

  /**
   * @returns {number}
   */
  getNumOfElementsPerRow() {
    let rootWidth = this.#root.getBoundingClientRect().width;
    let childWidth = this.#root.children[0].getBoundingClientRect().width;
    return Math.floor(rootWidth / childWidth);
  }

  /**
   * @type {IntersectionObserverCallback}
   */
  #handleIntersectionObserver(observedEntries) {
    observedEntries.forEach((entry) => {
      if (entry.isIntersecting) {
        let element = entry.target;
        this.#observedMap.forEach((e, name) => {
          if (e !== element) return;
          if (name === "pages-start") {
            this.#unobserve("pages-start");
            this.#unobserve("pages-end");

            let minPageNum = Math.min(...this.#pager.cache.keys());
            if (minPageNum > 1) {
              let range = this.#pager.previous();
              this.#observe("pages-start", this.#getRootItem(range.start));
              this.#observe("pages-end", this.#getRootItem(range.end));
            } else {
              let maxPageNum = Math.max(...this.#pager.cache.keys());
              let range = this.#pager.getExistsCacheRange(maxPageNum);
              this.#observe("pages-end", this.#getRootItem(range.end));
            }
          } else if (name === "pages-end") {
            this.#unobserve("pages-start");
            this.#unobserve("pages-end");
            let maxPageNum = Math.max(...this.#pager.cache.keys());
            if (maxPageNum < this.#pager.maxNum) {
              let range = this.#pager.next();
              this.#observe("pages-start", this.#getRootItem(range.start));
              this.#observe("pages-end", this.#getRootItem(range.end));
            } else {
              let minPageNum = Math.min(...this.#pager.cache.keys());
              let range = this.#pager.getExistsCacheRange(minPageNum);
              this.#observe("pages-start", this.#getRootItem(range.start));
            }
          }
        });
      }
    });
  }

  /**
   * @param {string} name
   * @param {Element} element
   */
  #observe(name, element) {
    if (this.#observedMap.get(name) === element) return;

    this.#observedMap.set(name, element);
    let hasOther = false;
    this.#observedMap.forEach((e, n) => {
      if (name !== n && e === element) {
        hasOther = true;
      }
    });
    if (!hasOther) this.#observer.observe(element);
  }

  /**
   * @param {string} name
   */
  #unobserve(name) {
    let ele = this.#observedMap.get(name);
    if (!ele) return;

    this.#observedMap.set(name, null);
    let hasOther = false;
    this.#observedMap.forEach((e, n) => {
      if (name !== n && e === ele) {
        hasOther = true;
      }
    });

    if (!hasOther) this.#observer.unobserve(ele);
  }

  /**
   * @param {number} index
   * @returns {Element}
   */
  #getRootItem(index) {
    let e = this.#root.children.item(index);
    if (!e) throw new Error(`root has not element on index ${index}`);
    return e;
  }

  /**
   * @type EventListener
   */
  #handleRootWheelEvent(e) {
    if (!(e instanceof WheelEvent)) return;
    if (e.ctrlKey) return;
    // 禁用默认的滚动行为
    e.preventDefault();

    let offset = 80;

    // 根据滚轮方向更新 scrollTop
    if (e.deltaY > 0) {
      //box.scrollTo({ top: box.scrollTop + offset, behavior: "smooth" });
      this.#root.scrollTop += offset;
    } else {
      //box.scrollTo({ top: box.scrollTop - offset, behavior: "smooth" });
      this.#root.scrollTop -= offset;
    }
  }
}
