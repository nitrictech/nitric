'use client'

import React, { useEffect, useId } from 'react'
import { Collapsible, CollapsibleTrigger } from '../ui/collapsible'
import { ChevronDownIcon } from '@heroicons/react/24/outline'
import { cn } from '@/lib/utils'
import { NavLink } from './NavLink'
import { usePathname } from 'next/navigation'
import { Button } from '../ui/button'
import { NavEntry, NavGroup } from '@/config/types'

interface Props {
  group: NavGroup
  className?: string
}

const checkIfActive = (items: NavEntry[], pathname: string): boolean => {
  return items.some((link) => {
    if ('href' in link && link.href === pathname) {
      return true
    }

    if ('items' in link && Array.isArray(link.items)) {
      return checkIfActive(link.items, pathname)
    }

    return false
  })
}

const CollapsibleNavItem: React.FC<Props> = ({ group, className }) => {
  const pathname = usePathname()

  // overriding radix id as we don't need their CollapsibleContent component
  const navItemId = useId()

  const { title, items, icon: Icon } = group

  const [isOpen, setIsOpen] = React.useState(checkIfActive(items, pathname))

  useEffect(() => {
    const isActive = checkIfActive(items, pathname)

    // only open if the group is active
    if (isActive) {
      setIsOpen(isActive)
    }
  }, [pathname])

  return (
    <Collapsible
      open={isOpen}
      onOpenChange={setIsOpen}
      className={cn('group/nav-collapse space-y-2 pl-2', className)}
    >
      <CollapsibleTrigger asChild aria-controls={navItemId}>
        <Button
          variant="link"
          className={cn(
            'flex w-full justify-between gap-2 py-1 pl-0 pr-3 text-sm transition hover:no-underline',
            isOpen
              ? 'text-zinc-900 dark:text-white'
              : 'text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-white',
          )}
        >
          <div className="flex items-center">
            {Icon && <Icon className="mr-2 h-4 w-4" />}
            <span>{title}</span>
          </div>
          <ChevronDownIcon className={cn('h-4 w-4', isOpen && 'rotate-180')} />
        </Button>
      </CollapsibleTrigger>

      <div
        id={navItemId}
        hidden={!isOpen}
        className="relative mt-2 space-y-2 overflow-hidden group-data-[state='closed']/nav-collapse:h-0 group-data-[state='open']/nav-collapse:h-auto"
      >
        <ul role="list">
          {items.map((entry, idx) =>
            'href' in entry ? (
              <li key={entry.title}>
                <NavLink
                  isAnchorLink
                  key={entry.title}
                  href={entry.href}
                  active={entry.href === pathname}
                >
                  {entry.title}
                </NavLink>
              </li>
            ) : (
              <li key={entry.title}>
                <div
                  aria-hidden="true"
                  className={cn(
                    'absolute bottom-0 left-2 top-0 w-[1px]',
                    'bg-zinc-300 dark:bg-zinc-700',
                  )}
                />{' '}
                <CollapsibleNavItem group={entry} className="pl-7" />
              </li>
            ),
          )}
        </ul>
      </div>
    </Collapsible>
  )
}

export default CollapsibleNavItem
