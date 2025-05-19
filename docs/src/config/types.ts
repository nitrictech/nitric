export interface BaseNavItem {
  title: string
  icon?: React.ComponentType<{ className?: string }>
}

export interface NavItem extends BaseNavItem {
  href: string
  items?: NavItem[]
  breadcrumbRoot?: boolean
}

export interface NavGroup extends BaseNavItem {
  items: NavEntry[]
}

export type NavEntry = NavItem | NavGroup
