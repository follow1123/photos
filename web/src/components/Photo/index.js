// @ts-ignore
import templateText from "@components/Photo/template.html?raw";
// @ts-ignore
import stylesText from "@components/Photo/styles.css?raw";
// @ts-ignore
import imgLoadingUri from "@/assets/image_loading.png";

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
  /** @type {string | null} */
  #photoId = null;
  /** @type {string | null} */
  #originalUri = null;
  /** @type {string | null} */
  #compressedUri = null;
  /** @type {boolean} */
  #loading = false;

  /** @type {HTMLImageElement} */
  #img;

  constructor() {
    super();
    // 创建 Shadow DOM
    const shadow = this.attachShadow({ mode: "open" });
    shadow.append(template.content.cloneNode(true));

    let imgEle = shadow.querySelector("img");
    if (!imgEle) throw new Error("invalid template");

    this.#img = imgEle;
    this.#img.src = imgLoadingUri;
    this.#img.addEventListener("click", this.#dispatchPreviewEvent.bind(this));
  }

  get photoId() {
    if (!this.#photoId) throw new Error("no photo id");
    return this.#photoId;
  }

  set photoId(value) {
    this.#photoId = value;
    this.#originalUri = `http://localhost:8080/photo/${value}/preview/original`;
    this.#compressedUri = `http://localhost:8080/photo/${value}/preview/compressed`;
  }

  unload() {
    this.#img.src = imgLoadingUri;
    this.#loading = false;
  }

  load() {
    if (!this.#compressedUri) throw new Error("compressed uri not exists");
    this.#img.src = this.#compressedUri;
    this.#loading = true;
  }

  /** @type {EventListener} */
  #dispatchPreviewEvent(e) {
    e.stopPropagation();
    if (!this.#loading) return;
    let ce = new CustomEvent("preview", {
      detail: { uri: this.#originalUri },
    });
    this.dispatchEvent(ce);
  }
}
