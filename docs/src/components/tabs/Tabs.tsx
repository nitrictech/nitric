'use client'

import React, { PropsWithChildren, ReactElement } from 'react'
import {
  Tabs as BaseTabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import { useTabs } from './TabsContext'

export interface TabProps extends PropsWithChildren {
  syncKey?: string
  value?: string
  ref?: React.Ref<HTMLDivElement>
  onValueChange?: (value: string) => void
}

export interface TabItemProps extends PropsWithChildren {
  label: string
}

// for a11y, remove spaces from the tab label
const removeSpaces = (str: string) => str.replace(/\s/g, '')

export const Tabs: React.FC<TabProps> = ({
  children,
  value,
  onValueChange,
  ref,
  syncKey,
}) => {
  const tabs = React.Children.toArray(children) as ReactElement<TabItemProps>[]

  const { set, get } = useTabs()

  return (
    <BaseTabs
      defaultValue={tabs[0] ? removeSpaces(tabs[0].props.label) : undefined}
      value={syncKey ? get(syncKey) : value}
      onValueChange={syncKey ? (value) => set(syncKey, value) : onValueChange}
      ref={ref}
    >
      <TabsList className="relative mx-0 mb-2 mt-auto h-12 w-full rounded-b-none bg-transparent p-0">
        {tabs.map((tab) => (
          <TabsTrigger
            value={removeSpaces(tab.props.label)}
            key={removeSpaces(tab.props.label)}
            className="group/tab relative h-12 hover:text-zinc-600 data-[state=active]:bg-transparent data-[state=active]:text-primary data-[state=active]:shadow-none dark:hover:text-zinc-200 dark:data-[state=active]:text-primary-300"
          >
            {tab.props.label}
            <div className="absolute inset-x-0 bottom-0 z-10 h-px bg-primary opacity-0 transition-opacity group-data-[state=active]/tab:opacity-100 dark:bg-primary-300" />
          </TabsTrigger>
        ))}
        <div className="absolute inset-x-0 bottom-0 h-px bg-zinc-200 dark:bg-zinc-300/10" />
      </TabsList>
      {children}
    </BaseTabs>
  )
}

export const TabItem: React.FC<TabItemProps> = ({ children, label }) => {
  return <TabsContent value={removeSpaces(label)}>{children}</TabsContent>
}
