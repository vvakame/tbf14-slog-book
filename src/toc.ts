import fs from "fs";

const config = require("../vivliostyle.config");

const entry = config.entry;
if (!Array.isArray(entry)) {
  throw new Error(`entry is not array`);
}

const result = entry
  .map((entry: any): string | null => {
    if (typeof entry === "string") {
      return entry;
    } else if (entry.rel === "contents") {
      return null;
    } else if (typeof entry.path === "string") {
      return entry.path as string;
    }
    return null;
  })
  .filter((s): s is string => !!s)
  .filter((s) => s.endsWith(".md"))
  .filter((s) => !s.includes("colophon"))
  .map((filePath) => {
    const content = fs.readFileSync(filePath, { encoding: "utf8" });
    return {
      filePath,
      content,
    };
  })
  .map((file) => {
    let filePath = file.filePath;
    filePath = filePath.split("/")[1] || "";
    filePath = filePath.replace(".md", ".html");
    return {
      ...file,
      filePath,
    };
  })
  .map((file) => {
    return file.content
      .split("\n")
      .filter((s) => s.trim())
      .filter((s) => s.startsWith("#"))
      .map((s) => {
        const result = /^(#+)\s+(.*)$/.exec(s);
        if (!result) {
          return "";
        }
        const level = result[1].length;
        const title = result[2];
        let className = "";
        switch (level) {
          case 1:
            className = "toc-chapter";
            break;
          case 2:
            className = "toc-section";
            break;
          case 3:
            className = "toc-subsection";
            break;
        }
        let id = title.toLowerCase();
        return `${"  ".repeat(level - 1)}- <a class="${className}" href="${
          file.filePath
        }#${id}">${title}</a>`;
      })
      .join("\n");
  })
  .join("\n");

console.log(result);
