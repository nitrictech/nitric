import React, { Suspense } from 'react'
import { LanguageSwitchClient } from './LanguageSwitch.client'
import { Button } from '../ui/button'
import { LanguageIcon } from '../icons/LanguageIcon'
import { languages } from '@/lib/constants'

export const LanguageSwitch = () => {
  return (
    <Suspense
      fallback={
        <ul className="flex gap-x-4">
          {languages.map((name) => (
            <li
              key={name}
              className={
                'cursor-pointer grayscale transition-all hover:grayscale-0'
              }
            >
              <Button variant="unstyled">
                <LanguageIcon name={name} />
                <span className="sr-only">set language to {name}</span>
              </Button>
            </li>
          ))}
        </ul>
      }
    >
      <LanguageSwitchClient />
    </Suspense>
  )
}
