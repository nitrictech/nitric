import React from 'react'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from './ui/breadcrumb'
import { Doc, Guide } from '@/content'
import { getNavInfo, NavInfo } from '@/lib/getNavInfo'
import Link from 'next/link'

interface Props {
  doc: Doc | Guide
  className?: string
}

const Breadcrumbs: React.FC<Props> = ({ doc, className }) => {
  // Implement your component logic here

  const docs = doc.slug.split('/')

  const isGuide = docs[0] === 'guides'

  let navInfo: NavInfo = getNavInfo(doc)

  // If the doc is a guide, generate the navInfo for the guides
  if (isGuide && !navInfo?.navItem) {
    navInfo = {
      prevItem: null,
      navItem: {
        title: doc.title,
        href: `/${doc.slug}`,
        breadcrumbParentItem: {
          title: 'Guides',
          href: '/guides',
        },
      },
      nextItem: null,
    }
  }

  if (
    docs.length === 1 ||
    !navInfo ||
    !navInfo.navItem?.breadcrumbParentItem ||
    navInfo.navItem.breadcrumbRoot
  ) {
    return null
  }

  const { breadcrumbParentItem } = navInfo.navItem

  return (
    <Breadcrumb className={className}>
      <BreadcrumbList>
        {'href' in breadcrumbParentItem ? (
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href={breadcrumbParentItem.href}>
                {breadcrumbParentItem.title}
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
        ) : null}
        <BreadcrumbSeparator />
        <BreadcrumbItem>
          <BreadcrumbPage>{navInfo.navItem.title}</BreadcrumbPage>
        </BreadcrumbItem>
      </BreadcrumbList>
    </Breadcrumb>
  )
}

export default Breadcrumbs
