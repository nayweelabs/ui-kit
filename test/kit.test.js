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
  focus() { El.focused = this; }
  showModal() { this.open = true; }
  remove() { if (this.parent) { this.parent.children.splice(this.parent.children.indexOf(this), 1); this.parent = null; } }
  appendChild(c) {
    if (c.parent) c.parent.children.splice(c.parent.children.indexOf(c), 1);
    c.parent = this;
    this.children.push(c);
    return c;
  }
  set textContent(t) { this.text = t; }
  get textContent() { return this.text || ""; }
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
    body, createElement: (tag) => new El(tag),
    addEventListener: (type, fn) => { handlers[type] = fn; },
  };
  const window = {
    location, document, fetch: (url) => {
      fetches.push(url);
      const text = fetched[url];
      const headers = { get: (k) => (k === "UI-Read-At" ? fetched.readAt || null : null) };
      return text === undefined ? new Promise(() => {}) : Promise.resolve({ ok: true, status: 200, headers, text: () => Promise.resolve(text) });
    },
    history: { replaceState: (_s, _t, url) => go(String(url)) },
    URLSearchParams, URL,
  };
  window.window = window;
  vm.runInNewContext(source, { ...window, localStorage: { getItem: () => null, setItem() {} },
    FormData: function (form) { return form.fields; }, matchMedia: () => ({ matches: false }), setTimeout, clearTimeout });
  const fire = (type, target, more) => { const e = { target, ...more, preventDefault() { e.prevented = true; } }; handlers[type](e); return e; };
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

test("Read again on a lazy section reloads only it with fresh=1, shows it is busy, then says Updated and when it was read", async () => {
  const again = new El("a", { class: "ui-iconbtn", href: "/?pr_sort=newest&fresh=1", "data-ui-again": "", "data-ui-carry": "fresh" });
  const icon = new El("svg");
  again.appendChild(icon);
  const read = new El("time", { class: "ui-read", "data-ui-read": "", hidden: "" });
  read.hidden = true;
  const done = new El("span", { class: "ui-again-done", role: "status", "aria-live": "polite" });
  const tools = new El("span", { class: "ui-section-tools", "data-ui-tools": "" }, [read, done, again]);
  const head = new El("div", { class: "ui-section-head" }, [tools]);
  const box = new El("div", { class: "ui-lazy", "data-ui-src": "/part/focus/decisions?pr_sort=newest" }, [new El("section", { class: "ui-section" }, [head])]);
  El.parse["<loaded>"] = () => [new El("section", { class: "ui-section" }, [new El("div", { class: "ui-section-head" }), new El("p", { id: "new" })])];
  const p = page("http://localhost/?pr_sort=newest", [box], { "/part/focus/decisions?pr_sort=newest": "<loaded>", "/part/focus/decisions?pr_sort=newest&fresh=1": "<loaded>", readAt: "2026-10-29T09:41:00Z" });
  await new Promise((r) => setTimeout(r, 0));
  assert.equal(read.hidden, false, "the first load says when the sources were read");
  assert.equal(read.getAttribute("datetime"), "2026-10-29T09:41:00Z");
  assert.match(read.textContent, /^read \d/);
  const e = p.fire("click", icon);
  assert.ok(e.prevented);
  assert.equal(p.fetches.at(-1), "/part/focus/decisions?pr_sort=newest&fresh=1");
  assert.equal(box.getAttribute("data-ui-src"), "/part/focus/decisions?pr_sort=newest", "the next plain load reads from the cache again");
  assert.equal(box.querySelector("section").getAttribute("aria-busy"), "true");
  assert.ok(box.querySelector("section").hasAttribute("data-ui-reloading"), "the old section dims under a bar");
  assert.equal(again.getAttribute("aria-busy"), "true", "the button spins");
  assert.equal(p.fire("click", icon).prevented, true);
  assert.equal(p.fetches.length, 2, "a busy button asks nothing more");
  await new Promise((r) => setTimeout(r, 0));
  assert.ok(box.querySelector("[id=new]"), "the section shows what the server read");
  assert.equal(box.querySelector(".ui-section-head [data-ui-again]"), again);
  assert.equal(again.getAttribute("aria-busy"), null);
  assert.equal(done.textContent, "Updated");
});

test("a section link reloads only its section at the link's address, and carries it to Read again", () => {
  const refresh = new El("a", { href: "/roads?road=r6&fresh=1", "data-ui-carry": "fresh" });
  const label = new El("text");
  const link = new El("a", { href: "/roads?road=r6&open=kit%2Fcards", "data-ui-section-link": "" }, [label]);
  const box = new El("div", { class: "ui-lazy", "data-ui-src": "/part/roads/road-graph?road=r6" }, [new El("section", { class: "ui-section" }, [link])]);
  const p = page("http://localhost/roads?road=r6", [refresh, box]);
  const e = p.fire("click", label);
  assert.ok(e.prevented);
  assert.equal(p.location.search, "?road=r6&open=kit%2Fcards");
  assert.equal(p.fetches.at(-1), "/part/roads/road-graph?road=r6&open=kit%2Fcards");
  assert.equal(refresh.getAttribute("href"), "/roads?road=r6&open=kit%2Fcards&fresh=1");
});

// A POST form that records what it sends: requestSubmit fires submit again and sends unless it is stopped, as a browser does.
function postForm(attrs, children) {
  const form = new El("form", { method: "post", ...attrs }, children);
  form.method = "post";
  form.sent = [];
  form.requestSubmit = (submitter) => {
    const fields = form.querySelectorAll("input[type=hidden]").map((i) => i.getAttribute("name") + "=" + i.getAttribute("value"));
    if (!form.page.fire("submit", form, { submitter }).prevented) form.sent.push(submitter || fields);
  };
  return form;
}

// answer closes the kit's confirm dialog the way the viewer did: "confirm", "cancel", or "" for Escape.
function answer(dlg, value) {
  dlg.returnValue = value;
  dlg.open = false;
  dlg.onclose();
}

test("a button marked data-ui-confirm asks in the kit's dialog first, and Cancel or Escape sends nothing", () => {
  const btn = new El("button", { class: "ui-btn", type: "submit", name: "do", value: "disconnect", "data-variant": "danger", "data-ui-confirm": "Disconnect GitHub?" });
  btn.textContent = "Disconnect";
  const form = postForm({ action: "/settings" }, [btn]);
  const p = page("http://localhost/settings", [form]);
  form.page = p;
  assert.ok(p.fire("submit", form, { submitter: btn }).prevented, "the form waits for the answer");
  const dlg = form.parent.querySelector("dialog.ui-dialog");
  assert.ok(dlg.open);
  assert.equal(dlg.querySelector(".ui-dialog-title").textContent, "Disconnect GitHub?");
  const yes = dlg.querySelector("button[value='confirm']");
  assert.equal(yes.textContent, "Disconnect", "the confirm button repeats the verb");
  assert.equal(yes.getAttribute("data-variant"), "danger");
  assert.equal(El.focused, dlg.querySelector("button[value='cancel']"), "Cancel gets first focus");
  assert.equal(dlg.querySelector(".ui-check").hidden, true, "no Remember my answer unless asked for");
  answer(dlg, "cancel");
  p.fire("submit", form, { submitter: btn });
  answer(dlg, "");
  assert.deepEqual(form.sent, []);
});

test("confirm sends the form once, with the button's name and value; a form without the mark is sent at once", () => {
  const btn = new El("button", { type: "submit", name: "do", value: "disconnect", "data-ui-confirm": "Disconnect GitHub?" });
  const form = postForm({}, [btn]);
  const plain = postForm({}, []);
  const p = page("http://localhost/settings", [form, plain]);
  form.page = p;
  p.fire("submit", form, { submitter: btn });
  const dlg = form.parent.querySelector("dialog");
  answer(dlg, "confirm");
  assert.deepEqual(form.sent, [btn], "the browser sends the button's name and value with it");
  assert.equal(form.querySelector("input[name=confirm]"), null);
  assert.equal(p.fire("submit", form, { submitter: btn }).prevented, true, "the next press asks again");
  assert.equal(p.fire("submit", plain).prevented, undefined);
});

test("Remember my answer sends confirm=remember with the form, so the server can remember", () => {
  const form = postForm({ action: "/focus", "data-ui-confirm": "Focus on kanban?", "data-ui-confirm-label": "Focus", "data-ui-remember": "" }, []);
  const p = page("http://localhost/", [form]);
  form.page = p;
  p.fire("submit", form);
  const dlg = form.parent.querySelector("dialog");
  assert.equal(dlg.querySelector("button[value='confirm']").textContent, "Focus");
  const remember = dlg.querySelector(".ui-check");
  assert.equal(remember.hidden, false);
  remember.querySelector("input").checked = true;
  answer(dlg, "confirm");
  assert.deepEqual(form.sent, [["confirm=remember"]]);
  assert.equal(form.querySelector("input[name=confirm]"), null, "a form sent again later does not carry it");
});

test("a name to type keeps the confirm button off until the text is exactly the name", () => {
  const btn = new El("button", { type: "submit", "data-variant": "danger", "data-ui-confirm": "Delete card 20261102-e472?", "data-ui-confirm-name": "20261102-e472", "data-ui-remember": "" });
  btn.textContent = "Delete card";
  const form = postForm({}, [btn]);
  const p = page("http://localhost/", [form]);
  form.page = p;
  p.fire("submit", form, { submitter: btn });
  const dlg = form.parent.querySelector("dialog");
  const field = dlg.querySelector("[data-ui-name]");
  const typed = field.querySelector("input");
  const yes = dlg.querySelector("button[value='confirm']");
  assert.equal(field.hidden, false);
  assert.equal(field.querySelector("code").textContent, "20261102-e472", "it says Type 20261102-e472 to go ahead");
  assert.equal(El.focused, typed);
  assert.equal(dlg.querySelector(".ui-check").hidden, false, "Remember my answer is offered here too");
  assert.equal(yes.disabled, true);
  typed.value = "20261102-e47";
  typed.oninput();
  assert.equal(yes.disabled, true, "almost the name is not the name");
  typed.value = "20261102-e472";
  typed.oninput();
  assert.equal(yes.disabled, false);
  answer(dlg, "confirm");
  assert.deepEqual(form.sent, [btn]);
  p.fire("submit", form, { submitter: btn });
  assert.equal(typed.value, "", "the next question starts empty");
  assert.equal(yes.disabled, true);
});
