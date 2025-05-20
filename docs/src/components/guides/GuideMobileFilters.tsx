import React from 'react'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../ui/dialog'
import { Button } from '../ui/button'
import { AdjustmentsHorizontalIcon } from '@heroicons/react/24/outline'
import { GuideFilters } from './GuideFilters'

interface GuideMobileFiltersProps {
  allTags: string[]
}

const GuideMobileFilters: React.FC<GuideMobileFiltersProps> = ({ allTags }) => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="icon" aria-label="Apply guide filters">
          <AdjustmentsHorizontalIcon className="h-6 w-6" />
        </Button>
      </DialogTrigger>
      <DialogContent className="max-h-[calc(100vh-5rem)] overflow-y-auto">
        <DialogHeader className="mb-10">
          <DialogTitle className="sr-only">Set Filters</DialogTitle>
        </DialogHeader>
        <GuideFilters allTags={allTags} />
        <DialogClose asChild>
          <Button className="mt-10 w-full">Apply</Button>
        </DialogClose>
      </DialogContent>
    </Dialog>
  )
}

export default GuideMobileFilters
