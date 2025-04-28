import { resolve } from "path";

/** @type {import('vite').UserConfig} */
export default {
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"), // 使用 @ 代表 src 目录
      "@components": resolve(__dirname, "./src/components"),
    },
  },
};
