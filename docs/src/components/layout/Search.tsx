'use client'

import { useCallback, useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import Link from 'next/link'
import { DocSearchModal, useDocSearchKeyboardEvents } from '@docsearch/react'
import { useRouter } from 'next/navigation'

const docSearchConfig = {
  appId: process.env.NEXT_PUBLIC_DOCSEARCH_APP_ID || '',
  apiKey: process.env.NEXT_PUBLIC_DOCSEARCH_API_KEY || '',
  indexName: process.env.NEXT_PUBLIC_DOCSEARCH_INDEX_NAME || '',
}

const cleanUrl = (href: string) => {
  // We transform the absolute URL into a relative URL to
  // work better on localhost, preview URLS.
  const url = new URL(href, window.location.origin)
  if (url.hash === '#overview') url.hash = ''

  return url.href.replace(url.origin, '').replace(/^\/docs/, '')
}

function Hit({ hit, children }: { hit: any; children: React.ReactNode }) {
  return <Link href={cleanUrl(hit.url)}>{children}</Link>
}

function SearchIcon(props: React.ComponentPropsWithoutRef<'svg'>) {
  return (
    <svg aria-hidden="true" viewBox="0 0 20 20" {...props}>
      <path d="M16.293 17.707a1 1 0 0 0 1.414-1.414l-1.414 1.414ZM9 14a5 5 0 0 1-5-5H2a7 7 0 0 0 7 7v-2ZM4 9a5 5 0 0 1 5-5V2a7 7 0 0 0-7 7h2Zm5-5a5 5 0 0 1 5 5h2a7 7 0 0 0-7-7v2Zm8.707 12.293-3.757-3.757-1.414 1.414 3.757 3.757 1.414-1.414ZM14 9a4.98 4.98 0 0 1-1.464 3.536l1.414 1.414A6.98 6.98 0 0 0 16 9h-2Zm-1.464 3.536A4.98 4.98 0 0 1 9 14v2a6.98 6.98 0 0 0 4.95-2.05l-1.414-1.414Z" />
    </svg>
  )
}

export function Search() {
  const router = useRouter()
  let [isOpen, setIsOpen] = useState(false)
  let [modifierKey, setModifierKey] = useState<string>()

  const onOpen = useCallback(() => {
    setIsOpen(true)
  }, [setIsOpen])

  const onClose = useCallback(() => {
    setIsOpen(false)
  }, [setIsOpen])

  useDocSearchKeyboardEvents({ isOpen, onOpen, onClose })

  useEffect(() => {
    setModifierKey(
      /(Mac|iPhone|iPod|iPad)/i.test(navigator.platform) ? '⌘' : 'Ctrl ',
    )
  }, [])

  return (
    <>
      <button
        type="button"
        className="lg:w-92 group hidden h-6 w-6 items-center justify-center sm:justify-start md:h-auto md:w-52 md:flex-none md:rounded-lg md:py-1 md:pl-4 md:pr-3.5 md:text-sm md:ring-1 md:ring-zinc-200 md:hover:ring-zinc-300 dark:md:bg-zinc-800/70 dark:md:ring-inset dark:md:ring-white/5 dark:md:hover:bg-zinc-700/40 dark:md:hover:ring-zinc-500 lg:flex"
        onClick={onOpen}
      >
        <SearchIcon className="h-5 w-5 flex-none fill-zinc-400 group-hover:fill-zinc-500 dark:fill-zinc-500 md:group-hover:fill-zinc-400" />
        <span className="sr-only transition-colors dark:group-hover:text-white md:not-sr-only md:ml-2 md:text-zinc-500 md:dark:text-zinc-400">
          Search
        </span>
        {modifierKey && (
          <kbd className="ml-auto hidden font-medium text-zinc-400 dark:text-zinc-500 md:block">
            <kbd className="font-sans">{modifierKey}</kbd>
            <kbd className="font-sans">K</kbd>
          </kbd>
        )}
      </button>
      <button
        type="button"
        className="group flex h-5 w-5 items-center justify-center rounded-lg md:py-2.5 lg:hidden"
        onClick={onOpen}
      >
        <span className="sr-only">Search Docs</span>
        <SearchIcon className="h-5 w-5 flex-none fill-zinc-400 group-hover:fill-zinc-500 dark:fill-zinc-500 md:group-hover:fill-zinc-400" />
      </button>
      {isOpen &&
        createPortal(
          <DocSearchModal
            {...docSearchConfig}
            initialScrollY={window.scrollY}
            onClose={onClose}
            hitComponent={Hit}
            getMissingResultsUrl={({ query }: { query: string }) =>
              `https://github.com/nitrictech/nitric/issues/new?title=Missing+results+for+docs+query+%22${encodeURIComponent(
                query,
              )}%22`
            }
            navigator={{
              navigate({ itemUrl }) {
                router.push(itemUrl)
              },
            }}
            transformItems={(items) => {
              return items.map((item) => ({
                ...item,
                url: cleanUrl(item.url),
              }))
            }}
          />,
          document.body,
        )}
    </>
  )
}
