import { test } from "node:test";
import assert from "node:assert/strict";
import { createElement as h } from "react";
import { renderToStaticMarkup as render } from "react-dom/server";
import * as kit from "./index.js";

test("pill draws the same classes and tone as the Go partial", () => {
  assert.equal(render(h(kit.Pill, { tone: "warn", label: "3 behind" })), '<span class="ui-pill" data-tone="warn">3 behind</span>');
});

test("count hides zero, shows 99+ with the true number for screen readers", () => {
  assert.equal(render(h(kit.Count, { n: 0 })), "");
  const many = render(h(kit.Count, { n: 128, tone: "accent" }));
  assert.match(many, /aria-hidden="true">99\+</);
  assert.match(many, /class="ui-sr">128</);
});

test("status dot always carries its word", () => {
  assert.equal(render(h(kit.Status, { tone: "good", label: "current" })), '<span class="ui-status"><span class="ui-dot" data-tone="good"></span>current</span>');
});

test("product card shows the star for focus and a failed fact in place", () => {
  const card = render(h(kit.ProductCard, { name: "designer", href: "#", focus: true, problem: { title: "Can’t read the board" }, facts: [{ label: "Board", error: "can’t read the board" }] }));
  assert.match(card, /class="ui-card" data-focus="" data-state="failed"/);
  assert.match(card, /class="ui-star" data-chip=""/);
  assert.match(card, /role="alert"/);
  assert.match(card, /class="ui-fact-err"/);
});

test("products draws the list view with the same facts as columns", () => {
  const list = render(h(kit.Products, { view: "list", columns: ["Installed"], items: [{ name: "hubman", href: "#", facts: [{ label: "Installed", text: "current" }] }] }));
  assert.match(list, /class="ui-list"/);
  assert.match(list, /<th scope="col">Installed<\/th>/);
  assert.doesNotMatch(list, /ui-cards/);
});

test("PR item names the repo and number, and links the files count", () => {
  const pr = render(h(kit.PrItem, { repo: "checkman", number: "5", href: "#pr", title: "t", files: { text: "3", href: "#files" }, added: "48", removed: "12" }));
  assert.match(pr, /class="ui-repo" href="#pr">checkman <b>#5<\/b>/);
  assert.match(pr, /<a class="ui-num" href="#files">3<\/a> files/);
  assert.match(pr, /\+48.*−12/);
});

test("section in each of its states", () => {
  assert.match(render(h(kit.Section, { title: "PRs", loading: true })), /aria-busy="true".*ui-skel-group/);
  assert.match(render(h(kit.Section, { title: "PRs", problem: { title: "Can’t read" } })), /role="alert"/);
  assert.match(render(h(kit.Section, { title: "PRs", empty: { title: "No open pull requests." } })), /ui-empty/);
});

test("secret field is a password with Show, and saved shows only who is connected", () => {
  const f = render(h(kit.Field, { id: "t", label: "Token", secret: true }));
  assert.match(f, /type="password"/);
  assert.match(f, /autoComplete|autocomplete="off"/);
  assert.match(f, />Show</);
  assert.doesNotMatch(render(h(kit.Field, { id: "t", label: "Token", saved: "Connected as bot" })), /<input/);
});

test("breadcrumb's last item is the current page, not a link", () => {
  const c = render(h(kit.Crumb, { items: [{ label: "Products", href: "/p" }, { label: "kanban", href: "/k" }] }));
  assert.match(c, /<span aria-current="page">kanban<\/span>/);
});

test("shell draws bar, header, body and footer", () => {
  const s = render(h(kit.Shell, { bar: { app: "Control room", pages: [{ label: "Focus", href: "/", current: true }] }, header: { title: "Focus" }, footer: { version: "v0.4.0" } }, "body"));
  for (const want of ["ui-bar", 'aria-current="page"', "<h1>Focus</h1>", "body", "ui-footer"]) assert.ok(s.includes(want), want);
});
