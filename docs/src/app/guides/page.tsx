import GuidePage from '@/components/guides/GuidePage'
import GuidesFeatured from '@/components/guides/GuidesFeatured'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Heading } from '@/components/ui/heading'
import { allGuides } from '@/content'
import { BASE_URL } from '@/lib/constants'
import Link from 'next/link'
import { Metadata } from 'next/types'

export const metadata: Metadata = {
  title: 'Guides',
  description:
    'Guides and tutorials for the Nitric cloud application framework.',
  openGraph: {
    siteName: 'Nitric Docs',
    locale: 'en_US',
    url: `${BASE_URL}/docs/guides`,
    images: [
      {
        url: `${BASE_URL}/docs/og?title=${encodeURIComponent('Guides')}&description=${encodeURIComponent('Guides and tutorials for the Nitric cloud application framework.')}`,
        alt: 'Nitric Docs',
      },
    ],
  },
  alternates: {
    canonical: `${BASE_URL}/docs/guides`,
  },
}

export default function GuidesPage() {
  const allTags = allGuides
    .reduce((acc: string[], guide) => {
      if (guide.tags) {
        guide.tags.forEach((tag) => {
          if (!acc.includes(tag)) {
            acc.push(tag)
          }
        })
      }
      return acc
    }, [])
    .sort()

  return (
    <>
      <div className="mx-auto flex h-full max-w-7xl flex-col gap-y-10 px-4 py-16">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link href={'/'}>Docs</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>Guides</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
        <Heading level={1}>Guides</Heading>
        <div className="xl:mr-6">
          <GuidesFeatured />
        </div>
      </div>
      <div className="-mx-2 border-t px-4 sm:-mx-6 lg:-mx-8">
        <div className="mx-auto max-w-7xl px-4">
          <GuidePage allTags={allTags} />
        </div>
      </div>
    </>
  )
}
