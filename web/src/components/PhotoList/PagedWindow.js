import Pager from "@components/PhotoList/Pager";

/**
 * @template E
 * @typedef {Object} Options
 * @property {Element} root 窗口元素
 * @property {import("@components/PhotoList/Pager").PageManager<E>} pageMgr
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
  /** @type {Pager<E>} */
  #pager;
  /** @type {ResizeObserver} */
  #resizeObserver;

  /** @type {Map<string, Element>} */
  #observedMap;

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
    this.#observedMap = new Map();

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

    this.#pager = new Pager({
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
      this.#observedMap.set("pages-end", this.#getRootItem(r.end));
      this.#observedMap.forEach((e) => this.#observer.observe(e));
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
      let pageSize = this.#calcPageSize();
      console.log(
        "current page size: ",
        this.#pager.pageSize,
        ", new page size: ",
        pageSize,
      );
      if (this.#pager.pageSize === pageSize) return;

      //let startIdx = this.#observedMap.has("pages-start")
      //  ? this.#pager.minRange.start
      //  : 0;
      //let endIdx = this.#observedMap.has("pages-end")
      //  ? this.#pager.maxRange.end
      //  : this.#root.children.length - 1;
      //console.log(startIdx, endIdx);
      //
      //let rootRect = this.#root.getBoundingClientRect();
      //console.log(rootRect);
      //let topLeftIdx = null;
      //let bottomRightIdx = null;
      //for (let i = startIdx; i <= endIdx; i++) {
      //  let child = this.#root.children[i];
      //  let childRect = child.getBoundingClientRect();
      //  if (childRect.top > rootRect.bottom) continue;
      //  //console.log(childRect);
      //  if (topLeftIdx === null) {
      //    topLeftIdx = i;
      //  } else {
      //    let curTopDistance = childRect.bottom - rootRect.top;
      //    let curLeftDistance = childRect.right - rootRect.left;
      //    let topLeft = this.#root.children[topLeftIdx].getBoundingClientRect();
      //    let topDistance = childRect.bottom - topLeft.top;
      //    let leftDistance = childRect.right - topLeft.left;
      //    if (
      //      curTopDistance <= topDistance &&
      //      curLeftDistance <= leftDistance
      //    ) {
      //      //console.log("new top left: ", child, childRect);
      //      topLeftIdx = i;
      //    }
      //  }
      //  if (bottomRightIdx === null) {
      //    //debugger;
      //    bottomRightIdx = i;
      //  } else {
      //    let curBottomDistance = rootRect.bottom - childRect.top;
      //    let curRightDistance = rootRect.right - childRect.left;
      //    let bottomRight =
      //      this.#root.children[bottomRightIdx].getBoundingClientRect();
      //    let bottomDistance = rootRect.bottom - bottomRight.top;
      //    let rightDistance = rootRect.right - bottomRight.left;
      //
      //    if (
      //      curRightDistance <= rightDistance &&
      //      curBottomDistance <= bottomDistance
      //    ) {
      //      //console.log("new bottom right: ", child, childRect);
      //      bottomRightIdx = i;
      //    }
      //  }
      //}
      //console.log("左上角元素:", topLeftIdx, this.#root.children[topLeftIdx]);
      //console.log(
      //  "右下角元素:",
      //  bottomRightIdx,
      //  this.#root.children[bottomRightIdx],
      //);

      console.log(
        "before resize: ",
        this.#pager.minRange.start,
        this.#pager.maxRange.end,
      );
      this.#pager.resize(pageSize).then((r) => {
        this.#observedMap.forEach((e) => this.#observer.unobserve(e));
        if (this.#pager.hasPrevious()) {
          this.#observedMap.set("pages-start", this.#getRootItem(r.start));
        }
        if (this.#pager.hasNext()) {
          this.#observedMap.set("pages-end", this.#getRootItem(r.end));
        }
        this.#observedMap.forEach((e) => this.#observer.observe(e));
        console.log(
          "after resize: ",
          this.#pager.minRange.start,
          this.#pager.maxRange.end,
        );
      });
    });
  }

  /**
   * @type {IntersectionObserverCallback}
   */
  #handleIntersectionObserver(observedEntries) {
    observedEntries.forEach((entry) => {
      if (entry.isIntersecting) {
        let element = entry.target;
        if (this.#observedMap.get("pages-start") === element) {
          this.#observedMap.forEach((e) => this.#observer.unobserve(e));
          if (this.#pager.hasPrevious()) {
            this.#pager.previous().then((r) => {
              this.#observedMap.set("pages-start", this.#getRootItem(r.start));
              this.#observedMap.set("pages-end", this.#getRootItem(r.end));
              this.#observedMap.forEach((e) => this.#observer.observe(e));
            });
          } else {
            let range = this.#pager.maxRange;
            this.#observedMap.set("pages-end", this.#getRootItem(range.end));
            this.#observedMap.forEach((e) => this.#observer.observe(e));
          }
        } else if (this.#observedMap.get("pages-end") === element) {
          this.#observedMap.forEach((e) => this.#observer.unobserve(e));
          if (this.#pager.hasNext()) {
            this.#pager.next().then((r) => {
              this.#observedMap.set("pages-start", this.#getRootItem(r.start));
              this.#observedMap.set("pages-end", this.#getRootItem(r.end));
              this.#observedMap.forEach((e) => this.#observer.observe(e));
            });
          } else {
            let range = this.#pager.minRange;
            this.#observedMap.set(
              "pages-start",
              this.#getRootItem(range.start),
            );
            this.#observedMap.forEach((e) => this.#observer.observe(e));
          }
        }
      }
    });
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
