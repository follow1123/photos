// @ts-ignore
import templateText from "@components/PhotoList/template.html?raw";
// @ts-ignore
import stylesText from "@components/PhotoList/styles.css?raw";

import PagedWindow from "@components/PhotoList/PagedWindow";
import Photo from "@components/Photo";

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

  /** @type {PagedWindow} */
  #pagedWindow;
  constructor() {
    super();
    // 创建 Shadow DOM
    const shadow = this.attachShadow({ mode: "open" });
    shadow.append(template.content.cloneNode(true));

    let originalImgEle = shadow.querySelector("img[original]");
    if (!(originalImgEle instanceof HTMLImageElement))
      throw new Error("invalid template");

    let divEle = shadow.querySelector("#container");
    if (!(divEle instanceof HTMLDivElement))
      throw new Error("invalid template");
    let dialogEle = shadow.querySelector("dialog");
    if (!(dialogEle instanceof HTMLDialogElement))
      throw new Error("invalid template");

    this.#originalImg = originalImgEle;
    this.#dialog = dialogEle;
    this.#container = divEle;
    /** @type {import("@components/PhotoList/CachedPager").PageLoader<HTMLElement>} */
    let pageLoader = {
      load: this.#handleLoadPage.bind(this),
      unload: (next) => {
        let photo = null;
        while ((photo = next()) !== null) {
          if (!(photo instanceof Photo)) return;
          photo.removeAttribute("photo-id");
        }
      },
    };

    this.#pagedWindow = new PagedWindow(
      /** @type {import("@components/PhotoList/fixedSizeViewer").ViewerOptions<>} */
      {
        root: this.#container,
        total: 200,
        pageSize: 20,
        elementProvider: () => document.createElement("p-photo"),
        pageLoader: pageLoader,
      },
    );

    this.#originalImg.addEventListener(
      "load",
      this.#handleOriginalImgLoaded.bind(this),
    );
  }

  connectedCallback() {
    this.#pagedWindow.init();
  }

  /**
   * @param {number} pageNum
   * @param {number} pageSize
   * @param {() => HTMLElement | null} next
   */
  #handleLoadPage(pageNum, pageSize, next) {
    fetch(`http://localhost:8080/photo?pageNum=${pageNum}&pageSize=${pageSize}`)
      .then((resp) => resp.json())
      .then((data) => {
        Array.from(data.list).forEach((item) => {
          let photo = next();
          if (!(photo instanceof Photo)) throw new Error("error photo");
          photo.setAttribute("photo-id", item.id);
          photo.addEventListener(
            "preview",
            this.#handlePhotoPreview.bind(this),
          );
        });
      });
  }
  /**
   * @type EventListener
   */
  #handlePhotoPreview(e) {
    if (!(e instanceof CustomEvent)) throw new Error("not custom event");
    this.#originalImg.src = e.detail.originalUri;
  }

  #handleOriginalImgLoaded() {
    this.#dialog.showModal();
  }
}
