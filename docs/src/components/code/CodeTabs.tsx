import React, { Suspense } from 'react'
import { TabProps, Tabs } from '../tabs/Tabs'
import { CodeTabsClient } from './CodeTabs.client'

export const CodeTabs: React.FC<TabProps & { fallback?: React.ReactNode }> = ({
  children,
  fallback,
  ...props
}) => {
  return (
    <Suspense fallback={fallback}>
      <CodeTabsClient {...props}>{children}</CodeTabsClient>
    </Suspense>
  )
}
