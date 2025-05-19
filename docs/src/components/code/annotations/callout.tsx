// from: https://codehike.org/docs/code/callout
import { InlineAnnotation, AnnotationHandler } from 'codehike/code'

export const callout: AnnotationHandler = {
  name: 'callout',
  transform: (annotation: InlineAnnotation) => {
    const { name, query, lineNumber, fromColumn, toColumn, data } = annotation
    return {
      name,
      query,
      fromLineNumber: lineNumber,
      toLineNumber: lineNumber,
      data: { ...data, column: (fromColumn + toColumn) / 2 },
    }
  },
  Block: ({ annotation, children }) => {
    const { column } = annotation.data

    return (
      <>
        {children}
        <div
          style={{
            minWidth: `${column + 4}ch`,
            boxShadow: 'inset 0 1px 0 0 #ffffff0d',
          }}
          className="relative -ml-[1ch] mt-1 w-fit whitespace-break-spaces rounded-md border bg-zinc-700 px-2 text-zinc-200"
        >
          <div
            style={{ left: `${column}ch` }}
            className="absolute -top-[0px] size-2 -translate-y-1/2 rotate-45 border-l-4 border-t-4 border-zinc-700"
          />
          {annotation.query}
        </div>
      </>
    )
  },
}
