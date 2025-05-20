import React, { PropsWithChildren } from 'react'
import Link from 'next/link'
import { cn } from '@/lib/utils'
import { Tag } from '../Tag'
import { Button } from '../ui/button'
import { NavItem } from '@/config/types'

interface NavLinkProps extends PropsWithChildren<Omit<NavItem, 'title'>> {
  href: string
  tag?: string
  active?: boolean
  isAnchorLink?: boolean
  className?: string
}

export const NavLink: React.FC<NavLinkProps> = ({
  active,
  isAnchorLink,
  href,
  children,
  tag,
  icon: Icon,
  className,
}) => {
  const isExternal = href.startsWith('http')

  return (
    <Button variant="link" asChild>
      <Link
        href={href}
        aria-current={active ? 'page' : undefined}
        target={isExternal ? '_blank' : undefined}
        rel={isExternal ? 'noopener noreferrer' : undefined}
        className={cn(
          'relative flex w-full justify-between gap-2 py-1 pr-3 text-sm no-underline transition hover:no-underline',
          isAnchorLink ? 'pl-7' : 'pl-2',
          active
            ? 'bg-zinc-100/60 text-zinc-900 dark:bg-zinc-800/60 dark:text-white'
            : 'text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-white',
          className,
        )}
      >
        {isAnchorLink && (
          <div
            aria-hidden="true"
            className={cn(
              'absolute bottom-0 left-2 top-0 z-10 w-[1px]',
              active ? 'bg-primary-500' : 'bg-zinc-300 dark:bg-zinc-700',
            )}
          />
        )}
        <span className="flex items-center truncate">
          {Icon && <Icon className="mr-2 h-4 w-4" />} {children}
        </span>
        {tag && (
          <Tag variant="small" color="zinc">
            {tag}
          </Tag>
        )}
      </Link>
    </Button>
  )
}
