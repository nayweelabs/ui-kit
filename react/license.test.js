import { test } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { createElement as h } from "react";
import { renderToStaticMarkup as render } from "react-dom/server";
import { License, LicenseFeatures } from "./index.js";

const read = (name) => readFileSync(new URL(`../testdata/golden/${name}`, import.meta.url), "utf8");
// same is markup with each tag's attributes in one order, since React and html/template write them in their own.
const same = (html) => html.replace(/\n/g, "").replace(/<([a-z]+)((?:\s+[^\s=>]+(?:="[^"]*")?)*)>/g, (_, tag, attrs) => `<${tag}${(attrs.match(/\s+[^\s=>]+(?:="[^"]*")?/g) || []).map((a) => a.trim()).sort().map((a) => " " + a).join("")}>`);
const parts = { "ui/license": License, "ui/license-features": LicenseFeatures };

// The cases are the Go preview's own (written by go test -run Golden -update), so React must draw each Go golden.
for (const [name, { template, props }] of Object.entries(JSON.parse(read("license-cases.json")))) {
  test(`license ${name} draws the same markup as the Go partial`, () => {
    assert.equal(same(render(h(parts[template], props))), same(read(`license--${name}.html`)));
  });
}

test("license draws a state and a closing word it has no words for as given", () => {
  const html = render(h(License, { state: "public", build: { release: true }, features: [{ id: "share", tier: "enterprise", on: false, by: "public_mode", from: "none" }] }));
  assert.match(html, /<span class="ui-pill">public<\/span><b class="ui-license-title">public<\/b>/);
  assert.match(html, /<td>public_mode<\/td>/);
  assert.match(html, /<span class="ui-dot" data-tone="none"><\/span>Off/);
});

test("a button that goes to another site opens in a new tab, as the Go partial does", () => {
  assert.match(render(h(License, { state: "none", build: { release: true }, trial: "https://example.com/trial-key" })), /href="https:\/\/example.com\/trial-key" target="_blank" rel="noopener noreferrer"/);
});
