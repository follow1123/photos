import CachedPager from "@components/PhotoList/CachedPager";

/**
 * @template E
 * @typedef {Object} Options
 * @property {Element} root 窗口元素
 * @property {import("@components/PhotoList/CachedPager").PageManager<E>} pageMgr
 * @property {() => E} elementProvider
 */

/**
 * 分页窗口
 * @template E
 * @class
 */
export default class PagedWindow {
  /** @type {Element} */
  #root;
  /** @type {IntersectionObserver} */
  #observer;
  /** @type {CachedPager<E>} */
  #pager;
  /** @type {ResizeObserver} */
  #resizeObserver;

  /** @type {Element | null} */
  #pagesStartElement = null;
  /** @type {Element | null} */
  #pagesEndElement = null;

  /**
   * @constructor
   * @param {Options<E>} options
   */
  constructor(options) {
    // 参数校验
    let opt = options || {};
    if (!opt.root) throw new Error("window root cannot be null");
    if (!opt.elementProvider)
      throw new Error("element provider cannot be null");
    if (!opt.pageMgr) throw new Error("page loader cannot be null");

    // 变量初始化
    this.#root = opt.root;

    /** @type {import("@components/PhotoList/FixedSizeInfiniteList").ElementManager<E>} */
    let manager = {
      createElement: opt.elementProvider,
      addElement: (element) =>
        this.#root.appendChild(/** @type {HTMLElement} */ (element)),
      getElement: (index) => {
        let e = this.#root.children[index];
        if (e === undefined)
          throw new Error(`root has not element on index ${index}`);
        return /** @type {E} */ (e);
      },
      moveToTop: (start, end) =>
        this.#root.prepend(
          ...Array.from(this.#root.children).slice(start, end),
        ),
      moveToBottom: (start, end) =>
        this.#root.append(...Array.from(this.#root.children).slice(start, end)),
    };

    this.#pager = new CachedPager({
      elementMgr: manager,
      pageMgr: opt.pageMgr,
    });

    this.#observer = new IntersectionObserver(
      this.#handleIntersectionObserver.bind(this),
      { root: opt.root, threshold: 0.5 },
    );
    this.#resizeObserver = new ResizeObserver(
      this.#handleRootResize.bind(this),
    );
  }

  init() {
    this.#pager.init();

    let pageSize = this.#calcPageSize();
    console.log(pageSize);
    this.#pager.pageSize = pageSize;

    this.#pager.next().then((r) => {
      this.#observe(null, this.#getRootItem(r.end));
    });
    this.#resizeObserver.observe(this.#root);
  }

  deinit() {
    this.#observer.disconnect();
    this.#resizeObserver.disconnect();
  }

  /**
   * @returns {number}
   */
  #calcPageSize() {
    let rootRect = this.#root.getBoundingClientRect();
    let firstChild = this.#root.children.item(0);
    if (!firstChild) throw new Error("child elements not initialized");
    let childRect = firstChild.getBoundingClientRect();
    let numOfElementsPerRow = Math.floor(rootRect.width / childRect.width);
    let numOfElementsPerColumn = Math.ceil(rootRect.height / childRect.height);
    let maxVisibleElements =
      numOfElementsPerRow * numOfElementsPerColumn + numOfElementsPerRow;

    return (
      Math.floor(
        (maxVisibleElements + numOfElementsPerRow) / numOfElementsPerRow,
      ) * numOfElementsPerRow
    );
  }

  /** @type {ResizeObserverCallback} */
  #handleRootResize(entries) {
    entries.forEach(() => {
      this.#pager.resize(this.#calcPageSize());
    });
  }

  /**
   * @type {IntersectionObserverCallback}
   */
  #handleIntersectionObserver(observedEntries) {
    observedEntries.forEach((entry) => {
      if (entry.isIntersecting) {
        let element = entry.target;
        if (element === this.#pagesStartElement) {
          this.#unobserve(true, true);
          if (this.#pager.hasPrevious()) {
            this.#pager.previous().then((r) => {
              this.#observe(
                this.#getRootItem(r.start),
                this.#getRootItem(r.end),
              );
            });
          } else {
            let range = this.#pager.maxPageRange;
            this.#observe(null, this.#getRootItem(range.end));
          }
        } else if (element === this.#pagesEndElement) {
          this.#unobserve(true, true);
          if (this.#pager.hasNext()) {
            this.#pager.next().then((r) => {
              this.#observe(
                this.#getRootItem(r.start),
                this.#getRootItem(r.end),
              );
            });
          } else {
            let range = this.#pager.minPageRange;
            this.#observe(this.#getRootItem(range.start), null);
          }
        }
      }
    });
  }

  /**
   * @param {Element | null} start
   * @param {Element | null} end
   */
  #observe(start, end) {
    if (start) {
      this.#pagesStartElement = start;
      this.#observer.observe(start);
    }
    if (end) {
      this.#pagesEndElement = end;
      this.#observer.observe(end);
    }
  }

  /**
   * @param {boolean} start
   * @param {boolean} end
   */
  #unobserve(start, end) {
    if (start) {
      if (this.#pagesStartElement) {
        this.#observer.unobserve(this.#pagesStartElement);
        this.#pagesStartElement = null;
      }
    }
    if (end) {
      if (this.#pagesEndElement) {
        this.#observer.unobserve(this.#pagesEndElement);
        this.#pagesEndElement = null;
      }
    }
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
}
