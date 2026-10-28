// Thin React wrappers over the kit's classes. Plain ES modules with createElement, so consumers need no build step;
// types are in index.d.ts. Each wrapper draws the same markup as the Go partial of the same name; kit.css does the rest.
import { Component, createElement as h, Fragment, useEffect, useRef, useState } from "react";

const cx = (...names) => names.filter(Boolean).join(" ");
const tone = (t) => (t ? t : undefined);

const paths = {
  star: [["path", { d: "M12 2.6l2.9 6 6.5.8-4.8 4.5 1.2 6.5L12 17.3l-5.8 3.1 1.2-6.5L2.6 9.4l6.5-.8z" }]],
  alert: [["path", { d: "M8 1.8l6.6 11.6H1.4z", strokeWidth: 1.4 }], ["path", { d: "M8 6.2v3.4M8 11.4v.2", strokeWidth: 1.6 }]],
  check: [["circle", { cx: 8, cy: 8, r: 6.3, strokeWidth: 1.4 }], ["path", { d: "M5.2 8.2l1.9 1.9 3.8-4", strokeWidth: 1.5 }]],
  info: [["circle", { cx: 8, cy: 8, r: 6.3, strokeWidth: 1.4 }], ["path", { d: "M8 7.2v4M8 4.8v.2", strokeWidth: 1.6 }]],
  search: [["circle", { cx: 7, cy: 7, r: 4.5, strokeWidth: 1.5 }], ["path", { d: "M10.5 10.5L14 14", strokeWidth: 1.5 }]],
  cards: [["rect", { x: 1.5, y: 1.5, width: 5.5, height: 5.5, rx: 1 }], ["rect", { x: 9, y: 1.5, width: 5.5, height: 5.5, rx: 1 }], ["rect", { x: 1.5, y: 9, width: 5.5, height: 5.5, rx: 1 }], ["rect", { x: 9, y: 9, width: 5.5, height: 5.5, rx: 1 }]],
  list: [["path", { d: "M2 3.5h12M2 8h12M2 12.5h12", strokeWidth: 1.5 }]],
  theme: [["circle", { cx: 8, cy: 8, r: 5.6, strokeWidth: 1.4 }], ["path", { d: "M8 2.4a5.6 5.6 0 0 1 0 11.2z", fill: "currentColor", stroke: "none" }]],
  close: [["path", { d: "M4 4l8 8M12 4l-8 8", strokeWidth: 1.6 }]],
  github: [["path", { d: "M8 0c4.42 0 8 3.58 8 8a8.01 8.01 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38 0-.27.01-1.13.01-2.2 0-.75-.25-1.23-.54-1.48 1.78-.2 3.65-.88 3.65-3.95 0-.88-.31-1.59-.82-2.15.08-.2.36-1.02-.08-2.12 0 0-.67-.22-2.2.82-.64-.18-1.32-.27-2-.27-.68 0-1.36.09-2 .27-1.53-1.03-2.2-.82-2.2-.82-.44 1.1-.16 1.92-.08 2.12-.51.56-.82 1.28-.82 2.15 0 3.06 1.86 3.75 3.64 3.95-.23.2-.44.55-.51 1.07-.46.21-1.61.55-2.33-.66-.15-.24-.6-.83-1.23-.82-.67.01-.27.38.01.53.34.19.73.9.82 1.13.16.45.68 1.31 2.69.94 0 .67.01 1.3.01 1.49 0 .21-.15.45-.55.38A8 8 0 0 1 0 8c0-4.42 3.58-8 8-8z", fill: "currentColor", stroke: "none" }]],
};

export function Icon({ name }) {
  if (name === "star") return h("svg", { viewBox: "0 0 24 24", "aria-hidden": true, focusable: "false" }, h("path", paths.star[0][1]));
  const parts = paths[name] || [];
  return h("svg", { viewBox: "0 0 16 16", "aria-hidden": true, focusable: "false", fill: "none", stroke: "currentColor", strokeLinecap: "round", strokeLinejoin: "round", strokeWidth: 1.4 },
    parts.map(([tag, props], i) => h(tag, { key: i, ...props })));
}

export function Status({ tone: t, label }) {
  return h("span", { className: "ui-status" }, h("span", { className: "ui-dot", "data-tone": t || "none" }), label);
}

export function Pill({ tone: t, variant, label, title, href, children }) {
  const props = { className: "ui-pill", "data-tone": tone(t), "data-variant": variant || undefined, title };
  return href ? h("a", { ...props, href }, label ?? children) : h("span", props, label ?? children);
}

export function Count({ tone: t, n, label, showZero }) {
  if (!n && !showZero) return null;
  const many = n > 99;
  return h("span", { className: "ui-count", "data-tone": tone(t), title: label || undefined },
    label || many ? h(Fragment, null, h("span", { "aria-hidden": true }, many ? "99+" : String(n)), h("span", { className: "ui-sr" }, label || String(n))) : String(n));
}

export function Star({ label = "Focus now", text, chip, corner }) {
  return h("span", { className: "ui-star", "data-chip": chip ? "" : undefined, "data-corner": corner ? "" : undefined, role: "img", "aria-label": label, title: label },
    h(Icon, { name: "star" }), text ? h("span", { "aria-hidden": true }, text) : null);
}

export function Num({ text, href, command }) {
  return href ? h("a", { className: "ui-num", href, title: command }, text) : h("span", { className: "ui-num", title: command }, text);
}

export function Button({ label, variant, size, type = "button", href, kbd, disabled, busy, onClick, children, ...rest }) {
  const props = { className: "ui-btn", "data-variant": variant || undefined, "data-size": size || undefined, onClick, ...rest };
  const body = [label ?? children, kbd ? h(Fragment, { key: "k" }, " ", h("kbd", null, kbd)) : null];
  if (href) return h("a", { ...props, href, "aria-disabled": disabled || busy ? "true" : undefined, "aria-busy": busy ? "true" : undefined }, ...body);
  return h("button", { ...props, type, disabled: disabled || busy, "aria-busy": busy ? "true" : undefined, "aria-keyshortcuts": kbd }, ...body);
}

export function ErrorBox({ title, cause, command }) {
  return h("div", { className: "ui-box", "data-tone": "bad", role: "alert" },
    h(Icon, { name: "alert" }), h("b", { className: "ui-box-title" }, title),
    cause ? h("p", null, cause) : null, command ? h("code", { className: "ui-cmd" }, command) : null);
}

export function Notice({ tone: t = "info", title, text, link }) {
  const icon = t === "good" ? "check" : t === "warn" ? "alert" : "info";
  return h("div", { className: "ui-box", "data-tone": t, role: "status" },
    h(Icon, { name: icon }), h("b", { className: "ui-box-title" }, title),
    text || link ? h("p", null, text, link ? h(Fragment, null, " ", h("a", { href: link.href }, link.label)) : null) : null);
}

export function Empty({ title, text, command, action, small }) {
  if (small) return h("div", { className: "ui-empty", "data-size": "sm" }, title || "empty");
  return h("div", { className: "ui-empty" }, h("b", { className: "ui-empty-title" }, title),
    text ? h("span", null, text) : null, command ? h("code", null, command) : null, action ? h(Button, action) : null);
}

export function Skeleton({ shape = "card", lines = 3 }) {
  return h("div", { className: "ui-skel-group", "data-shape": shape, role: "status", "aria-label": "Loading" },
    h("span", { className: "ui-skel", "data-size": "title" }), Array.from({ length: lines }, (_, i) => h("span", { key: i, className: "ui-skel" })));
}

function Fact({ fact }) {
  if (fact.error) return h("span", { className: "ui-fact-err" }, h(Icon, { name: "alert" }), fact.error);
  return h(Fragment, null,
    fact.tone ? h("span", { className: "ui-dot", "data-tone": fact.tone }) : null,
    fact.num ? h(Num, fact.num) : null,
    fact.text ? " " + fact.text : null,
    fact.branch ? h(Fragment, null, " ", h("span", { className: "ui-branch" }, h("b", null, fact.branch))) : null,
    fact.pill ? h(Fragment, null, " ", h(Pill, fact.pill)) : null);
}

export function ProductCard({ name, href, kind, variant, focus, marker, waiting, problem, facts = [], foot, openHref }) {
  return h("article", { className: "ui-card", "data-focus": focus ? "" : undefined, "data-state": problem ? "failed" : undefined, "data-variant": variant || undefined },
    h("div", { className: "ui-card-head" },
      h("div", { className: "ui-card-title" }, h("a", { href }, name), kind ? h("span", { className: "ui-card-kind" }, kind) : null),
      focus ? h(Star, { chip: true, text: "focus" }) : marker ? h(Pill, marker) : waiting ? h(Count, waiting) : null),
    problem ? h(ErrorBox, problem) : null,
    facts.length ? h("dl", { className: "ui-kv" }, facts.map((f) => h(Fragment, { key: f.label }, h("dt", null, f.label), h("dd", null, h(Fact, { fact: f }))))) : null,
    h("div", { className: "ui-card-foot" }, h("span", null, foot), h("a", { href: openHref || href }, "Open", h("span", { className: "ui-sr" }, " " + name))));
}

export function ProductRow({ name, href, kind, focus, facts = [] }) {
  return h("tr", null,
    h("td", null, h("a", { className: "ui-item", href }, name), kind ? h("span", { className: "ui-item-sub" }, kind) : null),
    facts.map((f) => h("td", { key: f.label }, h("span", { className: "ui-cell" }, h(Fact, { fact: f })))),
    h("td", null, focus ? h(Star, {}) : null));
}

// The products view: the app keeps the view in state (and in ?view=); the server sends the items in order.
export function Products({ view = "cards", onView, toolbar, columns = [], items = [], empty, problem, loading }) {
  let body;
  if (loading) body = h("div", { className: "ui-cards" }, [0, 1, 2].map((i) => h(Skeleton, { key: i })));
  else if (problem) body = h(ErrorBox, problem);
  else if (!items.length) body = empty ? h(Empty, empty) : null;
  else if (view === "list") body = h("div", { className: "ui-list" }, h("table", null,
    h("thead", null, h("tr", null, h("th", { scope: "col" }, "Product"), columns.map((c) => h("th", { scope: "col", key: c }, c)), h("th", { scope: "col" }, h("span", { className: "ui-sr" }, "Focus")))),
    h("tbody", null, items.map((p) => h(ProductRow, { key: p.name, ...p })))));
  else body = h("div", { className: "ui-cards" }, items.map((p) => h(ProductCard, { key: p.name, ...p })));
  return h("div", { className: "ui-products" }, toolbar ? h(Toolbar, { ...toolbar, view, onView }) : null, body);
}

export function PrItem(pr) {
  const { repo, number, href, repoHref, author, initials, age, draft, title, verb, branch, base = "main", files, added, removed, conflict, review, ci, open = "Open on GitHub", compact, problem } = pr;
  if (problem) return h("div", { className: "ui-pr" }, h(ErrorBox, problem));
  return h("div", { className: "ui-pr", "data-compact": compact ? "" : undefined },
    h("div", { className: "ui-pr-top" },
      compact ? h(Fragment, null, h("span", { className: "ui-pr-verb" }, verb), " ") : null,
      h("a", { className: "ui-repo", href: href || repoHref }, repo, " ", h("b", null, "#" + number)),
      draft ? h("span", { className: "ui-pill", "data-variant": "tag" }, "draft") : null,
      !compact && initials ? h("span", { className: "ui-avatar", "aria-hidden": true }, initials) : null,
      !compact && author ? h("span", null, author) : null,
      age ? h("span", null, "· " + age) : null),
    h("a", { className: "ui-pr-title", href }, title),
    compact ? null : h("div", { className: "ui-pr-meta" },
      branch ? h("span", { className: "ui-branch" }, h("b", null, branch), " → " + base) : null,
      files ? h("span", null, h(Num, files), " files") : null,
      added || removed ? h("span", null, h("span", { className: "ui-add" }, "+" + (added || "0")), " ", h("span", { className: "ui-del" }, "−" + (removed || "0"))) : null,
      conflict ? h("span", { className: "ui-status", "data-tone": "warn" }, h("span", { className: "ui-dot", "data-tone": "warn" }), conflict) : null),
    h("div", { className: "ui-pr-side" },
      h("span", { className: "ui-pr-pills" }, review ? h(Pill, review) : null, ci ? h(Pill, ci) : null),
      compact ? null : h("a", { className: "ui-btn", "data-size": "sm", href }, open)));
}

export function PrList({ items = [], empty, loading }) {
  if (loading) return h(Skeleton, { shape: "rows", lines: 6 });
  if (!items.length) return empty ? h(Empty, empty) : null;
  return h("div", { className: "ui-prs" }, items.map((p, i) => h(PrItem, { key: p.repo + "#" + (p.number || i), ...p })));
}

const pct = (f) => (f * 100).toFixed(2) + "%";

export function Milestones({ label, total, track = [], items = [], empty }) {
  return h("div", { className: "ui-miles" },
    h("span", { className: "ui-label" }, label, total ? h(Fragment, null, " · ", h("span", { className: "ui-num" }, total)) : null),
    !items.length ? (empty ? h(Empty, empty) : null) : h(Fragment, null,
      h("div", { className: "ui-track", "aria-hidden": true }, track.map((m, i) => m.kind === "span"
        ? h("span", { key: i, className: "ui-track-span", style: { left: pct(m.at), width: pct(m.width) } })
        : h("span", { key: i, className: "ui-track-mark", "data-kind": m.kind || "point", style: { left: pct(m.at) } }))),
      h("ol", { className: "ui-milelist" }, items.map((m) => h("li", { className: "ui-mile", key: m.name },
        h("b", null, m.name, m.focus ? h(Star, {}) : null), h("span", null, m.when, m.days ? " · " + m.days : ""))))));
}

// Section fails on its own: an error thrown while drawing its children shows its error box, not a blank page.
class Boundary extends Component {
  constructor(props) { super(props); this.state = { failed: null }; }
  static getDerivedStateFromError(error) { return { failed: error }; }
  render() {
    if (!this.state.failed) return this.props.children;
    return h(ErrorBox, { title: this.props.failTitle || "Can’t draw this section.", cause: String(this.state.failed.message || this.state.failed) });
  }
}

export function Section({ id, title, count, why, problem, empty, loading, more, children, failTitle }) {
  return h("section", { className: "ui-section", id, "aria-labelledby": id ? id + "-title" : undefined, "aria-busy": loading ? "true" : undefined },
    h("div", { className: "ui-section-head" },
      h("h2", { className: "ui-section-title", id: id ? id + "-title" : undefined }, title, count ? h(Fragment, null, " ", h(Count, count)) : null),
      why ? h("span", { className: "ui-why" }, why) : null),
    loading ? h(Skeleton, { shape: "rows", lines: 3 }) : h(Fragment, null,
      problem ? h(ErrorBox, problem) : null,
      children ? h(Boundary, { failTitle }, children) : empty && !problem ? h(Empty, empty) : null,
      more ? h("a", { className: "ui-section-more", href: more.href }, more.label) : null));
}

function Cell({ cell }) {
  return h(Fragment, null,
    cell.href ? h("a", { className: "ui-item", href: cell.href }, cell.text) : cell.id ? h("button", { className: "ui-id", type: "button", "data-ui-copy": cell.id, title: "Copy " + cell.id, onClick: () => navigator.clipboard && navigator.clipboard.writeText(cell.id) }, cell.id) : cell.text,
    cell.sub ? h("span", { className: "ui-item-sub" }, cell.sub) : null,
    cell.status ? h(Status, cell.status) : null, cell.pill ? h(Pill, cell.pill) : null, cell.star ? h(Star, cell.star) : null);
}

export function Table({ label, columns = [], rows = [], empty, loading, small }) {
  return h("div", { className: "ui-list", "data-size": small ? "sm" : undefined, role: label ? "region" : undefined, "aria-label": label, tabIndex: label ? 0 : undefined },
    loading ? h(Skeleton, { shape: "rows", lines: 4 }) : h(Fragment, null,
      h("table", null,
        h("thead", null, h("tr", null, columns.map((c) => h("th", { key: c.label, scope: "col", "data-num": c.num ? "" : undefined }, c.hidden ? h("span", { className: "ui-sr" }, c.label) : c.label)))),
        h("tbody", null, rows.map((r, i) => h("tr", { key: i }, r.error
          ? h("td", { colSpan: columns.length, "data-error": "" }, h("span", { className: "ui-fact-err" }, h(Icon, { name: "alert" }), r.error))
          : r.cells.map((c, j) => h("td", { key: j, "data-num": c.num ? "" : undefined }, h(Cell, { cell: c }))))))),
      !rows.length && empty ? h(Empty, empty) : null));
}

export function Ticket({ title, href, summary, tags = [], turn, age, id, focus }) {
  return h("article", { className: "ui-ticket", "data-focus": focus ? "" : undefined },
    h("div", { className: "ui-ticket-head" }, href ? h("a", { className: "ui-ticket-title", href }, title) : h("b", { className: "ui-ticket-title" }, title), focus ? h(Star, {}) : null),
    summary ? h("span", { className: "ui-ticket-summary" }, summary) : null,
    tags.length ? h("div", { className: "ui-ticket-tags" }, tags.map((t) => h(Pill, { key: t.label, variant: "tag", ...t }))) : null,
    h("div", { className: "ui-ticket-foot" },
      h("span", null, turn ? h(Fragment, null, "turn ", h("b", null, turn)) : null, age ? h(Fragment, null, " · ", h("span", { className: "ui-num" }, age)) : null),
      id ? h("button", { className: "ui-id", type: "button", title: "Copy " + id, onClick: () => navigator.clipboard && navigator.clipboard.writeText(id) }, id) : null));
}

export function Field({ id, name, label, type = "text", value, defaultValue, onChange, placeholder, hint, error, secret, saved, savedAction }) {
  const [shown, setShown] = useState(false);
  const described = error ? id + "-err" : hint ? id + "-hint" : undefined;
  const input = { id, name: name || id, placeholder, "aria-invalid": error ? "true" : undefined, "aria-describedby": described };
  let control;
  if (saved) control = h("div", { className: "ui-saved" }, h(Status, { tone: "good", label: saved }), savedAction ? h(Button, savedAction) : null);
  else if (secret) control = h("div", { className: "ui-input", "data-invalid": error ? "" : undefined },
    h("input", { ...input, type: shown ? "text" : "password", "data-secret": "", autoComplete: "off", spellCheck: false, autoCapitalize: "off", onChange }),
    h("button", { className: "ui-btn", "data-variant": "quiet", "data-size": "sm", type: "button", "aria-controls": id, "aria-pressed": String(shown), onClick: () => setShown(!shown) }, shown ? "Hide" : "Show"));
  else control = h("div", { className: "ui-input", "data-invalid": error ? "" : undefined }, h("input", { ...input, type, value, defaultValue, onChange }));
  return h("div", { className: "ui-field" }, h("label", { htmlFor: id }, label), control,
    error ? h("span", { className: "ui-err", id: id + "-err" }, error) : hint ? h("span", { className: "ui-hint", id: id + "-hint" }, hint) : null);
}

export function ViewToggle({ view = "cards", onView }) {
  const btn = (v, icon, label) => h("button", { type: "button", "data-view": v, "aria-pressed": String(view === v), onClick: () => onView && onView(v) }, h(Icon, { name: icon }), label);
  return h("div", { className: "ui-seg", role: "group", "aria-label": "View" }, btn("cards", "cards", "Cards"), btn("list", "list", "List"));
}

export function Toolbar({ search, filters = [], clear, count, view, onView, primary, onSubmit }) {
  return h("form", { className: "ui-toolbar", role: "search", onSubmit: (e) => { e.preventDefault(); onSubmit && onSubmit(new FormData(e.currentTarget)); } },
    search ? h("label", { className: "ui-search" }, h(Icon, { name: "search" }), h("input", { type: "search", name: search.name || "q", defaultValue: search.value, placeholder: search.placeholder, "aria-label": search.label || search.placeholder, onChange: search.onChange })) : null,
    filters.map((f) => h("select", { key: f.name, className: "ui-select", name: f.name, "aria-label": f.label, "data-active": f.active ? "" : undefined, defaultValue: (f.options.find((o) => o.selected) || {}).value, onChange: f.onChange },
      f.options.map((o) => h("option", { key: o.value, value: o.value }, o.label)))),
    clear ? h(Button, { variant: "quiet", size: "sm", label: clear.label || "Clear", href: clear.href, onClick: clear.onClick }) : null,
    h("span", { className: "ui-grow" }),
    count ? h("span", { className: "ui-toolbar-count" }, count) : null,
    onView ? h(ViewToggle, { view, onView }) : null,
    primary ? h(Button, primary) : null);
}

export function Tabs({ label, items = [], onSelect }) {
  return h("div", { className: "ui-tabs", role: "tablist", "aria-label": label },
    items.map((t) => t.href
      ? h("a", { key: t.label, className: "ui-tab", role: "tab", href: t.href, "aria-selected": String(!!t.selected), "aria-current": t.selected ? "page" : undefined, tabIndex: t.selected ? undefined : -1 }, t.label, t.count ? h(Fragment, null, " ", h(Count, t.count)) : null)
      : h("button", { key: t.label, className: "ui-tab", role: "tab", type: "button", "aria-controls": t.panel, "aria-selected": String(!!t.selected), tabIndex: t.selected ? undefined : -1, onClick: () => onSelect && onSelect(t) }, t.label, t.count ? h(Fragment, null, " ", h(Count, t.count)) : null)));
}

export function Crumb({ items = [] }) {
  return h("nav", { "aria-label": "Breadcrumb" }, h("ol", { className: "ui-crumb" }, items.map((l, i) =>
    h("li", { key: i }, l.current || i === items.length - 1 ? h("span", { "aria-current": "page" }, l.label) : h("a", { href: l.href }, l.label)))));
}

// Dialog and Drawer use the native <dialog>: it traps focus and closes on Escape. `open` shows it as a modal.
function useModal(open, onClose) {
  return (el) => {
    if (!el) return;
    if (open && !el.open && el.showModal) el.showModal();
    if (!open && el.open) el.close();
    el.onclose = onClose || null;
  };
}

export function Dialog({ id, title, text, cancel = "Cancel", confirm, open, onClose, onConfirm }) {
  return h("dialog", { className: "ui-dialog", id, "aria-labelledby": id + "-title", ref: useModal(open, onClose) },
    h("form", { method: "dialog", onSubmit: (e) => { if (e.nativeEvent.submitter && e.nativeEvent.submitter.value === "confirm" && onConfirm) onConfirm(); } },
      h("h2", { className: "ui-dialog-title", id: id + "-title" }, title), text ? h("p", { className: "ui-dialog-text" }, text) : null,
      h("div", { className: "ui-dialog-acts" }, h("button", { className: "ui-btn", type: "submit", value: "cancel" }, cancel), h(Button, { type: "submit", value: "confirm", ...confirm }))));
}

export function Drawer({ id, title, open, onClose, children }) {
  return h("dialog", { className: "ui-drawer", id, "aria-labelledby": id + "-title", ref: useModal(open, onClose) },
    h("div", { className: "ui-drawer-head" }, h("h2", { className: "ui-dialog-title", id: id + "-title" }, title),
      h("form", { method: "dialog" }, h("button", { className: "ui-iconbtn", type: "submit", "aria-label": "Close" }, h(Icon, { name: "close" })))),
    children);
}

// Light is a status light for the top bar: a dot, an icon and a word that opens a popover with the facts and an
// Actions menu. Escape, a click outside or focus leaving closes it. An action links (href), runs onSelect, opens
// the dialog of id opens, or asks once more inside the popover (confirm: a form to action, or onConfirm).
export function Light({ id, tone: t, icon, label, href = "#", title, text, details = [], actions = [], open: initial = false, children }) {
  const [open, setOpen] = useState(initial);
  const [step, setStep] = useState(-1);
  const box = useRef(null);
  const btn = useRef(null);
  const close = (back) => { setOpen(false); setStep(-1); if (back && btn.current) btn.current.focus(); };
  useEffect(() => {
    if (!open) return undefined;
    const away = (e) => { if (box.current && !box.current.contains(e.target) && !e.target.closest("dialog")) close(false); };
    document.addEventListener("click", away);
    document.addEventListener("focusin", away);
    return () => { document.removeEventListener("click", away); document.removeEventListener("focusin", away); };
  }, [open]);
  const pick = (a, i) => {
    if (a.confirm) return setStep(i);
    if (a.opens) { const d = document.getElementById(a.opens); if (d && d.showModal && !d.open) d.showModal(); }
    if (a.onSelect) a.onSelect();
  };
  const item = (a, i) => a.href
    ? h("a", { key: i, className: "ui-menu-item", href: a.href, "data-variant": a.variant }, a.label)
    : h("button", { key: i, className: "ui-menu-item", type: "button", "data-variant": a.variant, "aria-haspopup": a.opens ? "dialog" : undefined, "aria-expanded": a.confirm ? String(step === i) : undefined, onClick: () => pick(a, i) }, a.label);
  const c = step >= 0 ? actions[step].confirm : null;
  return h("div", { className: "ui-light", id, ref: box, onKeyDown: (e) => { if (e.key === "Escape" && open && !e.target.closest("dialog")) { e.preventDefault(); close(true); } } },
    h("a", { className: "ui-light-btn", href, ref: btn, "aria-controls": id + "-pop", "aria-expanded": String(open), "aria-haspopup": "dialog", title, onClick: (e) => { e.preventDefault(); open ? close(true) : setOpen(true); } },
      h("span", { className: "ui-dot", "data-tone": t || "none" }), icon ? h(Icon, { name: icon }) : null, h("span", { className: "ui-light-word" }, label), title ? h("span", { className: "ui-sr" }, ": " + title) : null),
    h("div", { className: "ui-pop", id: id + "-pop", role: "dialog", "aria-labelledby": id + "-pop-title", tabIndex: -1, hidden: !open, ref: (el) => { if (el && open && !el.contains(document.activeElement)) el.focus(); } },
      h("b", { className: "ui-pop-title", id: id + "-pop-title" }, title || label),
      text ? h("p", { className: "ui-pop-text" }, text) : null,
      details.length ? h("dl", { className: "ui-kv" }, details.map((d) => h(Fragment, { key: d.label }, h("dt", null, d.label), h("dd", null, d.value)))) : null,
      actions.length ? h("details", { className: "ui-menu" }, h("summary", { className: "ui-btn", "data-size": "sm" }, "Actions"), h("div", { className: "ui-menu-list" }, actions.map(item))) : null,
      c ? h("form", { className: "ui-pop-step", method: c.action ? "post" : undefined, action: c.action, onSubmit: c.onConfirm ? (e) => { e.preventDefault(); c.onConfirm(); close(true); } : undefined },
        h("p", null, c.text), (c.hidden || []).map((f) => h("input", { key: f.name, type: "hidden", name: f.name, value: f.value })),
        h("div", { className: "ui-dialog-acts" }, h("button", { className: "ui-btn", type: "button", "data-size": "sm", autoFocus: true, onClick: () => setStep(-1) }, "Cancel"),
          h("button", { className: "ui-btn", type: "submit", "data-variant": actions[step].variant, "data-size": "sm" }, c.label))) : null),
    children);
}

// The shell. Theme follows the system; ThemeToggle switches and remembers it, the same key kit.js uses.
export function ThemeToggle() {
  const flip = () => {
    const root = document.documentElement;
    let cur = root.getAttribute("data-theme");
    if (!cur) cur = matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
    const next = cur === "dark" ? "light" : "dark";
    root.setAttribute("data-theme", next);
    try { localStorage.setItem("ui-kit.theme", next); } catch (e) { /* storage blocked: the switch still works for this page */ }
  };
  return h("button", { className: "ui-iconbtn", type: "button", "aria-label": "Switch light or dark theme", title: "Switch light or dark theme", onClick: flip }, h(Icon, { name: "theme" }));
}

export function Bar({ app, sub, href = "/", pages = [], apps, status = [], lights = [], noTheme, children }) {
  return h("header", { className: "ui-bar ui-gutter" }, h("div", { className: "ui-wrap" },
    apps && apps.length
      ? h("details", { className: "ui-apps" }, h("summary", { className: "ui-brand" }, app, sub ? h(Fragment, null, " ", h("small", null, sub)) : null),
        h("nav", { className: "ui-apps-menu", "aria-label": "Apps" }, apps.map((a) => h("a", { key: a.href, href: a.href, "aria-current": a.current ? "page" : undefined }, a.label))))
      : h("a", { className: "ui-brand", href }, app, sub ? h(Fragment, null, " ", h("small", null, sub)) : null),
    pages.length ? h("nav", { className: "ui-nav", "aria-label": "Pages" }, pages.map((p) => h("a", { key: p.href, href: p.href, "aria-current": p.current ? "page" : undefined }, p.label, p.count ? h(Fragment, null, " ", h(Count, p.count)) : null))) : null,
    h("div", { className: "ui-bar-end" },
      status.map((s) => h("span", { key: s.label, className: "ui-status", title: s.label }, h("span", { className: "ui-dot", "data-tone": s.tone || "none" }), h("span", { className: "ui-status-word" }, s.label))),
      lights.map((l) => h(Light, { key: l.id, ...l })), children, noTheme ? null : h(ThemeToggle))));
}

export function PageHeader({ crumb, title, lede, actions = [] }) {
  return h("div", { className: "ui-head" },
    crumb && crumb.length ? h(Crumb, { items: crumb }) : null,
    h("div", { className: "ui-head-text" }, h("h1", null, title), lede ? h("p", { className: "ui-lede" }, lede) : null),
    actions.length ? h("div", { className: "ui-actions" }, actions.map((a) => h(Button, { key: a.label, ...a }))) : null);
}

export function Footer({ version, versionHref, branch, fresh, links = [] }) {
  return h("footer", { className: "ui-footer ui-gutter" }, h("div", { className: "ui-wrap" },
    version ? (versionHref ? h("a", { className: "ui-num", href: versionHref }, version) : h("span", { className: "ui-num" }, version)) : null,
    branch ? h("span", { className: "ui-branch" }, "on ", h("b", null, branch)) : null,
    fresh && fresh.label ? h("span", { className: "ui-status" }, h("span", { className: "ui-dot", "data-tone": fresh.tone || "none" }), fresh.label) : null,
    links.length ? h("span", { className: "ui-footer-end" }, links.map((l) => h("a", { key: l.href, href: l.href }, l.label))) : null));
}

export function Rail({ groups = [] }) {
  return h("nav", { className: "ui-rail", "aria-label": "All pages" }, groups.map((g, i) => h(Fragment, { key: i },
    g.label ? h("span", { className: "ui-label" }, g.label) : null,
    g.links.map((l) => h("a", { key: l.href + l.label, href: l.href, "aria-current": l.current ? "page" : undefined }, l.label, l.count ? h(Count, l.count) : null)))));
}

// Shell draws the frame inside the app's root; the app puts class "ui-page" on <body> once.
export function Shell({ bar, rail, header, notice, footer, wide, children }) {
  return h(Fragment, null,
    h(Bar, bar),
    h("div", { className: "ui-body", "data-rail": rail ? "" : undefined },
      rail ? h(Rail, { groups: rail }) : null,
      h("main", { className: "ui-main ui-gutter", id: "main" }, h("div", { className: "ui-wrap", style: wide ? { maxWidth: "none" } : undefined },
        header && header.title ? h(PageHeader, header) : null, notice ? h(Notice, notice) : null, children))),
    footer ? h(Footer, footer) : null);
}
