// Tests kit.js in a small stand-in DOM: elements, attribute and class selectors, events and the address.
import { test } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import vm from "node:vm";

const source = readFileSync(new URL("../js/kit.js", import.meta.url), "utf8");

// One compound selector: a tag, classes and attributes, like a.ui-tab[data-ui-carry] or [x='y'].
function compound(text) {
  const tag = (text.match(/^[a-z]+/) || [""])[0];
  const classes = [...text.matchAll(/\.([\w-]+)/g)].map((m) => m[1]);
  const attrs = [...text.matchAll(/\[([\w-]+)(?:=["']?([^"'\]]*)["']?)?\]/g)].map((m) => [m[1], m[2]]);
  return (el) => (!tag || el.tagName === tag) && classes.every((c) => (el.attrs.class || "").split(" ").includes(c)) &&
    attrs.every(([k, v]) => k in el.attrs && (v === undefined || el.attrs[k] === v));
}

function matcher(selector) {
  const alternatives = selector.split(",").map((s) => s.trim().split(/\s+/).map(compound));
  return (el) => alternatives.some((parts) => {
    if (!parts[parts.length - 1](el)) return false;
    let i = parts.length - 2;
    for (let p = el.parent; p && i >= 0; p = p.parent) if (parts[i](p)) i--;
    return i < 0;
  });
}

class El {
  constructor(tagName, attrs = {}, children = []) {
    this.tagName = tagName;
    this.attrs = { ...attrs };
    this.children = [];
    this.parent = null;
    this.value = attrs.value;
    children.forEach((c) => this.appendChild(c));
  }
  getAttribute(k) { return k in this.attrs ? this.attrs[k] : null; }
  setAttribute(k, v) { this.attrs[k] = String(v); }
  hasAttribute(k) { return k in this.attrs; }
  removeAttribute(k) { delete this.attrs[k]; }
  addEventListener() {}
  appendChild(c) {
    if (c.parent) c.parent.children.splice(c.parent.children.indexOf(c), 1);
    c.parent = this;
    this.children.push(c);
    return c;
  }
  set innerHTML(html) {
    this.children.forEach((c) => { c.parent = null; });
    this.children = [];
    (El.parse[html] || (() => []))().forEach((c) => this.appendChild(c));
  }
  all() { return this.children.flatMap((c) => (c.tagName === "template" ? [c] : [c, ...c.all()])); }
  querySelectorAll(sel) { return this.all().filter(matcher(sel)); }
  querySelector(sel) { return this.querySelectorAll(sel)[0] || null; }
  closest(sel) {
    const m = matcher(sel);
    for (let e = this; e; e = e.parent) if (m(e)) return e;
    return null;
  }
}
El.parse = {};

// page loads kit.js on a document whose body holds children, at the address href.
function page(href, children, fetched = {}) {
  const body = new El("body", {}, children);
  const handlers = {};
  const location = {};
  const go = (to) => {
    const u = new URL(to, location.href || "http://localhost/");
    Object.assign(location, { href: u.href, origin: u.origin, pathname: u.pathname, search: u.search });
  };
  go(href);
  const fetches = [];
  const document = {
    documentElement: new El("html"), readyState: "complete",
    querySelectorAll: (s) => body.querySelectorAll(s), querySelector: (s) => body.querySelector(s),
    getElementById: (id) => body.all().find((e) => e.attrs.id === id) || null,
    addEventListener: (type, fn) => { handlers[type] = fn; },
  };
  const window = {
    location, document, fetch: (url) => {
      fetches.push(url);
      const text = fetched[url];
      return text === undefined ? new Promise(() => {}) : Promise.resolve({ ok: true, status: 200, text: () => Promise.resolve(text) });
    },
    history: { replaceState: (_s, _t, url) => go(String(url)) },
    URLSearchParams, URL,
  };
  window.window = window;
  vm.runInNewContext(source, { ...window, localStorage: { getItem: () => null, setItem() {} },
    FormData: function (form) { return form.fields; }, matchMedia: () => ({ matches: false }), setTimeout });
  const fire = (type, target) => { const e = { target, preventDefault() { e.prevented = true; } }; handlers[type](e); return e; };
  return { location, fetches, fire };
}

// The control room's bar, product tabs and lazy GitHub light, drawn for /product/kanban?tab=prs.
function product() {
  const refresh = new El("a", { class: "ui-iconbtn", href: "/product/kanban?tab=prs&fresh=1", "data-ui-carry": "fresh" });
  const board = new El("a", { class: "ui-tab", href: "/product/kanban?tab=board", "data-ui-carry": "tab" });
  const light = new El("div", { class: "ui-lazy", "data-ui-src": "/part/light/github?back=%2Fproduct%2Fkanban%3Ftab%3Dprs", "data-ui-back": "back" });
  const back = new El("input", { type: "hidden", name: "back", value: "/product/kanban?tab=prs", "data-ui-back": "" });
  const outside = new El("a", { href: "https://github.com/x?y=1", "data-ui-carry": "" });
  const form = new El("form", { method: "get" });
  form.method = "get";
  form.fields = [["pr_sort", "newest"], ["tab", "prs"]];
  const section = new El("div", { class: "ui-lazy", "data-ui-src": "/part/product/kanban/github?tab=prs" }, [new El("section", { class: "ui-section" }, [form])]);
  return { els: { refresh, board, light, back, outside, form, section }, children: [refresh, board, light, back, outside, section] };
}

test("a sort that reloads only its section carries the new address to Read again, the tabs and the light's back", () => {
  const { els, children } = product();
  const p = page("http://localhost/product/kanban?tab=prs", children);
  const e = p.fire("submit", els.form);
  assert.ok(e.prevented);
  assert.equal(p.location.search, "?tab=prs&pr_sort=newest");
  assert.equal(els.refresh.getAttribute("href"), "/product/kanban?tab=prs&pr_sort=newest&fresh=1");
  assert.equal(els.board.getAttribute("href"), "/product/kanban?pr_sort=newest&tab=board");
  assert.equal(new URL(els.light.getAttribute("data-ui-src"), "http://x").searchParams.get("back"), "/product/kanban?tab=prs&pr_sort=newest");
  assert.equal(els.back.value, "/product/kanban?tab=prs&pr_sort=newest");
  assert.equal(els.outside.getAttribute("href"), "https://github.com/x?y=1", "a link to another site keeps its address");
  assert.deepEqual(p.fetches.slice(-1), ["/part/product/kanban/github?tab=prs&pr_sort=newest"]);
});

test("the cards or list toggle carries the view the same way", () => {
  const refresh = new El("a", { href: "/products?fresh=1", "data-ui-carry": "fresh" });
  const list = new El("button", { "data-view": "list" });
  const seg = new El("div", { "data-ui-views": "products" }, [new El("button", { "data-view": "cards" }), list]);
  const p = page("http://localhost/products", [refresh, seg]);
  p.fire("click", list);
  assert.equal(p.location.search, "?view=list");
  assert.equal(refresh.getAttribute("href"), "/products?view=list&fresh=1");
});

test("Read again on a lazy section reloads only it with fresh=1, keeps its address, and stays in the loaded head", async () => {
  const again = new El("a", { class: "ui-iconbtn", href: "/?pr_sort=newest&fresh=1", "data-ui-again": "", "data-ui-carry": "fresh" });
  const icon = new El("svg");
  again.appendChild(icon);
  const head = new El("div", { class: "ui-section-head" }, [again]);
  const box = new El("div", { class: "ui-lazy", "data-ui-src": "/part/focus/decisions?pr_sort=newest" }, [new El("section", { class: "ui-section" }, [head])]);
  El.parse["<loaded>"] = () => [new El("section", { class: "ui-section" }, [new El("div", { class: "ui-section-head" }), new El("p", { id: "new" })])];
  const p = page("http://localhost/?pr_sort=newest", [box], { "/part/focus/decisions?pr_sort=newest": "<loaded>", "/part/focus/decisions?pr_sort=newest&fresh=1": "<loaded>" });
  await new Promise((r) => setTimeout(r, 0));
  const e = p.fire("click", icon);
  assert.ok(e.prevented);
  assert.equal(p.fetches.at(-1), "/part/focus/decisions?pr_sort=newest&fresh=1");
  assert.equal(box.getAttribute("data-ui-src"), "/part/focus/decisions?pr_sort=newest", "the next plain load reads from the cache again");
  assert.equal(box.querySelector("section").getAttribute("aria-busy"), "true");
  await new Promise((r) => setTimeout(r, 0));
  assert.ok(box.querySelector("[id=new]"), "the section shows what the server read");
  assert.equal(box.querySelector(".ui-section-head [data-ui-again]"), again);
});
