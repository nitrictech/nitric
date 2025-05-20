import Link from 'next/link'
import clsx from 'clsx'
import { Table } from '@/components/ui/table'
import MermaidZoom from './MermaidZoom'

export {
  TableHead as th,
  TableHeader as thead,
  TableRow as tr,
  TableBody as tbody,
  TableCell as td,
  TableFooter as tfoot,
} from './ui/table'

export const table = (props: React.ComponentPropsWithoutRef<typeof Table>) => (
  <Table {...props} className="mt-4 text-base" />
)

export function a({
  href,
  children,
  ...props
}: React.ComponentPropsWithoutRef<typeof Link>) {
  const isExternal = href.toString().startsWith('http')

  return (
    <Link
      href={href}
      target={isExternal ? '_blank' : undefined}
      rel={isExternal ? 'noopener noreferrer' : undefined}
      {...props}
    >
      {children}
    </Link>
  )
}

export { Code } from './code/Code'

export { CodeSwitcher } from './code/CodeSwitcher'

function InfoIcon(props: React.ComponentPropsWithoutRef<'svg'>) {
  return (
    <svg viewBox="0 0 16 16" aria-hidden="true" {...props}>
      <circle cx="8" cy="8" r="8" strokeWidth="0" />
      <path
        fill="none"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.5"
        d="M6.75 7.75h1.5v3.5"
      />
      <circle cx="8" cy="4" r=".5" fill="none" />
    </svg>
  )
}

export function Note({ children }: { children: React.ReactNode }) {
  return (
    <div className="my-6 flex gap-2.5 rounded-2xl border border-primary-500/20 bg-primary-50/50 p-4 leading-6 text-primary-900 dark:border-primary-500/30 dark:bg-primary-500/5 dark:text-primary-200">
      <InfoIcon className="mt-1 h-4 w-4 flex-none fill-primary-500 stroke-white dark:fill-primary-200/20 dark:stroke-primary-200" />
      <div className="[&>:first-child]:mt-0 [&>:last-child]:mb-0">
        {children}
      </div>
    </div>
  )
}

export function Row({ children }: { children: React.ReactNode }) {
  return (
    <div className="grid grid-cols-1 items-start gap-x-16 gap-y-10 xl:max-w-none xl:grid-cols-2">
      {children}
    </div>
  )
}

export function Col({
  children,
  sticky = false,
}: {
  children: React.ReactNode
  sticky?: boolean
}) {
  return (
    <div
      className={clsx(
        '[&>:first-child]:mt-0 [&>:last-child]:mb-0',
        sticky && 'xl:sticky xl:top-24',
      )}
    >
      {children}
    </div>
  )
}

export { Properties, Property } from './Properties'

export { OSTabs } from '@/components/OSTabs'

export { HomeHeader } from '@/components/HomeHeader'

export { ShowIfLang } from '@/components/ShowIfLang'

export { LanguageSwitch } from '@/components/LanguageSwitch'

export { ImportCode } from '@/components/code/ImportCode'

export { Tabs, TabItem } from '@/components/tabs/Tabs'

export { CodeTabs } from '@/components/code/CodeTabs'

export { Mermaid } from 'mdx-mermaid/Mermaid'

export const svg = (props: React.ComponentPropsWithoutRef<'svg'>) => {
  const { id } = props

  if (id?.startsWith('mermaid-svg')) {
    return <MermaidZoom {...props} />
  }

  return <svg {...props} />
}

// see if we need to remove these
