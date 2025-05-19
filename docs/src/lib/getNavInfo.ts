import { navigation } from '@/config'
import { NavEntry, NavItem } from '@/config/types'
import { DocumentTypes } from '@/content'

interface NavItemWithBreadcrumb extends NavItem {
  breadcrumbParentItem?: NavEntry
}

export interface NavInfo {
  prevItem: NavItemWithBreadcrumb | null
  navItem: NavItemWithBreadcrumb | null
  nextItem: NavItemWithBreadcrumb | null
}

const flatten = (
  entries: NavEntry[],
  rootParent?: NavEntry,
  breadcrumbParentItem?: NavEntry,
): NavItemWithBreadcrumb[] => {
  return entries.flatMap((entry) => {
    if ('href' in entry) {
      if (entry.breadcrumbRoot) {
        breadcrumbParentItem = {
          ...entry,
          title: rootParent?.title || entry.title, // take title of parent if available
        }
      }

      return breadcrumbParentItem
        ? {
            ...entry,
            breadcrumbParentItem,
          }
        : entry
    }
    return flatten(entry.items, entry, breadcrumbParentItem)
  })
}

// This function should return an object that contains the current nav item, the previous nav item, and the next nav item.
export function getNavInfo(doc: DocumentTypes): NavInfo {
  const slug = doc.slug ? `/${doc.slug}` : '/'

  const flattenedNav = flatten(navigation)

  const index = flattenedNav.findIndex((item) => item.href === slug)

  return {
    navItem: flattenedNav[index] || null,
    prevItem: flattenedNav[index - 1] || null,
    nextItem: flattenedNav[index + 1] || null,
  }
}
