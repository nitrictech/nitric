import { BASE_URL } from '@/lib/constants'

// Function to construct the XML structure of the sitemap index.
export async function GET() {
  const sitemapIndexXML = buildSitemapIndex([`${BASE_URL}/docs/sitemap-0.xml`])

  // Return the sitemap index XML with the appropriate content type.
  return new Response(sitemapIndexXML, {
    headers: {
      'Content-Type': 'application/xml',
    },
  })
}

function buildSitemapIndex(sitemaps: string[]) {
  // XML declaration and opening tag for the sitemap index.
  let xml = '<?xml version="1.0" encoding="UTF-8"?>'
  xml += '<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">'

  // Iterate over each sitemap URL and add it to the sitemap index.
  for (const sitemapURL of sitemaps) {
    xml += '<sitemap>'
    xml += `<loc>${sitemapURL}</loc>` // Location tag specifying the URL of a sitemap file.
    xml += '</sitemap>'
  }

  // Closing tag for the sitemap index.
  xml += '</sitemapindex>'
  return xml
}
