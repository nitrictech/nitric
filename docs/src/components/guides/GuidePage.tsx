import React, { Suspense } from 'react'

import { GuideFilters } from './GuideFilters'
import GuideList from './GuideList'
import GuideMobileFilters from './GuideMobileFilters'
import { allGuides } from '@/content'

interface Props {
  allTags: string[]
}

const GuidePage: React.FC<Props> = ({ allTags }) => {
  return (
    <div className="gap-x-4 lg:grid lg:grid-cols-[280px,1fr]">
      <div className="hidden border-r pb-10 lg:block">
        <aside
          aria-label="Sidebar"
          className="mt-10 w-80 space-y-10 overflow-y-auto"
        >
          <GuideFilters allTags={allTags} />
        </aside>
      </div>
      <div className="relative -top-5 lg:hidden">
        <GuideMobileFilters allTags={allTags} />
      </div>
      <Suspense>
        <GuideList
          allGuides={allGuides}
          className="relative mx-2 my-4 w-full sm:px-8 lg:mb-8"
        />
      </Suspense>
    </div>
  )
}

export default GuidePage
