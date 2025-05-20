import type { MetadataRoute } from 'next'

export default function robots(): MetadataRoute.Robots {
  // staging robots.txt
  if (process.env.NEXT_PUBLIC_VERCEL_ENV !== 'production') {
    return {
      rules: {
        userAgent: '*',
        disallow: '/',
      },
    }
  }

  // production robots.txt
  return {
    rules: {
      userAgent: '*',
      allow: '/',
    },
    host: 'https://nitric.io/docs',
    sitemap: 'https://nitric.io/docs/sitemap.xml',
  }
}
