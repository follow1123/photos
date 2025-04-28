// @ts-ignore
import templateText from "@components/Photo/template.html?raw";
// @ts-ignore
import stylesText from "@components/Photo/styles.css?raw";

let tpl = new DOMParser()
  .parseFromString(templateText, "text/html")
  .querySelector("template");
if (!tpl) throw new Error("invalid template");
/** @type {HTMLTemplateElement} */
const template = tpl;

const style = document.createElement("style");
style.textContent = stylesText;
template.content.prepend(style);

export default class Photo extends HTMLElement {
  static observedAttributes = ["photo-id"];

  /** @type {string | null} */
  photoId = null;

  /** @type {HTMLImageElement} */
  img;

  constructor() {
    super();
    // 创建 Shadow DOM
    const shadow = this.attachShadow({ mode: "open" });
    shadow.append(template.content.cloneNode(true));

    let imgEle = shadow.querySelector("img");
    if (!(imgEle instanceof HTMLImageElement)) {
      throw new Error("invalid template");
    }

    this.img = imgEle;
    this.img.addEventListener("load", this.handleImgLoaded.bind(this));
    this.img.addEventListener("click", this.dispatchPreviewEvent.bind(this));
  }

  /**
   * @param {string} name
   * @param {string | null} oldValue
   * @param {string | null} newValue
   */
  attributeChangedCallback(name, oldValue, newValue) {
    if (name === "photo-id") {
      if (newValue === null || newValue === "") {
        this.clearImgSrc();
      } else {
        this.photoId = newValue;
        this.setImgSrc();
      }
    }
  }

  connectedCallback() {
    this.setImgSrc();
  }

  handleImgLoaded() {
    if (this.img.naturalHeight > this.img.naturalWidth) {
      this.img.classList.remove("wrap-heidht");
      this.img.classList.add("wrap-width");
    } else {
      this.img.classList.remove("wrap-width");
      this.img.classList.add("wrap-height");
    }
  }

  setImgSrc() {
    this.img.src = `http://localhost:8080/photo/${this.photoId}/preview/compressed`;
  }

  clearImgSrc() {
    this.img.removeAttribute("src");
    let ce = new CustomEvent("clear");
    this.dispatchEvent(ce);
  }

  /**
   * @param {Event} e
   */
  dispatchPreviewEvent(e) {
    console.log("image clicked");
    e.stopPropagation();
    let ce = new CustomEvent("preview", {
      detail: {
        originalUri: `http://localhost:8080/photo/${this.photoId}/preview/original`,
      },
    });
    this.dispatchEvent(ce);
  }
}
