'use client'

import { useState } from 'react'
import { Button } from './ui/button'
import { MinusCircleIcon, PlusCircleIcon } from '@heroicons/react/24/outline'
import { Tag } from './Tag'
import { cn } from '@/lib/utils'

export function Properties({
  children,
  nested,
}: {
  children: React.ReactNode
  nested?: boolean
}) {
  const [open, setOpen] = useState(false)

  const content = (
    <ul
      role="list"
      className="m-0 max-w-[calc(theme(maxWidth.lg)-theme(spacing.8))] list-none divide-y divide-zinc-900/5 p-0 dark:divide-white/5"
    >
      {children}
    </ul>
  )

  return nested ? (
    <li>
      <Button
        onClick={() => setOpen(!open)}
        variant="outline"
        className="mt-4 flex items-center gap-2"
      >
        {open ? (
          <MinusCircleIcon className="h-5 w-5" />
        ) : (
          <PlusCircleIcon className="h-5 w-5" />
        )}
        {open ? 'Hide' : 'Show'} accepted values
      </Button>
      <div className={cn('mt-4 rounded-lg border p-4', !open && 'sr-only')}>
        {content}
      </div>
    </li>
  ) : (
    <div className={'my-6'}>{content}</div>
  )
}

interface PropertyProps {
  name: string
  type: string
  children: React.ReactNode
  required: boolean
}

export function Property({ name, type, children, required }: PropertyProps) {
  return (
    <li className="m-0 px-0 py-4 first:pt-0 last:pb-0">
      <dl className="m-0 flex flex-wrap items-center gap-x-3 gap-y-2">
        <dt className="sr-only">Name</dt>
        <dd>
          <code>{name}</code>
        </dd>
        <dt className="sr-only">{required ? 'Required' : 'Optional'}</dt>
        <dd>
          {required ? (
            <Tag color="amber" className="rounded-md py-0.5 uppercase">
              Required
            </Tag>
          ) : (
            <span className="mt-0.5 font-mono text-xs text-zinc-400 dark:text-zinc-500">
              Optional
            </span>
          )}
        </dd>
        <dt className="sr-only">Type</dt>
        <dd className="mt-0.5 flex font-mono text-xs text-zinc-400 dark:text-zinc-500">
          {type}
        </dd>
        <dt className="sr-only">Description</dt>
        <dd className="w-full flex-none [&>:first-child]:mt-0 [&>:last-child]:mb-0">
          {children}
        </dd>
      </dl>
    </li>
  )
}
