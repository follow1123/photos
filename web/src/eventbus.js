/**
 * @typedef {function(*): void} EventHandler
 */

class EventBus {
  static id = 1;
  name = `event-bus-${EventBus.id}`;
  /** @type {Map<string, Map<Object, EventHandler>>} */
  #eventsMap = new Map();

  constructor() {}

  /**
   * @param {string} name
   * @param {EventHandler} handler
   * @param {Object} caller
   */
  on(name, handler, caller) {
    let events = this.#eventsMap.get(name);
    if (events) {
      if (events.has(caller)) {
        console.warn(
          `event: '${name}' in ${caller} already registered, override this function`,
        );
      }
      events.set(caller, handler);
    } else {
      events = new Map();
      events.set(caller, handler);
      this.#eventsMap.set(name, events);
    }
  }

  /**
   * @param {string} name
   * @param {any} [data]
   */
  emit(name, data) {
    let events = this.#eventsMap.get(name);
    if (events) {
      events.forEach((v, k) => v.call(k, data));
    } else {
      console.warn(`no one handle this event: '${name}'`);
    }
  }

  /**
   * @param {string} name
   */
  off(name) {
    this.#eventsMap.delete(name);
  }

  /**
   * @param {Object} caller
   */
  offCaller(caller) {
    this.#eventsMap.forEach((events, name) => {
      if (events.has(caller)) {
        events.delete(caller); // 删除 caller 对应的事件监听器
        if (events.size === 0) {
          this.#eventsMap.delete(name); // 如果没有监听器了，删除事件名条目
        }
      }
    });
    console.log(this.#eventsMap);
  }
}

export const eventBus = new EventBus();
