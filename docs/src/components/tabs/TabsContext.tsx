'use client'

import React, {
  createContext,
  PropsWithChildren,
  useCallback,
  useContext,
  useState,
} from 'react'

export interface TabsContextProps {
  set: (syncKey: string, value: string) => void
  get: (syncKey: string) => string | undefined
}

export const TabsContext = createContext<TabsContextProps>({
  set: () => {},
  get: () => '',
})

export const TabsProvider: React.FC<PropsWithChildren> = ({ children }) => {
  const [values, setValues] = useState(new Map<string, string>()) // Set the initial Tabs here

  const set = useCallback(
    (syncKey: string, value: string) => {
      setValues((prev) => {
        const next = new Map(prev)
        next.set(syncKey, value)
        return next
      })
    },
    [setValues],
  )

  const get = useCallback(
    (syncKey: string) => {
      return values.get(syncKey)
    },
    [values],
  )

  return (
    <TabsContext.Provider value={{ set, get }}>{children}</TabsContext.Provider>
  )
}

export const useTabs = (): TabsContextProps => {
  return useContext(TabsContext)
}
