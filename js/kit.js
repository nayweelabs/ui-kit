/* ui-kit/js/kit.js: theme switch, cards/list toggle, fold, tabs, lazy sections, carried links, status lights, and the small helpers of the field, id and dialog parts.
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
      carry();
    }
  }

  // Carried links: a link with data-ui-carry takes the page's query when the address changes without a
  // reload, and keeps its own values, the names the attribute lists. data-ui-back holds the page's address:
  // in a hidden input's value, or in the value it names of a lazy box's data-ui-src.
  function carry(within) {
    var scope = within || document;
    var here = location.pathname + location.search;
    scope.querySelectorAll("[data-ui-carry]").forEach(function (a) {
      var url = new URL(a.getAttribute("href"), location.href);
      if (url.origin !== location.origin) return;
      var q = new URLSearchParams(location.search);
      (a.getAttribute("data-ui-carry") || "").split(/\s+/).forEach(function (k) {
        if (!k) return;
        q.delete(k);
        url.searchParams.getAll(k).forEach(function (v) { q.append(k, v); });
      });
      var query = q.toString();
      a.setAttribute("href", url.pathname + (query ? "?" + query : "") + url.hash);
    });
    scope.querySelectorAll("[data-ui-back]").forEach(function (el) {
      var name = el.getAttribute("data-ui-back");
      if (!name) {
        el.value = here;
        el.setAttribute("value", here);
        return;
      }
      var url = new URL(el.getAttribute("data-ui-src"), location.href);
      url.searchParams.set(name, here);
      el.setAttribute("data-ui-src", url.pathname + url.search);
    });
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
  function initFolds(within) {
    (within || document).querySelectorAll("details.ui-fold[id]").forEach(function (d) {
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

  // Status light: its link opens its popover instead of the page; Escape, a click outside or focus
  // leaving closes it. An action with a step shows that step in the popover, Cancel first.
  function setPop(btn, open, focusBack) {
    var pop = document.getElementById(btn.getAttribute("aria-controls"));
    if (!pop) return false;
    pop.hidden = !open;
    btn.setAttribute("aria-expanded", String(open));
    if (open) {
      pop.focus();
    } else {
      pop.querySelectorAll("details.ui-menu[open]").forEach(function (d) { d.open = false; });
      pop.querySelectorAll(".ui-pop-step").forEach(function (f) { f.hidden = true; });
      pop.querySelectorAll("[data-ui-step]").forEach(function (b) { b.setAttribute("aria-expanded", "false"); });
      if (focusBack) btn.focus();
    }
    return true;
  }

  function closePops(except, focusBack) {
    document.querySelectorAll("[data-ui-light][aria-expanded='true']").forEach(function (btn) {
      if (!except || !btn.closest(".ui-light").contains(except)) setPop(btn, false, focusBack);
    });
  }

  function lightClick(t, e) {
    var btn = t.closest("[data-ui-light]");
    if (btn) {
      var open = btn.getAttribute("aria-expanded") !== "true";
      closePops(btn);
      if (setPop(btn, open, true)) e.preventDefault();
      return true;
    }
    var step = t.closest("[data-ui-step]");
    if (step) {
      var form = document.getElementById(step.getAttribute("aria-controls"));
      if (!form) return true;
      var menu = step.closest("details");
      if (menu) menu.open = false;
      form.hidden = false;
      step.setAttribute("aria-expanded", "true");
      var cancel = form.querySelector("[data-ui-step-cancel]");
      if (cancel) cancel.focus();
      return true;
    }
    var back = t.closest("[data-ui-step-cancel]");
    if (back) {
      var f = back.closest(".ui-pop-step");
      f.hidden = true;
      var opener = document.querySelector('[data-ui-step][aria-controls="' + f.id + '"]');
      if (opener) opener.setAttribute("aria-expanded", "false");
      var sum = f.closest(".ui-pop").querySelector("details.ui-menu > summary");
      if (sum) sum.focus();
      return true;
    }
    return false;
  }

  function onKey(e) {
    if (e.key !== "Escape" || e.target.closest("dialog")) return;
    var open = document.querySelector("[data-ui-light][aria-expanded='true']");
    if (open) {
      e.preventDefault();
      setPop(open, false, true);
    }
  }

  function onFocus(e) {
    if (!e.target.closest || e.target.closest("dialog")) return;
    closePops(e.target, false);
  }

  function onClick(e) {
    var t = e.target;
    if (!t.closest("dialog")) closePops(t, false);
    if (lightClick(t, e)) return;
    if (t.closest("[data-ui-theme-toggle]")) return toggleTheme();

    var again = t.closest(".ui-lazy [data-ui-again]");
    if (again && window.fetch) {
      e.preventDefault();
      if (again.getAttribute("aria-busy") === "true") return;
      return load(again.closest(".ui-lazy"), null, true);
    }

    // A section link changes only its lazy section's view: the address becomes its href, and the section reloads.
    var here = t.closest(".ui-lazy[data-ui-src] a[data-ui-section-link]");
    if (here && window.fetch && !e.metaKey && !e.ctrlKey && !e.shiftKey) {
      e.preventDefault();
      return reloadAt(here.closest(".ui-lazy"), new URL(here.getAttribute("href"), location.href));
    }

    var retry = t.closest(".ui-lazy [data-ui-retry]");
    if (retry && window.fetch) {
      e.preventDefault();
      return load(retry.closest(".ui-lazy"));
    }

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
      var menu = opener.closest("details.ui-menu");
      if (menu) menu.open = false;
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

  // Lazy sections: each .ui-lazy fetches its own HTML from data-ui-src, all in parallel, and shows it in
  // its place; a failed fetch shows the section's error box (from its <template data-ui-failed>) with Retry.
  // The placeholder's tools (when it was read, Read again) stay in the loaded section's head. Read again
  // reloads with fresh=1: the button spins, the section dims under a bar until the new HTML lands, then
  // "Updated" shows and is announced. The UI-Read-At header says when the section's sources were read.
  function failed(box, cause) {
    var html = box.uiFailed || "";
    box.innerHTML = html;
    var p = box.querySelector("[data-ui-cause]");
    if (p) p.textContent = cause + ". Everything else on this page is current.";
  }

  function keepTools(box) {
    var head = box.querySelector(".ui-section-head");
    if (!box.uiTools || !head) return;
    var drawn = head.querySelector("[data-ui-tools]");
    if (drawn === box.uiTools) return;
    if (drawn) drawn.remove();
    head.appendChild(box.uiTools);
  }

  function readAt(box, at) {
    var t = box.uiTools && box.uiTools.querySelector("[data-ui-read]");
    var when = at && new Date(at);
    if (!t || !when || isNaN(when)) return;
    t.setAttribute("datetime", at);
    t.setAttribute("title", "Its sources were read at " + when.toLocaleString());
    t.textContent = "read " + when.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    t.hidden = false;
  }

  function busy(box, on) {
    var btn = box.uiTools && box.uiTools.querySelector("[data-ui-again]");
    if (btn) {
      if (on) btn.setAttribute("aria-busy", "true");
      else btn.removeAttribute("aria-busy");
      if (on) btn.setAttribute("aria-disabled", "true");
      else btn.removeAttribute("aria-disabled");
    }
    var sec = box.querySelector(".ui-section");
    if (sec && on) sec.setAttribute("data-ui-reloading", "");
  }

  function said(box, text) {
    var s = box.uiTools && box.uiTools.querySelector(".ui-again-done");
    if (!s) return;
    s.textContent = text;
    clearTimeout(box.uiSaid);
    if (text) box.uiSaid = setTimeout(function () { s.textContent = ""; }, 4000);
  }

  function load(box, src, fresh) {
    if (!box.uiFailed) {
      var tpl = box.querySelector("template[data-ui-failed]");
      box.uiFailed = tpl ? tpl.innerHTML : "";
      box.uiTools = box.querySelector("[data-ui-tools]");
    }
    if (src) box.setAttribute("data-ui-src", src);
    src = box.getAttribute("data-ui-src");
    if (fresh) {
      var url = new URL(src, location.href);
      url.searchParams.set("fresh", "1");
      src = url.pathname + url.search;
      said(box, "");
      busy(box, true);
    }
    var sec = box.querySelector(".ui-section");
    if (sec) sec.setAttribute("aria-busy", "true");
    var seq = (box.uiSeq || 0) + 1;
    box.uiSeq = seq;
    return fetch(src, { credentials: "same-origin", headers: { "Accept": "text/html" } }).then(function (res) {
      return res.text().then(function (text) {
        if (box.uiSeq !== seq) return;
        busy(box, false);
        if (!res.ok) {
          failed(box, "The server answered " + res.status + (res.statusText ? " " + res.statusText : ""));
          return keepTools(box);
        }
        box.innerHTML = text;
        keepTools(box);
        readAt(box, res.headers && res.headers.get("UI-Read-At"));
        if (fresh) said(box, "Updated");
        initFolds(box);
        carry(box);
      });
    }, function (err) {
      if (box.uiSeq !== seq) return;
      busy(box, false);
      failed(box, "The page could not reach the server (" + err.message + ")");
      keepTools(box);
    });
  }

  function initLazy() {
    if (!window.fetch) return;
    document.querySelectorAll(".ui-lazy[data-ui-src]").forEach(function (box) { load(box); });
  }

  // A GET form inside a lazy section reloads only that section; the address keeps the form's values.
  function onSubmit(e) {
    var form = e.target;
    var box = form.closest && form.closest(".ui-lazy[data-ui-src]");
    if (!box || !window.fetch || !window.URLSearchParams || (form.method || "get").toLowerCase() !== "get") return;
    e.preventDefault();
    var fields = new URLSearchParams(new FormData(form));
    var page = new URL(location.href);
    fields.forEach(function (v, k) { page.searchParams.set(k, v); });
    reloadAt(box, page);
  }

  // reloadAt makes page the address and reloads only box, from its src with the page's query.
  function reloadAt(box, page) {
    if (window.history && history.replaceState) {
      history.replaceState(null, "", page);
      carry();
    }
    var src = new URL(box.getAttribute("data-ui-src"), location.href);
    src.search = page.search;
    load(box, src.pathname + src.search);
  }

  // A filter select in a toolbar applies on change, the same as pressing Apply.
  function onChange(e) {
    var sel = e.target.closest("select[data-ui-submit]");
    if (sel && sel.form) sel.form.requestSubmit ? sel.form.requestSubmit() : sel.form.submit();
  }

  function init() {
    initViews();
    initFolds();
    initLazy();
    document.addEventListener("click", onClick);
    document.addEventListener("keydown", tabKey);
    document.addEventListener("keydown", onKey);
    document.addEventListener("focusin", onFocus);
    document.addEventListener("change", onChange);
    document.addEventListener("submit", onSubmit);
  }

  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init);
  else init();
})();
