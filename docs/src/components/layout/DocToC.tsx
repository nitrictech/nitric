'use client'

import type { Doc, Guide } from '@/content'
import { useState, useEffect } from 'react'
import Link from 'next/link'
import { DocTracingBeam } from './DocTracingBeam'
import { Button } from '../ui/button'
import { GitHubIcon } from '../icons/GitHubIcon'
import { useMotionValue } from 'framer-motion'
import { cn } from '@/lib/utils'

interface Toc {
  url: string
  value: string
  depth: number
}

const DocToC = ({
  doc,
  startSteps,
}: {
  doc: Doc | Guide
  startSteps?: React.JSX.Element | undefined
}) => {
  const initial = 14
  const sectionSize = 28
  const offset = 10

  const y1 = useMotionValue(0)
  const y2 = useMotionValue(0)

  const [activeIndex, setActiveIndex] = useState(-1)

  useEffect(() => {
    const headings = Array.from(
      document.querySelectorAll<HTMLHeadingElement>('h2.md-content-header'),
    )

    const options = {
      root: null,
      rootMargin: '0px',
      threshold: 1, // Adjust based on when you want to highlight
    }

    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          const index = headings.findIndex(
            (heading) => heading.textContent === entry.target.textContent,
          )

          y2.set(initial + (index * sectionSize + offset))

          setActiveIndex(index)
        }
      })
    }, options)

    headings.forEach((h2) => observer.observe(h2))

    return () => {
      headings.forEach((h2) => observer.unobserve(h2))
    }
  }, [])

  return (
    <div className="hidden h-full min-w-64 md:block">
      <aside className="sticky top-[calc(var(--header-height)+1px+2rem)] max-h-[calc(100vh-var(--header-height)-3rem)] space-y-8">
        {doc.toc.length ? (
          <div className="relative flex flex-col">
            <p className="mb-2 font-mono text-sm uppercase dark:text-zinc-300">
              On this page
            </p>
            <DocTracingBeam y1={y1} y2={y2}>
              <ol className="flex flex-col gap-y-1 pl-4 text-sm font-medium">
                {doc.toc.map((item: Toc, i: number) => {
                  return (
                    <li key={item.url + i}>
                      <Link
                        href={item.url}
                        className={cn(
                          'text-2xs',
                          activeIndex === i
                            ? 'text-zinc-900 dark:text-zinc-200'
                            : 'text-muted-foreground transition-colors hover:text-zinc-900 dark:hover:text-zinc-200',
                        )}
                      >
                        {item.value}
                      </Link>
                    </li>
                  )
                })}
              </ol>
            </DocTracingBeam>
          </div>
        ) : null}
        <div>
          <Button asChild variant="unstyled">
            <a
              href={doc.editUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="group flex items-center gap-x-2 text-xs dark:text-zinc-300 dark:hover:text-white"
            >
              <GitHubIcon className="h-5 w-5 fill-zinc-800 transition group-hover:fill-zinc-900 dark:fill-zinc-300 dark:group-hover:fill-white" />
              Edit this page on GitHub
            </a>
          </Button>
        </div>
        {startSteps && (
          <div className="flex flex-col">
            <p className="mb-2 font-mono text-sm uppercase dark:text-zinc-300">
              Try this out
            </p>
            <p className="mb-2 text-2xs text-muted-foreground">
              Run these commands in your terminal.
            </p>
            {startSteps}
          </div>
        )}
      </aside>
    </div>
  )
}

export default DocToC
