// Types for the kit's React wrappers. The shapes match the Go types in types.go, in camelCase.
import type { ReactElement, ReactNode, ChangeEventHandler, MouseEventHandler, FormHTMLAttributes } from "react";

export type Tone = "good" | "warn" | "bad" | "info" | "accent" | "none" | "";

export interface Link { label: string; href: string; title?: string; current?: boolean; count?: CountProps }
export interface StatusProps { tone?: Tone; label: string }
export interface PillProps { tone?: Tone; variant?: "" | "outline" | "tag"; label?: string; title?: string; href?: string; pulse?: boolean; children?: ReactNode }
export interface CountProps { tone?: Tone; n: number; label?: string; showZero?: boolean; icon?: string; href?: string }
export interface StarProps { label?: string; text?: string; chip?: boolean; corner?: boolean }
export interface NumProps { text: string; href?: string; command?: string }
export interface ButtonProps {
  label?: string; variant?: "" | "primary" | "quiet" | "danger"; size?: "" | "sm"; type?: "button" | "submit" | "reset";
  href?: string; kbd?: string; disabled?: boolean; busy?: boolean; value?: string; name?: string;
  onClick?: MouseEventHandler; children?: ReactNode;
}
export interface Problem { title: string; cause?: string; command?: string }
export interface NoticeProps { tone?: "good" | "info" | "warn"; title: string; text?: string; link?: Link }
export interface EmptyProps { title?: string; text?: string; command?: string; action?: ButtonProps; small?: boolean }
export interface SkeletonProps { shape?: "card" | "rows"; lines?: number }
export interface Fact { label: string; tone?: Tone; num?: NumProps; text?: string; branch?: string; pill?: PillProps; error?: string }
export interface Product {
  name: string; href: string; kind?: string; variant?: "" | "planned" | "library"; focus?: boolean;
  marker?: PillProps; waiting?: CountProps; problem?: Problem; facts?: Fact[]; foot?: string; openHref?: string;
}
export interface Option { value: string; label: string; selected?: boolean }
export interface Select { name: string; label: string; options: Option[]; active?: boolean; onChange?: ChangeEventHandler<HTMLSelectElement> }
export interface Search { name?: string; label?: string; placeholder?: string; value?: string; onChange?: ChangeEventHandler<HTMLInputElement> }
export type View = "cards" | "list";
export interface ToolbarProps {
  search?: Search; filters?: Select[]; clear?: { label?: string; href?: string; onClick?: MouseEventHandler };
  count?: string; view?: View; onView?: (view: View) => void; primary?: ButtonProps; onSubmit?: (data: FormData) => void;
}
export interface ProductsProps {
  view?: View; onView?: (view: View) => void; toolbar?: ToolbarProps; columns?: string[]; items?: Product[];
  empty?: EmptyProps; problem?: Problem; loading?: boolean;
}
export interface PR {
  repo: string; number?: string; href?: string; repoHref?: string; author?: string; initials?: string; age?: string;
  draft?: boolean; title?: string; verb?: string; branch?: string; base?: string; files?: NumProps; added?: string; removed?: string;
  conflict?: string; review?: PillProps; ci?: PillProps; open?: string; compact?: boolean; problem?: Problem;
}
export interface PrListProps { items?: PR[]; empty?: EmptyProps; loading?: boolean }
export interface Mark { kind?: "point" | "past" | "today" | "span"; at: number; width?: number }
export interface Milestone { name: string; when: string; days?: string; focus?: boolean }
export interface MilestonesProps { label: string; total?: string; track?: Mark[]; items?: Milestone[]; empty?: EmptyProps }
export interface SectionProps {
  id?: string; title: string; count?: CountProps; why?: string; problem?: Problem; empty?: EmptyProps; loading?: boolean;
  more?: Link; failTitle?: string; children?: ReactNode;
}
export interface Column { label: string; num?: boolean; hidden?: boolean }
export interface Cell { text?: string; href?: string; sub?: string; num?: boolean; status?: StatusProps; pill?: PillProps; id?: string; star?: StarProps }
export interface Row { cells?: Cell[]; error?: string }
export interface TableProps { label?: string; columns?: Column[]; rows?: Row[]; empty?: EmptyProps; loading?: boolean; small?: boolean }
export interface TicketProps { title: string; href?: string; summary?: string; tags?: PillProps[]; turn?: string; age?: string; id?: string; focus?: boolean }
export interface FieldProps {
  id: string; name?: string; label: string; type?: string; value?: string; defaultValue?: string; placeholder?: string;
  onChange?: ChangeEventHandler<HTMLInputElement>; hint?: string; error?: string; secret?: boolean; saved?: string; savedAction?: ButtonProps;
}
export interface Tab { label: string; href?: string; panel?: string; count?: CountProps; selected?: boolean }
export interface TabsProps { label: string; items?: Tab[]; onSelect?: (tab: Tab) => void }
export interface DialogProps { id: string; title: string; text?: string; cancel?: string; confirm: ButtonProps; open?: boolean; onClose?: () => void; onConfirm?: () => void; children?: ReactNode }
export interface ConfirmFormProps extends Omit<FormHTMLAttributes<HTMLFormElement>, "id" | "onSubmit"> {
  id: string; question: string; label?: string; confirmName?: string; remember?: boolean; onSubmit?: (data: FormData) => void; children?: ReactNode;
}
export type ActionType = "create" | "edit" | "move" | "focus" | "archive" | "close" | "delete" | "settings" | "connect" | "run";
export type ConfirmStyle = "name" | "yes-no" | "never";
export type ConfirmSource = "this app" | "default";
export interface ConfirmSettingsProps {
  action?: string; types?: string[]; app?: Record<string, string>; system?: string[]; hidden?: { name: string; value: string }[]; save?: ButtonProps; onSubmit?: (data: FormData) => void;
}
export interface Location { label: string; path: string; note?: string }
export interface LicenseFile { kind: "org" | "product" | "trial"; for: string; until: string; daysLeft: number; ended: boolean; keyId: string; window?: { of: number } }
export interface RefusedFile { file: string; type: string; reason: string }
export interface FeatureValue { name: string; purpose: string }
export interface LicenseFeature { id: string; title: string; covers: string; needs?: FeatureValue[]; tier: "free" | "enterprise"; on: boolean; by?: string; startsOff?: true; missing?: FeatureValue[]; from: "org" | "product" | "both" | "trial" | "none"; until?: string }
export interface LicenseCounts { free: number; enterprise: number; on: number; locked: number }
export interface LicenseProps { state: "none" | "trial" | "trial-ended" | "licensed" | "ending" | "ended" | string; build: { release: boolean }; licenses?: LicenseFile[]; refused?: RefusedFile[]; features?: LicenseFeature[]; counts?: LicenseCounts; version?: string; trial?: string; actions?: ButtonProps[] }
export interface LocationsProps { home?: string; items?: Location[] }
export interface DrawerProps { id: string; title: string; open?: boolean; onClose?: () => void; children?: ReactNode }
export interface LightDetail { label: string; value: string }
export interface LightConfirm { text: string; label: string; action?: string; hidden?: { name: string; value: string }[]; onConfirm?: () => void }
export interface LightAction { label: string; variant?: "danger"; href?: string; opens?: string; onSelect?: () => void; confirm?: LightConfirm }
export interface LightProps { id: string; tone?: Tone; icon?: string; label: string; href?: string; title?: string; text?: string; details?: LightDetail[]; actions?: LightAction[]; open?: boolean; children?: ReactNode }
export interface BarActionProps { button: ButtonProps; action: string; target?: string; count?: CountProps }
export interface BarProps { app: string; sub?: string; href?: string; pages?: Link[]; apps?: Link[]; status?: StatusProps[]; lights?: LightProps[]; actions?: BarActionProps[]; noTheme?: boolean; children?: ReactNode }
export interface HeadProps { crumb?: Link[]; title: string; lede?: string; actions?: ButtonProps[] }
export interface Build { version?: string; branch?: string; commit?: string; dirty?: boolean; released?: boolean }
export interface FooterProps { version?: string; versionHref?: string; branch?: string; fresh?: StatusProps; links?: Link[] }
export interface RailGroup { label?: string; links: Link[] }
export interface ShellProps { build?: Build; bar: BarProps; rail?: RailGroup[]; header?: HeadProps; notice?: NoticeProps; footer?: FooterProps; wide?: boolean; children?: ReactNode }

export function Icon(props: { name: string }): ReactElement;
export function Status(props: StatusProps): ReactElement;
export function Pill(props: PillProps): ReactElement;
export function Count(props: CountProps): ReactElement | null;
export function Star(props: StarProps): ReactElement;
export function Num(props: NumProps): ReactElement;
export function Button(props: ButtonProps): ReactElement;
export function ErrorBox(props: Problem): ReactElement;
export function Notice(props: NoticeProps): ReactElement;
export function Empty(props: EmptyProps): ReactElement;
export function Skeleton(props: SkeletonProps): ReactElement;
export function ProductCard(props: Product): ReactElement;
export function ProductRow(props: Product): ReactElement;
export function Products(props: ProductsProps): ReactElement;
export function PrItem(props: PR): ReactElement;
export function PrList(props: PrListProps): ReactElement | null;
export function Milestones(props: MilestonesProps): ReactElement;
export function Section(props: SectionProps): ReactElement;
export function Table(props: TableProps): ReactElement;
export function Ticket(props: TicketProps): ReactElement;
export function Field(props: FieldProps): ReactElement;
export function ViewToggle(props: { view?: View; onView?: (view: View) => void }): ReactElement;
export function Toolbar(props: ToolbarProps): ReactElement;
export function Tabs(props: TabsProps): ReactElement;
export function Crumb(props: { items?: Link[] }): ReactElement;
export function Dialog(props: DialogProps): ReactElement;
export function ConfirmForm(props: ConfirmFormProps): ReactElement;
export const actionTypes: ActionType[];
export function resolveConfirm(app: string | undefined, def?: ConfirmStyle): { style: ConfirmStyle; source: ConfirmSource };
export function ConfirmSettings(props: ConfirmSettingsProps): ReactElement;
export function Locations(props: LocationsProps): ReactElement;
export function License(props: LicenseProps): ReactElement;
export function LicenseFeatures(props: { features?: LicenseFeature[]; counts?: LicenseCounts }): ReactElement;
export function Drawer(props: DrawerProps): ReactElement;
export function ThemeToggle(): ReactElement;
export function Light(props: LightProps): ReactElement;
export function Bar(props: BarProps): ReactElement;
export function BarAction(props: BarActionProps): ReactElement;
export function PageHeader(props: HeadProps): ReactElement;
export function PreviewBanner(props: Build): ReactElement | null;
export function Footer(props: FooterProps): ReactElement;
export function Rail(props: { groups?: RailGroup[] }): ReactElement;
export function Shell(props: ShellProps): ReactElement;
