import React from 'react'
import { cn } from '@/lib/utils'

type HeadingProps = React.HTMLAttributes<HTMLParagraphElement>

export const Subheading: React.FC<HeadingProps> = ({
  className,
  children,
  ...props
}) => {
  return (
    <p
      {...props}
      className={cn('mt-6 text-lg tracking-tight text-zinc-400', className)}
    >
      {children}
    </p>
  )
}
