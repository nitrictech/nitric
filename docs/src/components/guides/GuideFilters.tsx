import React, { Suspense } from 'react'
import { LanguageSwitch } from '../LanguageSwitch'
import GuideFilterCheckbox from './GuideFilterCheckbox'

interface Props {
  allTags: string[]
}

export const GuideFilters: React.FC<Props> = ({ allTags }) => {
  return (
    <>
      <LanguageSwitch />
      <ul className="space-y-4">
        {allTags.map((tag) => (
          <li key={tag}>
            <Suspense>
              <GuideFilterCheckbox tag={tag} />
            </Suspense>
          </li>
        ))}
      </ul>
    </>
  )
}
