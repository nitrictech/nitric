'use client'

import clsx from 'clsx'
import React, { FC, useEffect, useState } from 'react'

const STAR_COUNT_KEY = 'nitric_github_star_count'
const STAR_COUNT_TIMESTAMP_KEY = 'nitric_github_star_count_timestamp'
const CACHE_DURATION = 1000 * 60 * 60 * 24 // 24 hours

// Function to format star count like GitHub (e.g., 1.5k, 2M)
const formatStarCount = (count: number) => {
  if (count >= 1000000) {
    return (count / 1000000).toFixed(1).replace(/\.0$/, '') + 'M'
  } else if (count >= 1000) {
    const rounded = Math.round(count / 100) / 10 // Round to the nearest decimal
    return rounded.toFixed(1).replace(/\.0$/, '') + 'k'
  }
  return count.toString()
}

interface GitHubStarCountProps {
  className: string
  defaultStarCount: number
}

const GitHubStarCount: FC<GitHubStarCountProps> = ({
  className,
  defaultStarCount,
}) => {
  const [starCount, setStarCount] = useState(defaultStarCount)

  useEffect(() => {
    const fetchStarCount = async () => {
      try {
        const cachedStarCount = localStorage.getItem(STAR_COUNT_KEY)
        const cachedTimestamp = localStorage.getItem(STAR_COUNT_TIMESTAMP_KEY)
        const currentTime = new Date().getTime()

        if (
          cachedStarCount &&
          cachedTimestamp &&
          currentTime - parseInt(cachedTimestamp) < CACHE_DURATION &&
          // Ensure cached value is at least the new default value
          parseInt(cachedStarCount) >= defaultStarCount
        ) {
          setStarCount(JSON.parse(cachedStarCount))
        } else {
          // Fetch from GitHub API if cache is expired or doesn't exist
          const response = await fetch(
            'https://api.github.com/repos/nitrictech/nitric',
          )
          if (!response.ok) {
            throw new Error('Failed to fetch star count')
          }

          const repoData = await response.json()
          const starCount = repoData.stargazers_count

          // Update state with the new star count
          setStarCount(starCount)

          // Cache the result with the current timestamp
          localStorage.setItem(STAR_COUNT_KEY, JSON.stringify(starCount))
          localStorage.setItem(STAR_COUNT_TIMESTAMP_KEY, currentTime.toString())
        }
      } catch (error) {
        console.error('Error fetching star count:', error)
      }
    }

    fetchStarCount()
  }, [])

  return (
    <div className={clsx('text-xs uppercase dark:text-zinc-300', className)}>
      {formatStarCount(starCount)}
    </div>
  )
}

export default GitHubStarCount
