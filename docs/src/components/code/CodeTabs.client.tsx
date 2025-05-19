'use client'

import React from 'react'
import { TabProps, Tabs } from '../tabs/Tabs'
import useLang from '@/hooks/useLang'
import { Language, LANGUAGE_LABEL_MAP } from '@/lib/constants'

export const CodeTabsClient = React.forwardRef<HTMLDivElement, TabProps>(
  ({ children }, ref) => {
    const { currentLanguage, setCurrentLanguage } = useLang()

    return (
      <Tabs
        value={LANGUAGE_LABEL_MAP[currentLanguage]}
        onValueChange={(value) =>
          setCurrentLanguage(value.toLowerCase() as Language)
        }
        ref={ref}
      >
        {children}
      </Tabs>
    )
  },
)

CodeTabsClient.displayName = 'CodeTabs'
