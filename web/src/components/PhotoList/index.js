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

  /** @type {PagedWindow<Photo>} */
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
    /** @type {import("@components/PhotoList/Pager").PageManager<Photo>} */
    let pageLoader = {
      load: this.#handleLoadPage.bind(this),
      unload: (e) => e.removeAttribute("photo-id"),
      hide: (e) => e.classList.add("hidden"),
      show: (e) => e.classList.remove("hidden"),
    };

    this.#pagedWindow = new PagedWindow(
      /** @type {import("@components/PhotoList/PagedWindow").Options<Photo>} */
      ({
        root: this.#container,
        total: 100,
        pageSize: 20,
        elementProvider: () => {
          let photo = document.createElement("p-photo");
          if (!(photo instanceof Photo)) throw new Error("not photo element");
          return photo;
        },
        pageMgr: pageLoader,
      }),
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
   * @param {() => Photo | null} next
   * @returns {Promise<number>}
   */
  #handleLoadPage(pageNum, pageSize, next) {
    return fetch(
      `http://localhost:8080/photo?pageNum=${pageNum}&pageSize=${pageSize}`,
    )
      .then((resp) => resp.json())
      .then((data) => {
        Array.from(data.list).forEach((item) => {
          let photo = next();
          if (photo === null) return;
          photo.setAttribute("photo-id", item.id);
          photo.addEventListener(
            "preview",
            this.#handlePhotoPreview.bind(this),
          );
        });
        return data.total;
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
