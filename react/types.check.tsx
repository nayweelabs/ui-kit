// Compiles only if the types in index.d.ts fit the way an app uses the wrappers.
import { Shell, Section, PrList, Products, Pill, Count, Field, Dialog, ConfirmForm, ConfirmSettings, resolveConfirm, Light, type PR, type Product } from "./index";

const prs: PR[] = [{ repo: "checkman", number: "5", title: "chore: close the accepted review cards", review: { tone: "accent", label: "your review" }, ci: { tone: "good", label: "CI passing" } }];
const products: Product[] = [{ name: "hubman", href: "/p/hubman", focus: true, facts: [{ label: "Installed", tone: "good", num: { text: "v0.7.0" }, text: "current" }] }];

export const page = (
  <Shell build={{ branch: "feat/road-focus", commit: "758ac48", dirty: true, released: false }} bar={{ app: "Control room", pages: [{ label: "Focus", href: "/", current: true }] }} header={{ title: "Focus" }} footer={{ version: "v0.4.0" }}>
    <Section id="prs" title="Open pull requests" count={{ tone: "accent", n: 1 }}>
      <PrList items={prs} />
    </Section>
    <Products view="list" onView={() => {}} columns={["Installed"]} items={products} />
    <Pill tone="warn" label="3 behind" />
    <Count n={128} tone="accent" />
    <Field id="token" label="GitHub token" secret />
    <Dialog id="d" title="Archive 4 cards?" confirm={{ label: "Archive 4 cards", variant: "danger" }} />
    <ConfirmForm id="ask-save" question="Save the settings?" remember method="post" action="/settings" onSubmit={(data) => data.get("confirm")}><button type="submit">Save</button></ConfirmForm>
    <ConfirmForm id="ask-delete" question="Delete card 7?" confirmName="7" remember={resolveConfirm({ delete: "name" }, {}, "delete").style !== "never"}><button type="submit">Delete</button></ConfirmForm>
    <ConfirmSettings action="/settings" scope="global" types={["move", "delete"]} app={{ move: "never" }} global={{ delete: "yes-no" }} />
    <Light id="gh" tone="good" icon="github" label="GitHub" details={[{ label: "Account", value: "bot" }]} actions={[{ label: "Enter a new token…", opens: "d" }, { label: "Disconnect", variant: "danger", confirm: { text: "Removes it.", label: "Disconnect now", onConfirm: () => {} } }]} />
  </Shell>
);
