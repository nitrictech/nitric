'use client'

import clsx from 'clsx'
import { useEffect, useState } from 'react'
import { ClipboardIcon } from '../icons/ClipboardIcon'
import { cn } from '@/lib/utils'
import { CodeAnnotation } from 'codehike/code'

// remove - diff annotations from the code
const cleanCode = (code: string, annotations: CodeAnnotation[]) => {
  if (!annotations.length) return code

  const lines = code.split('\n')

  // Sort annotations by fromLineNumber in descending order to avoid conflicts
  const sortedAnnotations = annotations
    .filter((a) => 'fromLineNumber' in a)
    .sort((a, b) => b.fromLineNumber - a.fromLineNumber)

  for (const annotation of sortedAnnotations) {
    if (annotation.query === '-' && annotation.name === 'diff') {
      const { fromLineNumber, toLineNumber } = annotation
      const start = fromLineNumber - 1
      const end = toLineNumber

      // Remove the lines in the range [start, end)
      lines.splice(start, end - start)
    }
  }

  return lines.join('\n')
}

export function CopyButton({
  code,
  annotations,
  className,
}: {
  code: string
  annotations: CodeAnnotation[]
  className?: string
}) {
  const [copyCount, setCopyCount] = useState(0)
  const copied = copyCount > 0

  useEffect(() => {
    if (copyCount > 0) {
      const timeout = setTimeout(() => setCopyCount(0), 1000)
      return () => {
        clearTimeout(timeout)
      }
    }
  }, [copyCount])

  return (
    <button
      type="button"
      className={cn(
        'group/button absolute right-3.5 top-2.5 z-10 h-8 rounded-md px-1.5 py-1 text-2xs font-medium backdrop-blur transition',
        copied
          ? 'bg-primary-400/10 ring-primary-400/20'
          : 'bg-white/5 ring-1 ring-inset ring-zinc-300/80 hover:bg-white/7.5 dark:bg-white/2.5 dark:ring-zinc-300/10 dark:hover:bg-white/5',
        'opacity-0 focus:opacity-100 group-hover:opacity-100 group-focus:opacity-100',
        className,
      )}
      onClick={() => {
        window.navigator.clipboard
          .writeText(cleanCode(code, annotations))
          .then(() => {
            setCopyCount((count) => count + 1)
          })
      }}
    >
      <span
        aria-hidden={copied}
        className={clsx(
          'pointer-events-none flex items-center gap-0.5 text-zinc-600 transition duration-300 dark:text-zinc-400',
          copied && '-translate-y-1.5 opacity-0',
        )}
      >
        <ClipboardIcon className="size-4 fill-zinc-500/20 stroke-zinc-500 transition-colors group-hover/button:stroke-zinc-400" />
        Copy
      </span>
      <span
        aria-hidden={!copied}
        className={clsx(
          'pointer-events-none absolute inset-0 flex items-center justify-center text-primary-400 transition duration-300',
          !copied && 'translate-y-1.5 opacity-0',
        )}
      >
        Copied
      </span>
    </button>
  )
}
