// Fathom.tsx
'use client'

import { load, trackPageview } from 'fathom-client'
import { useEffect, Suspense } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'

const SITE_ID = process.env.NEXT_PUBLIC_FATHOM_SITE_ID

const enabled = !!SITE_ID

function TrackPageView() {
  const pathname = usePathname()
  const searchParams = useSearchParams()

  // Load the Fathom script on mount
  useEffect(() => {
    if (!enabled) {
      console.error('Fathom is not enabled')

      return
    }

    load(SITE_ID, {
      auto: false,
      // disable for preview branches in Vercel
      excludedDomains: ['vercel.app'],
    })
  }, [])

  // Record a pageview when route changes
  useEffect(() => {
    if (!pathname || !enabled) return

    trackPageview({
      url: pathname + searchParams?.toString(),
      referrer: document.referrer,
    })
  }, [pathname, searchParams])

  return null
}

export default function Fathom() {
  return (
    <Suspense fallback={null}>
      <TrackPageView />
    </Suspense>
  )
}
