import React from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

type Level = 1 | 2 | 3 | 4 | 5 | 6

const headingVariants = cva('text-balance font-bold font-display', {
  variants: {
    level: {
      1: 'text-foreground text-3xl sm:text-4xl lg:text-5xl',
      2: 'text-foreground text-2xl sm:text-3xl font-bold lg:text-4xl',
      3: 'text-foreground text-2xl lg:text-3xl',
      4: 'text-foreground text-xl lg:text-2xl',
      5: 'text-foreground text-lg lg:text-xl',
      6: 'text-foreground text-base lg:text-lg',
    },
  },
  defaultVariants: {
    level: 6,
  },
})

interface HeadingProps
  extends React.HTMLAttributes<HTMLHeadingElement>,
    VariantProps<typeof headingVariants> {}

export const Heading: React.FC<HeadingProps> = ({
  level,
  className,
  children,
  ...props
}) => {
  const Tag = `h${level}` as `h${Level}`

  return (
    <Tag {...props} className={cn(headingVariants({ level, className }))}>
      {children}
    </Tag>
  )
}
