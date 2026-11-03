// Tests ConfirmForm and ConfirmSettings in a browser-like DOM (jsdom): what the viewer does, and what the form sends.
import { test } from "node:test";
import assert from "node:assert/strict";
import { JSDOM } from "jsdom";

const dom = new JSDOM("<!doctype html><body></body>", { url: "http://localhost/" });
Object.assign(globalThis, { window: dom.window, document: dom.window.document, FormData: dom.window.FormData, IS_REACT_ACT_ENVIRONMENT: true });
globalThis.HTMLElement = dom.window.HTMLElement;
const { createElement: h, act } = await import("react");
const { createRoot } = await import("react-dom/client");
const kit = await import("./index.js");

// mount draws el into a fresh container and returns it with a cleanup.
async function mount(el) {
  const box = document.body.appendChild(document.createElement("div"));
  const root = createRoot(box);
  await act(async () => root.render(el));
  return { box, rerender: (next) => act(async () => root.render(next)), done: () => { root.unmount(); box.remove(); } };
}

// press submits form with button, as a click on it does.
const press = (form, button) => act(async () => form.requestSubmit(button));

// confirmForm is a ConfirmForm whose onSubmit records each send as its entries.
function confirmForm(props) {
  const sent = [];
  const el = h(kit.ConfirmForm, { id: "ask", question: "Delete card 7?", onSubmit: (data) => sent.push([...data.entries()]), ...props },
    h("button", { type: "submit", name: "do", value: "delete" }, "Delete card"));
  return { sent, el };
}

test("ConfirmForm waits for the answer: Cancel sends nothing, the verb sends once with confirm=yes", async () => {
  const { sent, el } = confirmForm();
  const m = await mount(el);
  const form = m.box.querySelector("form");
  await press(form, form.querySelector("button[name=do]"));
  assert.deepEqual(sent, [], "nothing is sent before the answer");
  const dlg = document.getElementById("ask");
  assert.equal(dlg.querySelector("button[value='confirm']").textContent, "Delete card");
  await press(dlg.querySelector("form"), dlg.querySelector("button[value='cancel']"));
  assert.deepEqual(sent, [], "Cancel sends nothing");
  await press(form, form.querySelector("button[name=do]"));
  await press(dlg.querySelector("form"), dlg.querySelector("button[value='confirm']"));
  assert.deepEqual(sent, [[["do", "delete"], ["confirm", "yes"]]]);
  assert.equal(form.querySelector("input[name=confirm]"), null, "the next send does not carry it");
  m.done();
});

test("ConfirmForm with Remember my answer ticked sends confirm=remember", async () => {
  const { sent, el } = confirmForm({ remember: true });
  const m = await mount(el);
  const form = m.box.querySelector("form");
  await press(form, form.querySelector("button[name=do]"));
  const dlg = document.getElementById("ask");
  await act(async () => dlg.querySelector(".ui-check input").click());
  await press(dlg.querySelector("form"), dlg.querySelector("button[value='confirm']"));
  assert.deepEqual(sent, [[["do", "delete"], ["confirm", "remember"]]]);
  m.done();
});

test("ConfirmSettings radios follow a changed app value", async () => {
  const m = await mount(h(kit.ConfirmSettings, { types: ["move"], app: {} }));
  const checked = () => m.box.querySelector("input[name='confirm.move']:checked").value;
  assert.equal(checked(), "default");
  await m.rerender(h(kit.ConfirmSettings, { types: ["move"], app: { move: "never" } }));
  assert.equal(checked(), "never");
  m.done();
});
