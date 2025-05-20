'use client'

import React, { useEffect } from 'react'
import { TabProps, Tabs } from './tabs/Tabs'
import { useTabs } from './tabs/TabsContext'

export type OS = 'Windows' | 'macOS' | 'Linux'

// this function detects the operating system of the user and returns the name that will be used to sort the instructions
function detectOS() {
  const userAgent =
    typeof window !== 'undefined' ? window.navigator.userAgent : ''
  if (userAgent.includes('Mac')) {
    return 'macOS'
  } else if (userAgent.includes('Win')) {
    return 'Windows'
  }

  return 'Linux'
}

export const OSTabs = React.forwardRef<HTMLDivElement, TabProps>(
  ({ children }, ref) => {
    const { set } = useTabs()

    useEffect(() => {
      // this will only run on the client side
      set('operating-system', detectOS())
    }, [])

    return (
      <Tabs syncKey="operating-system" ref={ref}>
        {children}
      </Tabs>
    )
  },
)

OSTabs.displayName = 'OSTabs'
