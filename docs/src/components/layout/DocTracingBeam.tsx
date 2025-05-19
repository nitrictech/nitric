'use client'

import React, { useEffect, useRef, useState } from 'react'
import { motion, useSpring, MotionValue } from 'framer-motion'
import { cn } from '@/lib/utils'

export const DocTracingBeam = ({
  children,
  className,
  y1,
  y2,
}: {
  children: React.ReactNode
  className?: string
  y1: MotionValue<number>
  y2: MotionValue<number>
}) => {
  const contentRef = useRef<HTMLDivElement>(null)
  const [svgHeight, setSvgHeight] = useState(0)

  useEffect(() => {
    if (contentRef.current) {
      setSvgHeight(contentRef.current.offsetHeight)
    }
  }, [])

  return (
    <motion.div
      className={cn('relative mx-auto h-full w-full max-w-4xl', className)}
    >
      <div className="absolute -left-4 top-1 mb-2 ml-1 md:-left-16">
        <svg
          viewBox={`0 0 20 ${svgHeight}`}
          width="20"
          height={svgHeight} // Set the SVG height
          className="ml-10 block"
          aria-hidden="true"
        >
          <motion.path
            d={`M 1 0V -36 l 18 24 V ${svgHeight * 0.98}`}
            fill="none"
            className={`stroke-zinc-400 dark:stroke-zinc-200`}
            strokeOpacity="0.16"
            transition={{
              duration: 10,
            }}
          ></motion.path>
          <motion.path
            d={`M 1 0V -36 l 18 24 V ${svgHeight * 0.98}`}
            fill="none"
            stroke="url(#gradient)"
            strokeWidth="1.25"
            className="motion-reduce:hidden"
            transition={{
              duration: 10,
            }}
          ></motion.path>
          <defs>
            <motion.linearGradient
              id="gradient"
              gradientUnits="userSpaceOnUse"
              x1="0"
              x2="0"
              y1={y1} // set y1 for gradient
              y2={useSpring(y2)} // set y2 for gradient
            >
              <stop stopColor="var(--primary-500)" stopOpacity="0"></stop>
              <stop stopColor="var(--primary-500)"></stop>
              <stop offset="0.325" stopColor="var(--secondary-700)"></stop>
              <stop
                offset="1"
                stopColor="var(--secondary-500)"
                stopOpacity="0"
              ></stop>
            </motion.linearGradient>
          </defs>
        </svg>
      </div>
      <div ref={contentRef}>{children}</div>
    </motion.div>
  )
}
