import { notFound } from 'next/navigation'
import MDXContent from '@/components/MDXContent'
import { Metadata } from 'next/types'
import { BASE_URL } from '@/lib/constants'
import { allDocuments } from '@/content'
import { getNavInfo } from '@/lib/getNavInfo'
import { title } from 'radash'
import { OpenGraph } from 'next/dist/lib/metadata/types/opengraph-types'

export async function generateMetadata({
  params,
}: {
  params: { slug: string[] }
}): Promise<Metadata | undefined> {
  const slug = params.slug ? decodeURI(params.slug.join('/')) : ''
  const doc = allDocuments.find((p) => p.slug === slug)

  if (!doc) {
    return
  }

  const { navItem } = getNavInfo(doc)

  const title_meta = doc.title_seo || navItem?.title || doc.title

  let seoTitle =
    doc.slug.split('/').length > 1
      ? `${title(doc.slug.split('/')[0])}: ${title_meta}`
      : title_meta

  if (navItem && navItem.breadcrumbParentItem?.title) {
    seoTitle = `${navItem.breadcrumbParentItem.title}: ${title_meta}`
  }

  const url = doc.slug ? `${BASE_URL}/docs/${doc.slug}` : `${BASE_URL}/docs`

  const openGraph: OpenGraph = {
    siteName: 'Nitric Docs',
    locale: 'en_US',
    url,
    images: !doc.slug // is home page
      ? [
          {
            url: `${BASE_URL}/docs/images/og_home.png`,
            alt: 'Nitric Docs',
          },
        ]
      : [
          {
            url: `${BASE_URL}/docs/og?title=${encodeURIComponent(title_meta)}&description=${encodeURIComponent(doc.description || '')}`,
            alt: 'Nitric Docs',
          },
        ],
  }

  if (doc.image && doc.image_alt) {
    openGraph.images = [
      {
        url: doc.image,
        alt: doc.image_alt,
      },
    ]
  }

  return {
    title: seoTitle,
    description: doc.description,
    openGraph,
    twitter: {
      card: 'summary_large_image',
    },
    alternates: {
      canonical: doc.canonical_url ? doc.canonical_url : url,
    },
    ...(doc.noindex && {
      robots: {
        index: false,
        follow: false,
      },
    }),
  }
}

export const generateStaticParams = async () => {
  return allDocuments.map((p) => ({
    slug: p.slug.split('/').map((name) => decodeURI(name)),
  }))
}

export default function Page({ params }: { params: { slug: string[] } }) {
  const slug = params.slug ? decodeURI(params.slug.join('/')) : ''

  const doc = allDocuments.find((p) => p.slug === slug)

  if (!doc) {
    return notFound()
  }

  return <MDXContent mdx={doc.body.raw} />
}
