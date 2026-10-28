// Compiles only if the types in index.d.ts fit the way an app uses the wrappers.
import { Shell, Section, PrList, Products, Pill, Count, Field, Dialog, type PR, type Product } from "./index";

const prs: PR[] = [{ repo: "checkman", number: "5", title: "chore: close the accepted review cards", review: { tone: "accent", label: "your review" }, ci: { tone: "good", label: "CI passing" } }];
const products: Product[] = [{ name: "hubman", href: "/p/hubman", focus: true, facts: [{ label: "Installed", tone: "good", num: { text: "v0.7.0" }, text: "current" }] }];

export const page = (
  <Shell bar={{ app: "Control room", pages: [{ label: "Focus", href: "/", current: true }] }} header={{ title: "Focus" }} footer={{ version: "v0.4.0" }}>
    <Section id="prs" title="Open pull requests" count={{ tone: "accent", n: 1 }}>
      <PrList items={prs} />
    </Section>
    <Products view="list" onView={() => {}} columns={["Installed"]} items={products} />
    <Pill tone="warn" label="3 behind" />
    <Count n={128} tone="accent" />
    <Field id="token" label="GitHub token" secret />
    <Dialog id="d" title="Archive 4 cards?" confirm={{ label: "Archive 4 cards", variant: "danger" }} />
  </Shell>
);
