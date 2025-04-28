// @ts-ignore
import templateText from "@components/PhotoList/template.html?raw";
// @ts-ignore
import stylesText from "@components/PhotoList/styles.css?raw";

import Photo from "@components/Photo";

let tpl = new DOMParser()
  .parseFromString(templateText, "text/html")
  .querySelector("template");
if (!tpl) throw new Error("invalid template");
/** @type {HTMLTemplateElement} */
const template = tpl;

const style = document.createElement("style");
style.textContent = stylesText;
template.content.prepend(style);

class ImagePool {
  /** @type {number} */
  #maxSize;

  /** @type {number} */
  #idleSize;

  /** @type {Array<Photo>} */
  #photos;

  /**
   * @param {number} maxSize
   * @param {number} idleSize
   * @param {EventListenerOrEventListenerObject} handlePreview
   * @param {EventListenerOrEventListenerObject} handleClear
   */
  constructor(maxSize, idleSize, handlePreview, handleClear) {
    this.#maxSize = maxSize;
    this.#idleSize = idleSize;

    /** @type {Array<Photo>} */
    let imgs = new Array(this.#maxSize);
    this.#photos = imgs;

    for (let i = 0; i < maxSize; i++) {
      let photo = document.createElement("p-photo");
      if (!(photo instanceof Photo)) {
        throw new Error("create photo error");
      }
      photo.addEventListener("preview", handlePreview);
      photo.addEventListener("clear", handleClear);
      this.#photos.push(photo);
    }
  }
}

/**
  *
我创建了一个 web component 这个组件的长和宽是固定的

其中里面有固定28个 img 标签用于存放图片

这个组件的可视窗口可以显示 16 个 img 元素，每列 4个 每行 4个

28 个 img 标签中有 24 个直接加载图片，4个空闲

我会创建一个div标签模拟滚动条

当模拟滚动时空闲的4个 img 开始指定 src 属性加载图片，而最上面的4个移除 src属性释放图片，保持28个加载4个空闲的状态
  */
export default class PhotoList extends HTMLElement {
  /** @type {HTMLDivElement} */
  container;

  /** @type {HTMLImageElement} */
  originalImg;

  /** @type {HTMLDialogElement} */
  dialog;

  constructor() {
    super();
    // 创建 Shadow DOM
    const shadow = this.attachShadow({ mode: "open" });
    shadow.append(template.content.cloneNode(true));

    let originalImgEle = shadow.querySelector("img[original]");
    if (!(originalImgEle instanceof HTMLImageElement)) {
      throw new Error("invalid template");
    }

    let divEle = shadow.querySelector("#container");
    if (!(divEle instanceof HTMLDivElement)) {
      throw new Error("invalid template");
    }
    let dialogEle = shadow.querySelector("dialog");
    if (!(dialogEle instanceof HTMLDialogElement)) {
      throw new Error("invalid template");
    }

    this.originalImg = originalImgEle;
    this.dialog = dialogEle;
    this.container = divEle;

    this.originalImg.addEventListener(
      "load",
      this.handleOriginalImgLoaded.bind(this),
    );
  }

  connectedCallback() {
    fetch("http://localhost:8080/photo?pageNum=1&pageSize=30")
      .then((resp) => resp.json())
      .then((data) => {
        let list = Array.from(data.list);
        list.forEach((item) => {
          let photo = document.createElement("p-photo");
          photo.setAttribute("photo-id", item.id);
          photo.addEventListener("preview", this.handlePhotoPreview.bind(this));
          this.container.append(photo);
        });
      })
      .catch((error) => {
        console.error("Error:", error); // 错误处理
      });
  }

  /**
   * @param {Event} e
   */
  handlePhotoPreview(e) {
    console.log("handle preview in photo list");
    if (!(e instanceof CustomEvent)) {
      throw new Error("not custom event");
    }
    this.originalImg.src = e.detail.originalUri;
  }

  handleClearPhoto() {}

  handleOriginalImgLoaded() {
    this.dialog.showModal();
  }
}
