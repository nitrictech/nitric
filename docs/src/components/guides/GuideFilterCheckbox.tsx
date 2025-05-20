'use client'

import React from 'react'
import { Checkbox } from '../ui/checkbox'
import useParams from '@/hooks/useParams'
import { dash } from 'radash'

interface GuideFilterCheckboxProps {
  tag: string
}

const GuideFilterCheckbox: React.FC<GuideFilterCheckboxProps> = ({ tag }) => {
  const { searchParams, setParams } = useParams()
  const selectedTags = searchParams.get('tags')?.split(',') || []

  return (
    <div className="flex items-center space-x-4">
      <Checkbox
        id={dash(tag)}
        aria-label={tag}
        checked={selectedTags.includes(tag)}
        onCheckedChange={(checked) => {
          if (checked) {
            setParams('tags', [...selectedTags, tag].join(','))
          } else {
            setParams(
              'tags',
              selectedTags
                .filter((selectedTag) => selectedTag !== tag)
                .join(','),
            )
          }
        }}
        className="h-5 w-5 border-primary-400 data-[state=checked]:bg-primary"
      />
      <label
        htmlFor={dash(tag)}
        className="cursor-pointer text-base font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        {tag}
      </label>
    </div>
  )
}

export default GuideFilterCheckbox
