import { test } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { createElement as h } from "react";
import { renderToStaticMarkup as render } from "react-dom/server";
import { License, LicenseFeatures } from "./index.js";

// The same samples as sampleLicense in preview/cases.go, so React must draw each golden of the Go partial.
const covers = { events: "GET /api/events, the live stream of what changed", telemetry: "traces, metrics and logs sent to another system over OpenTelemetry" };
const lockedEvents = { id: "events", covers: covers.events, tier: "enterprise", on: false, by: "license", startsOff: true, from: "none" };
const brokenTelemetry = { id: "telemetry", covers: covers.telemetry, tier: "enterprise", on: true, startsOff: true, missing: [{ name: "OTEL_EXPORTER_OTLP_ENDPOINT", purpose: "where telemetry is sent" }], from: "org", until: "2027-11-19" };
const add = (label, primary) => [{ label, href: "#add-license", variant: primary ? "primary" : undefined }];
const files = {
  trial: [{ kind: "trial", for: "Ada Example", until: "2026-12-03", keyId: "prod-2026-a" }],
  "trial-ended": [{ kind: "trial", for: "Ada Example", until: "2026-12-03", ended: true, keyId: "prod-2026-a" }],
  licensed: [{ kind: "org", for: "Example Ltd", until: "2027-11-19", keyId: "prod-2026-a" }],
  ending: [{ kind: "org", for: "Example Ltd", until: "2026-12-01", keyId: "prod-2026-a" }],
  dev: [{ kind: "product", for: "Ada Example", until: "2027-01-01", keyId: "dev-1" }],
};
const unlock = { trial: ["trial", "2026-12-03"], licensed: ["org", "2027-11-19"], ending: ["org", "2026-12-01"], dev: ["product", "2027-01-01"] };
const acts = { none: add("Add a license"), "trial-ended": add("Add a license"), trial: add("Add a license", true), licensed: add("Replace the license"), ending: add("Add the renewed license", true), dev: add("Add a dev license") };

function sample(state) {
  const [from, until] = unlock[state] || [];
  let events = lockedEvents, telemetry = { id: "telemetry", covers: covers.telemetry, tier: "enterprise", on: false, by: "license", startsOff: true, from: "none" };
  if (from) {
    events = { id: "events", covers: covers.events, tier: "enterprise", on: true, startsOff: true, from, until };
    telemetry = { id: "telemetry", covers: covers.telemetry, tier: "enterprise", on: false, by: "default", startsOff: true, from, until };
  }
  if (state === "licensed") telemetry = brokenTelemetry;
  return {
    state, release: state !== "dev", version: state === "dev" ? "dev-3d91aef" : "v9.4.1", trial: "https://example.com/trial-key",
    licenses: files[state] || [], actions: acts[state],
    countdown: state === "trial" ? { left: 10, of: 14 } : state === "ending" ? { left: 12, of: 14 } : undefined,
    features: [events, telemetry,
      { id: "edit", covers: "apply, edit, connect, delete, focus, the design files and word lists", tier: "free", on: true },
      { id: "layout", covers: "moving, tidying and arranging boxes and lines", tier: "free", on: true },
      { id: "tracing", covers: "saved requests, runs and suites", tier: "free", on: false, by: "default", startsOff: true },
      { id: "workshop", covers: "reading a project’s code into documents", tier: "free", on: false, by: "switch" }],
  };
}

// same is markup with each tag's attributes in one order, since React and html/template write them in their own.
const same = (html) => html.replace(/\n/g, "").replace(/<([a-z]+)((?:\s+[^\s=>]+(?:="[^"]*")?)*)>/g, (_, tag, attrs) => `<${tag}${(attrs.match(/\s+[^\s=>]+(?:="[^"]*")?/g) || []).map((a) => a.trim()).sort().map((a) => " " + a).join("")}>`);
const golden = (name) => same(readFileSync(new URL(`../testdata/golden/license--${name}.html`, import.meta.url), "utf8"));

for (const state of ["none", "trial", "trial-ended", "licensed", "ending", "dev"]) {
  test(`license ${state} draws the same markup as the Go partial`, () => {
    assert.equal(same(render(h(License, sample(state)))), golden(state));
  });
}

test("a locked row and an on-but-broken row draw the same markup as the Go partial", () => {
  assert.equal(same(render(h(LicenseFeatures, { items: [lockedEvents] }))), golden("locked-row"));
  assert.equal(same(render(h(LicenseFeatures, { items: [brokenTelemetry] }))), golden("broken-row"));
});
