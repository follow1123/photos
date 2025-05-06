import FixedSizeViewer from "@components/PhotoList/fixedSizeViewer";
import Pager from "@components/PhotoList/pager";

/**
 * @typedef {(pageNum: number, pageSize: number, next: () => Element | null) => void} LoadingPage
 */

/**
 * @typedef {Object} PagedWindowOptions
 * @property {Element} root 窗口元素
 * @property {number} pageSize 每页大小
 * @property {number} total 总共要显示的元素
 * @property {() => HTMLElement} elementProvider 创建元素
 * @property {LoadingPage} loadingPageFn 加载页面
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

  /** @type {FixedSizeViewer<HTMLElement>} */
  #viewer;

  /** @type {Pager} */
  #pager;
  /** @type {Map<number, number>} */
  #pageRangeMap;

  /** @type {LoadingPage} */
  #loadingPage;

  /**
   * @constructor
   * @param {PagedWindowOptions} options
   */
  constructor(options) {
    // 参数校验
    let opt = options || {};
    if (!opt.elementProvider) throw new Error("element provider not be null");
    if (!opt.root) throw new Error("window root not be null");
    if (!opt.total) throw new Error("total not be null");
    if (!opt.pageSize) throw new Error("page size not be null");
    if (!opt.loadingPageFn)
      throw new Error("loading page function not be null");

    // 变量初始化
    this.#root = opt.root;
    this.#loadingPage = opt.loadingPageFn;
    this.#pageRangeMap = new Map();
    this.#observedMap = new Map();

    this.#pager = new Pager({ total: opt.total, pageSize: opt.pageSize });

    this.#viewer = new FixedSizeViewer({
      total: opt.total,
      numOfElements: 128,
      createElementFn: opt.elementProvider,
      addElementFn: (element) => this.#root.appendChild(element),
      getElementFn: (index) => {
        let e = this.#root.children.item(index);
        if (!(e instanceof HTMLElement))
          throw new Error(`root has not element on index ${index}`);
        return e;
      },
      moveToTopFn: (start, end) => {
        this.#root.prepend(
          ...Array.from(this.#root.children).slice(start, end),
        );
        this.#organizePageRange(start, end);
      },
      moveToBottomFn: (start, end) => {
        this.#root.append(...Array.from(this.#root.children).slice(start, end));
        this.#organizePageRange(start, end);
      },
    });
    this.#observer = new IntersectionObserver(
      this.#handleIntersectionObserver.bind(this),
      { root: opt.root, threshold: 0.5 },
    );

    //this.#root.addEventListener("wheel", this.#handleRootWheelEvent.bind(this));
  }

  init() {
    //debugger;
    this.#viewer.init();

    this.#observeViewRange(null, this.#viewer.getNext());

    let pageEndIdx = this.#pager.size - 1;
    this.#pageRangeMap.set(1, pageEndIdx);
    this.#observePageRange(null, pageEndIdx);
    this.#handleLoadingPage(1, 0, pageEndIdx);
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
          if (name.endsWith("view")) {
            this.#handleObservedView(name.startsWith("next"));
            console.log("element counter: ", this.#viewer.counter);
          }
          if (name.endsWith("page")) {
            this.#handleObservedPage(name.startsWith("next"));
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
   * @param {boolean} next
   */
  #handleObservedView(next) {
    //debugger;
    console.log("handleObservedView");
    if (next) {
      this.#observeViewRange(
        this.#viewer.getPrevious(),
        this.#viewer.hasNext() ? this.#viewer.next() : null,
      );
    } else {
      this.#observeViewRange(
        this.#viewer.hasPrevious() ? this.#viewer.previous() : null,
        this.#viewer.getNext(),
      );
    }
  }

  /**
   * @param {Element | null} start 起始元素
   * @param {Element | null} end 结束元素
   */
  #observeViewRange(start, end) {
    this.#unobserve("next-view");
    this.#unobserve("prev-view");

    if (start) {
      this.#observe("prev-view", start);
    }
    if (end) {
      this.#observe("next-view", end);
    }
    console.log("view range start: ", start, " end: ", end);
  }

  /**
   * @param {number} start 起始下标
   * @param {number} end 结束下标
   */
  #organizePageRange(start, end) {
    console.log("clean index, start: ", start, " end: ", end);

    if (this.#pageRangeMap.size <= 1)
      throw new Error("range map size cannot be 1");

    let invalidPageNums = new Array();
    this.#pageRangeMap.forEach((pageEnd, pageNum) => {
      let pageStart = pageEnd - (this.#pager.size - 1);
      if (pageStart < end && start < pageEnd) {
        invalidPageNums.push(pageNum);
      } else {
        let newPageEnd =
          start === 0
            ? pageEnd - this.#viewer.numOfMovedElements
            : pageEnd + this.#viewer.numOfMovedElements;
        this.#pageRangeMap.set(pageNum, newPageEnd);
      }
    });

    invalidPageNums.forEach((n) => this.#pageRangeMap.delete(n));

    let minPageNum = Math.min(...this.#pageRangeMap.keys());
    let maxPageNum = Math.max(...this.#pageRangeMap.keys());
    let minPageEndIdx = this.#pageRangeMap.get(minPageNum);
    let maxPageEndIdx = this.#pageRangeMap.get(maxPageNum);

    if (!maxPageEndIdx)
      throw new Error(`unloaded page number: ${maxPageEndIdx}`);
    if (!minPageEndIdx)
      throw new Error(`unloaded page number: ${minPageEndIdx}`);

    this.#observePageRange(minPageEndIdx, maxPageEndIdx);
  }

  /**
   * @param {number | null} minPageEndIdx 最小页结束下标
   * @param {number | null} maxPageEndIdx 最大页结束下标
   */
  #observePageRange(minPageEndIdx, maxPageEndIdx) {
    console.log("page range start: ", minPageEndIdx, " end: ", maxPageEndIdx);
    console.log(
      "page range map: ",
      JSON.stringify(Array.from(this.#pageRangeMap.entries())),
    );
    this.#unobserve("next-page");
    this.#unobserve("prev-page");
    if (minPageEndIdx) {
      let e = this.#getRootItem(minPageEndIdx - (this.#pager.size - 1));
      this.#observe("prev-page", e);
    }

    if (maxPageEndIdx) {
      let e = this.#getRootItem(maxPageEndIdx);
      this.#observe("next-page", e);
    }
  }
  /**
   * @param {number} pageNum
   * @param {number} start
   * @param {number} end
   */
  #handleLoadingPage(pageNum, start, end) {
    let loadingIdx = start;
    let nextElement = () =>
      loadingIdx <= end ? this.#getRootItem(loadingIdx++) : null;
    this.#loadingPage(pageNum, this.#pager.size, nextElement);
  }

  /**
   * @param {boolean} next
   */
  #handleObservedPage(next) {
    let minPageNum = Math.min(...this.#pageRangeMap.keys());
    let maxPageNum = Math.max(...this.#pageRangeMap.keys());
    let minPageEndIdx = this.#pageRangeMap.get(minPageNum);
    let maxPageEndIdx = this.#pageRangeMap.get(maxPageNum);

    if (!maxPageEndIdx)
      throw new Error(`unloaded page number: ${maxPageEndIdx}`);
    if (!minPageEndIdx)
      throw new Error(`unloaded page number: ${minPageEndIdx}`);

    if (next) {
      let nextPageNum = maxPageNum + 1;
      let endPageEndIdx = null;

      if (this.#pager.hasNext(nextPageNum)) {
        let nextPageEndIdx = maxPageEndIdx + this.#pager.size;

        this.#handleLoadingPage(
          nextPageNum,
          nextPageEndIdx - (this.#pager.size - 1),
          nextPageEndIdx,
        );

        this.#pageRangeMap.set(nextPageNum, nextPageEndIdx);
        endPageEndIdx = nextPageEndIdx;
      }

      this.#observePageRange(minPageEndIdx, endPageEndIdx);
    } else {
      let prevPageNum = minPageNum - 1;
      let startPageEndIdx = null;
      if (this.#pager.hasPrevious(prevPageNum)) {
        let prevPageEndIdx = minPageEndIdx - this.#pager.size;

        this.#handleLoadingPage(
          prevPageNum,
          prevPageEndIdx - (this.#pager.size - 1),
          prevPageEndIdx,
        );

        this.#pageRangeMap.set(prevPageNum, prevPageEndIdx);
        startPageEndIdx = prevPageEndIdx;
      }
      this.#observePageRange(startPageEndIdx, maxPageEndIdx);
    }
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
