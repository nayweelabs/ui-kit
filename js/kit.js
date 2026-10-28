/* ui-kit/js/kit.js: theme switch, cards/list toggle, fold, tabs, and the small helpers of the field, id and dialog parts.
   No app logic and no data. Every part works without it; this only saves a reload. Load it in <head> without defer. */
(function () {
  "use strict";
  var root = document.documentElement;
  var KEY = "ui-kit.";

  function store(k, v) {
    try {
      if (v === undefined) return localStorage.getItem(KEY + k);
      localStorage.setItem(KEY + k, v);
    } catch (e) { return null; }
    return null;
  }

  // Theme: follow the system until the viewer picks one; set before the first paint.
  var saved = store("theme");
  if (saved === "light" || saved === "dark") root.setAttribute("data-theme", saved);

  function toggleTheme() {
    var cur = root.getAttribute("data-theme");
    if (!cur) cur = window.matchMedia && matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
    var next = cur === "dark" ? "light" : "dark";
    root.setAttribute("data-theme", next);
    store("theme", next);
  }

  // View toggle: the choice lives in the address (?view=list) and in the browser.
  function setView(name, view, push) {
    document.querySelectorAll('[data-ui-view-panel="' + name + '"]').forEach(function (p) {
      p.hidden = p.getAttribute("data-view") !== view;
    });
    document.querySelectorAll('[data-ui-views="' + name + '"] [data-view]').forEach(function (b) {
      b.setAttribute("aria-pressed", String(b.getAttribute("data-view") === view));
    });
    store("view." + name, view);
    if (push && window.history && history.replaceState) {
      var url = new URL(location.href);
      url.searchParams.set("view", view);
      history.replaceState(null, "", url);
    }
  }

  function initViews() {
    var fromURL = new URLSearchParams(location.search).get("view");
    document.querySelectorAll("[data-ui-views]").forEach(function (seg) {
      var name = seg.getAttribute("data-ui-views");
      var view = fromURL || store("view." + name);
      if (view === "cards" || view === "list") setView(name, view, false);
    });
  }

  // Fold: a <details class="ui-fold" id="…"> remembers whether it was open.
  function initFolds() {
    document.querySelectorAll("details.ui-fold[id]").forEach(function (d) {
      var s = store("fold." + d.id);
      if (s === "open") d.open = true;
      if (s === "closed") d.open = false;
      d.addEventListener("toggle", function () { store("fold." + d.id, d.open ? "open" : "closed"); });
    });
  }

  // Tabs: arrow keys move between tabs; a tab without an address switches its panel.
  function selectTab(tab) {
    var list = tab.closest('[role="tablist"]');
    list.querySelectorAll('[role="tab"]').forEach(function (t) {
      var on = t === tab;
      t.setAttribute("aria-selected", String(on));
      t.tabIndex = on ? 0 : -1;
      var panel = t.getAttribute("aria-controls") && document.getElementById(t.getAttribute("aria-controls"));
      if (panel) panel.hidden = !on;
    });
  }

  function tabKey(e) {
    var tab = e.target.closest('[role="tab"]');
    if (!tab) return;
    var tabs = Array.prototype.slice.call(tab.closest('[role="tablist"]').querySelectorAll('[role="tab"]'));
    var i = tabs.indexOf(tab), n = -1;
    if (e.key === "ArrowRight") n = (i + 1) % tabs.length;
    else if (e.key === "ArrowLeft") n = (i - 1 + tabs.length) % tabs.length;
    else if (e.key === "Home") n = 0;
    else if (e.key === "End") n = tabs.length - 1;
    if (n < 0) return;
    e.preventDefault();
    tabs[n].focus();
    if (!tabs[n].getAttribute("href")) selectTab(tabs[n]);
  }

  function onClick(e) {
    var t = e.target;
    if (t.closest("[data-ui-theme-toggle]")) return toggleTheme();

    var viewBtn = t.closest("[data-ui-views] [data-view]");
    if (viewBtn) {
      e.preventDefault();
      return setView(viewBtn.closest("[data-ui-views]").getAttribute("data-ui-views"), viewBtn.getAttribute("data-view"), true);
    }

    var tab = t.closest('button[role="tab"][aria-controls]');
    if (tab) return selectTab(tab);

    var reveal = t.closest("[data-ui-reveal]");
    if (reveal) {
      var input = document.getElementById(reveal.getAttribute("data-ui-reveal"));
      if (!input) return;
      var show = input.type === "password";
      input.type = show ? "text" : "password";
      reveal.textContent = show ? "Hide" : "Show";
      reveal.setAttribute("aria-pressed", String(show));
      return;
    }

    var opener = t.closest("[data-ui-open]");
    if (opener) {
      var dlg = document.getElementById(opener.getAttribute("data-ui-open"));
      if (dlg && dlg.showModal && !dlg.open) dlg.showModal();
      return;
    }

    var copy = t.closest("[data-ui-copy]");
    if (copy && navigator.clipboard) {
      navigator.clipboard.writeText(copy.getAttribute("data-ui-copy")).then(function () {
        copy.setAttribute("data-copied", "");
        setTimeout(function () { copy.removeAttribute("data-copied"); }, 1200);
      });
    }
  }

  // A filter select in a toolbar applies on change, the same as pressing Apply.
  function onChange(e) {
    var sel = e.target.closest("select[data-ui-submit]");
    if (sel && sel.form) sel.form.requestSubmit ? sel.form.requestSubmit() : sel.form.submit();
  }

  function init() {
    initViews();
    initFolds();
    document.addEventListener("click", onClick);
    document.addEventListener("keydown", tabKey);
    document.addEventListener("change", onChange);
  }

  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init);
  else init();
})();
