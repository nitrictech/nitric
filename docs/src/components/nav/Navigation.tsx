'use client'

import { usePathname } from 'next/navigation'

import { NavLink } from './NavLink'
import { cn } from '@/lib/utils'
import { NavigationGroup } from './NavigationGroup'
import { navigation } from '@/config'

export function Navigation(props: React.ComponentPropsWithoutRef<'nav'>) {
  const pathname = usePathname()

  return (
    <nav {...props}>
      <ul role="list">
        {navigation.map((entry, groupIndex) =>
          'href' in entry ? (
            <li key={entry.title}>
              <NavLink
                className={cn('pl-2', groupIndex === 0 ? 'md:mt-0' : '')}
                active={entry.href === pathname}
                {...entry}
              >
                {entry.title}
              </NavLink>
            </li>
          ) : (
            <NavigationGroup
              key={entry.title}
              group={entry}
              className={groupIndex === 0 ? 'md:mt-0' : ''}
            />
          ),
        )}
      </ul>
    </nav>
  )
}
