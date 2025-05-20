import { HighlightedCode, Pre as CodeHikePre } from 'codehike/code'
import React from 'react'
import { callout } from './annotations/callout'
import { CopyButton } from './CopyButton'
import { fold } from './annotations/fold'
import {
  collapse,
  collapseContent,
  collapseTrigger,
} from './annotations/collapse'
import { tokenTransitions } from './annotations/token-transitions'
import { cn } from '@/lib/utils'
import { meta } from './meta'
import { mark } from './Mark'
import { diff } from './Diff'

export interface HandlerProps {
  enableTransitions?: boolean
}

type Props = {
  highlighted: HighlightedCode
  showPanel?: boolean
  className?: string
  hideBashPanel?: boolean
} & HandlerProps

const Pre: React.FC<Props> = ({
  highlighted,
  showPanel,
  hideBashPanel,
  enableTransitions,
  className,
}) => {
  const { title } = meta(highlighted)

  let handlers = [
    callout,
    fold,
    collapse,
    collapseTrigger,
    collapseContent,
    mark,
    diff,
  ]

  // Getting an issue with color transitions on the -code-text token
  // Note: this also won't work with the Tabs component, only the CodeSwitcherSelect
  if (enableTransitions) {
    handlers = [...handlers, tokenTransitions]
  }

  const showFileNamePanel = showPanel && !!title

  const isBash = highlighted.lang === 'shellscript' && !hideBashPanel

  return (
    <>
      {isBash && !showFileNamePanel && (
        <div className="relative flex h-9 items-center justify-start pr-12 font-display text-2xs text-zinc-300 sm:text-xs">
          {/* one-off breakpoint to hide the filename on extremely narrow screens - to avoid interfering with the lang select */}
          <div className="relative flex h-full items-center justify-center gap-x-1 px-4">
            <div className="size-[6px] rounded-full bg-zinc-300 dark:bg-zinc-300/10" />
            <div className="size-[6px] rounded-full bg-zinc-300 dark:bg-zinc-300/10" />
            <div className="size-[6px] rounded-full bg-zinc-300 dark:bg-zinc-300/10" />
          </div>
          <div className="absolute bottom-0 h-px w-full bg-zinc-200 dark:bg-zinc-300/5" />
        </div>
      )}
      {showFileNamePanel && (
        <div className="relative flex h-10 items-center justify-start pr-12 font-display text-2xs dark:text-zinc-300 sm:text-xs">
          {/* one-off breakpoint to hide the filename on extremely narrow screens - to avoid interfering with the lang select */}
          <div className="relative flex h-full items-center border-r dark:border-zinc-300/5">
            <span className="hidden whitespace-nowrap px-4 py-2 min-[320px]:block">
              {title}
            </span>
            <div className="absolute bottom-0 z-10 h-px w-full bg-primary-300" />
          </div>
          <div className="absolute bottom-0 h-px w-full bg-zinc-200 dark:bg-zinc-300/5" />
        </div>
      )}
      <CopyButton
        code={highlighted.code}
        annotations={highlighted.annotations}
        className={cn(showPanel && 'top-12', isBash && 'top-11')}
      />
      <CodeHikePre
        code={highlighted}
        handlers={handlers}
        className={cn(
          'overflow-auto overscroll-x-contain p-4 text-sm',
          showPanel && !title && 'pt-7', // add padding to ensure the code doesn't touch the top of the panel
          className,
        )}
        style={{
          ...highlighted.style,
          fontSize: '0.875rem',
          background: 'transparent !important',
        }}
      />
    </>
  )
}

export default Pre
