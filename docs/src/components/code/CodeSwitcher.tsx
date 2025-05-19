import React, { Suspense } from 'react'
import { RawCode } from 'codehike/code'
import CodeContainer from './CodeContainer'
import { CodeSwitcherSelect } from './CodeSwitcherSelect'
import Pre, { HandlerProps } from './Pre'
import { highlight } from './highlight'
import { TabItem, Tabs } from '../tabs/Tabs'
import { LANGUAGE_LABEL_MAP, languages } from '@/lib/constants'
import { CodeTabs } from './CodeTabs'
import { meta } from './meta'

export async function CodeSwitcher({
  code,
  tabs,
  ...props
}: {
  code: RawCode[]
  tabs?: boolean
  showPanel?: boolean
  className?: string
} & HandlerProps) {
  const highlighted = await Promise.all(
    code.map((codeblock) => highlight(codeblock)),
  )

  const missingLangs = languages.filter(
    (lang) => !highlighted.some((h) => h.lang === lang),
  )

  if (missingLangs.length) {
    throw Error(
      `CodeSwitcher missing languages: ${missingLangs.join(', ')} at: ${highlighted[0].meta || 'unknown'}`,
    )
  }

  if (tabs) {
    const children = highlighted.map((h) => (
      <TabItem key={h.lang} label={LANGUAGE_LABEL_MAP[h.lang]}>
        <CodeContainer>
          <Pre highlighted={h} {...props} showPanel={!!meta(h).title} />
        </CodeContainer>
      </TabItem>
    ))

    return (
      <CodeTabs fallback={<Tabs {...props}>{children}</Tabs>}>
        {children}
      </CodeTabs>
    )
  }

  return (
    <CodeContainer>
      <Suspense
        fallback={<Pre highlighted={highlighted[0]} {...props} showPanel />}
      >
        <CodeSwitcherSelect highlighted={highlighted} {...props} />
      </Suspense>
    </CodeContainer>
  )
}
