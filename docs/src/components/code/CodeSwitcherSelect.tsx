'use client'

import React from 'react'
import { HighlightedCode } from 'codehike/code'
import Pre, { HandlerProps } from './Pre'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select'
import useLang from '@/hooks/useLang'
import { cn } from '@/lib/utils'
import { Language } from '@/lib/constants'

const languageMap: Record<string, string> = {
  javascript: 'JavaScript',
  python: 'Python',
  typescript: 'TypeScript',
  go: 'Go',
  dart: 'Dart',
  java: 'JVM',
  Kotlin: 'Kotlin',
  csharp: 'C#',
}

export function CodeSwitcherSelect({
  highlighted,
  showPanel = true,
  ...props
}: {
  highlighted: HighlightedCode[]
  className?: string
  showPanel?: boolean
} & HandlerProps) {
  const { currentLanguage, setCurrentLanguage } = useLang()

  const selectedCode = highlighted.find(
    (code) => code.lang === currentLanguage,
  )!

  return (
    <>
      {showPanel && (
        <>
          <Select
            value={currentLanguage}
            onValueChange={(lang) => {
              // This is necessary to fix a ui update delay causing 200ms lag
              setTimeout(() => setCurrentLanguage(lang as Language), 0)
            }}
          >
            <SelectTrigger
              aria-label="Switch code language"
              className="absolute right-3.5 top-[2.5px] z-10 hidden h-8 w-28 bg-white/5 text-2xs font-medium text-zinc-400 ring-1 ring-inset ring-zinc-300/10 hover:bg-white/7.5 dark:bg-white/2.5 dark:hover:bg-white/5 md:flex"
            >
              <SelectValue />
            </SelectTrigger>
            <SelectContent position="item-aligned">
              {highlighted.map(({ lang }, index) => (
                <SelectItem key={index} value={lang}>
                  {languageMap[lang] || lang}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <select
            aria-label="Switch code language"
            onChange={(e) => setCurrentLanguage(e.target.value as Language)}
            className={cn(
              'flex h-8 w-full items-center justify-between whitespace-nowrap rounded-md border border-input bg-transparent px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50',
              'absolute right-3.5 top-[2.5px] z-10 w-24 bg-white/5 font-sans text-2xs font-medium text-zinc-400 ring-1 ring-inset ring-zinc-300/10 hover:bg-white/7.5 dark:bg-white/2.5 dark:hover:bg-white/5 md:hidden',
            )}
            value={currentLanguage}
          >
            {highlighted.map(({ lang }, index) => (
              <option key={index} value={lang} className="font-sans text-sm">
                {languageMap[lang] || lang}
              </option>
            ))}
          </select>
        </>
      )}
      <Pre highlighted={selectedCode} {...props} showPanel={showPanel} />
    </>
  )
}
