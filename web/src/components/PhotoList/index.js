// @ts-ignore
import templateText from "@components/PhotoList/template.html?raw";
// @ts-ignore
import stylesText from "@components/PhotoList/styles.css?raw";

import Photo from "@components/Photo";
import InfiniteWindowList from "@components/PhotoList/InfiniteWindowList";

/** @typedef {import('@/eventbus').EventHandler} EventHandler */

let tpl = new DOMParser()
  .parseFromString(templateText, "text/html")
  .querySelector("template");
if (!tpl) throw new Error("invalid template");
/** @type {HTMLTemplateElement} */
const template = tpl;

const style = document.createElement("style");
style.textContent = stylesText;
template.content.prepend(style);

/**
 * @class
 */
export default class PhotoList extends HTMLElement {
  /** @type {HTMLDivElement} */
  #container;
  /** @type {HTMLImageElement} */
  #originalImg;
  /** @type {HTMLDialogElement} */
  #dialog;

  /** @type {InfiniteWindowList<Photo>} */
  #windowList;

  constructor() {
    super();
    // 创建 Shadow DOM
    const shadow = this.attachShadow({ mode: "open" });
    shadow.append(template.content.cloneNode(true));

    let originalImgEle = shadow.querySelector("img[original]");
    if (!originalImgEle) throw new Error("invalid template");
    let divEle = shadow.querySelector("#container");
    if (!divEle) throw new Error("invalid template");
    let dialogEle = shadow.querySelector("dialog");
    if (!dialogEle) throw new Error("invalid template");
    this.#originalImg = /** @type {HTMLImageElement} */ (originalImgEle);
    this.#dialog = dialogEle;
    this.#container = /** @type {HTMLDivElement} */ (divEle);

    /** @type {import("@components/PhotoList/InfiniteWindowList").ElementManager<Photo>} */
    let elementManager = {
      load: (e) => e.load(),
      unload: (e) => e.unload(),
      createElement: () =>
        /** @type {Photo} */ (document.createElement("p-photo")),
      queryElement: this.#handleLoadPage.bind(this),
    };
    this.#windowList = new InfiniteWindowList({
      root: this.#container,
      manager: elementManager,
    });

    this.#originalImg.addEventListener(
      "load",
      this.#handleOriginalImgLoaded.bind(this),
    );
  }

  connectedCallback() {
    this.#windowList.init();
  }

  /** @type {import("@components/PhotoList/InfiniteWindowList").QueryElementsFn<Photo>} */
  async #handleLoadPage(pageNum, pageSize, next) {
    return await fetch(
      `http://localhost:8080/photo?pageNum=${pageNum}&pageSize=${pageSize}`,
    )
      .then((resp) => resp.json())
      .then((data) => {
        Array.from(data.list).forEach((item) => {
          let photo = next();
          photo.photoId = item.id;
          photo.addEventListener(
            "preview",
            this.#handlePhotoPreview.bind(this),
          );
        });
        return data.total;
      });
  }

  /** @type {EventListener} */
  #handlePhotoPreview(e) {
    const ce = /** @type {CustomEvent} */ (e);
    this.#originalImg.src = ce.detail.uri;
  }

  #handleOriginalImgLoaded() {
    this.#dialog.showModal();
  }
}
