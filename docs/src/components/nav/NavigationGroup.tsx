import { cn } from '@/lib/utils'
import { usePathname } from 'next/navigation'
import { NavLink } from './NavLink'
import CollapsibleNavItem from './CollapsibleNavItem'
import { NavGroup } from '@/config/types'

export function NavigationGroup({
  group,
  className,
  isAnchorLinks,
}: {
  group: NavGroup
  className?: string
  isAnchorLinks?: boolean
}) {
  // If this is the mobile navigation then we always render the initial
  // state, so that the state does not change during the close animation.
  // The state will still update when we re-open (re-render) the navigation.
  const pathname = usePathname()

  return (
    <li className={cn('relative mt-6', className)}>
      <h2 className="pl-2 text-2xs font-semibold text-zinc-900 dark:text-white">
        {group.title}
      </h2>
      <div className="relative mt-3">
        <ul role="list" className="border-l border-transparent">
          {group.items.map((link) =>
            'href' in link ? (
              <li key={link.href} className="relative">
                <NavLink
                  {...link}
                  active={link.href === pathname}
                  isAnchorLink={isAnchorLinks}
                >
                  {link.title}
                </NavLink>
              </li>
            ) : (
              <li key={link.title}>
                <CollapsibleNavItem group={link} />
              </li>
            ),
          )}
        </ul>
      </div>
    </li>
  )
}
