'use client'

import { cn } from '@/lib/utils'
import React from 'react'
import { Button } from '../ui/button'
import { languages } from '@/lib/constants'
import { LanguageIcon } from '../icons/LanguageIcon'
import useParams from '@/hooks/useParams'

export const LanguageSwitchClient = () => {
  const { setParams, searchParams } = useParams()
  const selectedLangs = searchParams?.get('langs')?.split(',') || []

  const handleLanguageChange = (lang: string) => {
    if (!selectedLangs.includes(lang)) {
      setParams('langs', [...selectedLangs, lang].join(','))
    } else {
      setParams(
        'langs',
        selectedLangs.filter((selected) => selected !== lang).join(','),
      )
    }
  }

  return (
    <ul className="flex gap-x-4">
      {languages.map((name) => (
        <li
          key={name}
          className={cn(
            'cursor-pointer transition-all hover:grayscale-[50%]',
            !selectedLangs.includes(name) ? 'grayscale' : '',
          )}
        >
          <Button variant="unstyled" onClick={() => handleLanguageChange(name)}>
            <LanguageIcon name={name} />
            <span className="sr-only">set language to {name}</span>
          </Button>
        </li>
      ))}
    </ul>
  )
}
