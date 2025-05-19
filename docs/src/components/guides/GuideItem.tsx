import type { Guide } from '@/content'
import React from 'react'
import { cn } from '@/lib/utils'
import Link from 'next/link'
import { LanguageIcon } from '../icons/LanguageIcon'
import { Language } from '@/lib/constants'

interface Props {
  guide: Guide
  featured?: boolean
}

export const GuideItem: React.FC<Props> = ({ guide, featured }) => {
  return (
    <Link
      href={`/${guide.slug}`}
      className={cn(
        'group relative flex overflow-hidden rounded-lg border p-3 transition-colors hover:bg-zinc-700/5 dark:border-zinc-700 dark:text-white dark:hover:bg-zinc-700/30',
        featured ? 'flex-col lg:flex-row' : 'flex-col',
      )}
    >
      <div className="flex flex-col gap-y-4">
        <div>
          <div className="flex items-start justify-between gap-x-1">
            <p
              className={cn(
                'font-display text-xl font-semibold',
                featured ? 'lg:text-2xl xl:text-3xl' : '',
              )}
            >
              {guide.title}
            </p>
            <time
              dateTime={guide.published_at}
              className="text-nowrap text-2xs text-muted-foreground"
            >
              {new Date(guide.published_at).toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
              })}
            </time>
          </div>

          <p
            className={cn(
              'text-md mt-3 text-base text-foreground dark:text-foreground-light',
              featured ? 'lg:text-md' : '',
            )}
          >
            {guide.description}
          </p>
        </div>
        <div className="flex">
          {guide.tags?.length ? (
            <div className="flex flex-wrap items-center gap-2">
              {guide.tags.map((tag) => (
                <span
                  key={tag}
                  className="items-center rounded-md border bg-zinc-600/10 px-2.5 py-0.5 text-xs font-semibold text-foreground transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                >
                  {tag}
                </span>
              ))}
            </div>
          ) : null}
          {guide.languages?.length ? (
            <div className="ml-auto flex items-center gap-x-2">
              {guide.languages.map((lang) => (
                <LanguageIcon
                  key={lang}
                  name={lang as Language}
                  className="size-6"
                />
              ))}
            </div>
          ) : null}
        </div>
      </div>
    </Link>
  )
}
