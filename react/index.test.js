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

test("count says its label on hover and to screen readers", () => {
  const labelled = render(h(kit.Count, { n: 2, tone: "warn", label: "2 cards wait for you" }));
  assert.match(labelled, /class="ui-count" data-tone="warn" title="2 cards wait for you"/);
  assert.match(labelled, /class="ui-sr">2 cards wait for you</);
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

test("preview banner draws the same markup as the Go partial, and nothing for a release", () => {
  const banner = render(h(kit.PreviewBanner, { version: "v0.4.0-3-g758ac48", branch: "feat/road-focus", commit: "758ac48", dirty: true }));
  assert.equal(banner, '<div class="ui-preview-banner" role="note" aria-label="Preview build" title="PREVIEW · branch feat/road-focus · 758ac48 · not a release · uncommitted changes · version v0.4.0-3-g758ac48"><span class="ui-preview-banner-text"><b>PREVIEW</b> · branch feat/road-focus · 758ac48 · not a release · uncommitted changes</span></div>');
  assert.equal(render(h(kit.PreviewBanner, { branch: "main", commit: "9a1b2c3", released: true })), "");
  const shell = render(h(kit.Shell, { build: { branch: "feat/x", commit: "abc1234" }, bar: { app: "kanban" } }));
  assert.match(shell, /^<div class="ui-preview-banner"[^]*<header class="ui-bar/);
});

test("light draws the same link, dot and hidden popover as the Go partial", () => {
  const light = render(h(kit.Light, { id: "gh", tone: "warn", icon: "github", label: "GitHub", href: "/settings/github", title: "Connected, the last read failed.", details: [{ label: "Account", value: "bot" }], actions: [{ label: "Disconnect", variant: "danger", confirm: { text: "x", label: "Disconnect now" } }] }));
  assert.match(light, /<a class="ui-light-btn" href="\/settings\/github" aria-controls="gh-pop" aria-expanded="false" aria-haspopup="dialog"/);
  assert.match(light, /<span class="ui-dot" data-tone="warn"><\/span>/);
  assert.match(light, /class="ui-pop" id="gh-pop" role="dialog" aria-labelledby="gh-pop-title" tabindex="-1" hidden=""/);
  assert.match(light, /<dt>Account<\/dt><dd>bot<\/dd>/);
  assert.match(light, /<summary class="ui-btn" data-size="sm">Actions<\/summary>/);
  assert.match(light, /data-variant="danger" aria-expanded="false">Disconnect<\/button>/);
});

test("count with an icon is a badge whose label is its tooltip, and with href a link", () => {
  const badge = render(h(kit.Count, { n: 6, tone: "good", icon: "check", label: "6 approved, not merged yet" }));
  assert.match(badge, /^<span class="ui-count" data-tone="good" data-icon="" data-tip="6 approved, not merged yet" tabindex="0"><svg/);
  assert.match(badge, /class="ui-sr">6 approved, not merged yet</);
  const link = render(h(kit.Count, { n: 2, icon: "eye", label: "2 to review", href: "#review" }));
  assert.match(link, /^<a class="ui-count" data-icon="" data-tip="2 to review" href="#review">/);
});

test("confirm form draws the form and, beside it, the kit's dialog with the question and Remember my answer", () => {
  const f = render(h(kit.ConfirmForm, { id: "ask-save", question: "Save the settings?", label: "Save", remember: true, method: "post", action: "/settings" }, h("button", { type: "submit" }, "Save")));
  assert.match(f, /^<form action="\/settings" method="post"><button type="submit">Save<\/button><\/form><dialog class="ui-dialog" id="ask-save"/);
  assert.match(f, /<h2 class="ui-dialog-title" id="ask-save-title">Save the settings\?<\/h2><label class="ui-check"><input type="checkbox"\/><span>Remember my answer<\/span><\/label><div class="ui-dialog-acts"><button class="ui-btn" type="submit" value="cancel">Cancel<\/button><button class="ui-btn" data-variant="primary" value="confirm" type="submit">Save<\/button>/);
  assert.doesNotMatch(render(h(kit.ConfirmForm, { id: "a", question: "Disconnect?" })), /ui-check/);
});
