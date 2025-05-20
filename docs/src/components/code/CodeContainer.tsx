import React from 'react'

const CodeContainer: React.FC<React.PropsWithChildren> = ({ children }) => {
  return (
    <div className="w-full max-w-full">
      <div
        tabIndex={0}
        className={
          'not-prose group relative mb-6 w-full max-w-full overflow-hidden rounded-md border bg-code shadow-md shadow-zinc-300/5 dark:border-white/5 dark:shadow-zinc-800/10'
        }
      >
        {children}
      </div>
    </div>
  )
}

export default CodeContainer
