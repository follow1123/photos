// @ts-ignore
import templateText from "@components/SearchBar/template.html?raw";
// @ts-ignore
import stylesText from "@components/SearchBar/styles.css?raw";
import { eventBus } from "@/eventbus";
import Condition from "./Condition";

let tpl = new DOMParser()
  .parseFromString(templateText, "text/html")
  .querySelector("template");
if (!tpl) throw new Error("invalid template");
/** @type {HTMLTemplateElement} */
const template = tpl;

const style = document.createElement("style");
style.textContent = stylesText;
template.content.prepend(style);

export default class SearchBar extends HTMLElement {
  constructor() {
    super();
    // 创建 Shadow DOM
    const shadow = this.attachShadow({ mode: "open" });
    shadow.append(template.content.cloneNode(true));

    let inputEle = shadow.querySelector("input");
    if (!inputEle) throw new Error("invalid template");

    inputEle.addEventListener("keydown", (e) => {
      if (e.key === "Enter") {
        let input = /** @type {HTMLInputElement} */ (e.target);
        if (!input.value || input.value === "") return;
        console.log("Enter pressed!, value: ", input.value);
        eventBus.emit("query", new Condition({ desc: input.value }));
      }
    });
  }
}
