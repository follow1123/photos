/**
 * @typedef {Object} Options
 * @property {string} desc
 */

export default class Condition {
  /** @type {Options} */
  #opts;

  /**
   * @constructor
   * @param {Options} options
   */
  constructor(options) {
    if (!options) throw new Error("condition cannot be null");
    this.#opts = options;
  }

  /**
   * TODO: 后面换成复杂条件
   * @returns {string}
   */
  build() {
    return `&desc=${this.#opts.desc}`;
  }
}
