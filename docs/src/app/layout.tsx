import { type Metadata } from 'next'
import { Inter, Sora, JetBrains_Mono } from 'next/font/google'

import { Providers } from '@/app/providers'

import '@/styles/mermaid.css'
import '@/styles/tailwind.css'
import clsx from 'clsx'
import Fathom from './Fathom'
import { BaseLayout } from '@/components/layout/BaseLayout'

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap',
})

const sora = Sora({
  weight: ['400', '500', '600', '700'],
  variable: '--font-sora',
  subsets: ['latin'],
  display: 'swap',
})

const jetBrainsMono = JetBrains_Mono({
  weight: ['500', '600', '700'],
  variable: '--font-jetbrains-mono',
  display: 'swap',
  adjustFontFallback: false,
  subsets: ['latin'],
})

const isProd = process.env.NEXT_PUBLIC_VERCEL_ENV === 'production'

const defaultFallbackTitle =
  'Nitric Cloud-Native Framework | Nitric Documentation'

const defaultDescription =
  'Documentation for the Nitric cloud application framework.'

export const metadata: Metadata = {
  title: {
    template: '%s | Nitric Documentation',
    default: defaultFallbackTitle,
  },
  description: defaultDescription,
  robots: {
    index: isProd,
    follow: isProd,
  },
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html
      lang="en"
      className={clsx(
        'h-full',
        inter.variable,
        sora.variable,
        jetBrainsMono.variable,
      )}
      suppressHydrationWarning
    >
      {/* group/body is used by ShowIfLang */}
      <body className="group/body flex min-h-full antialiased">
        <Fathom />
        <Providers>
          <BaseLayout>{children}</BaseLayout>
        </Providers>
      </body>
    </html>
  )
}
