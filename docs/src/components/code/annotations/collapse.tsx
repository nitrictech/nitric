import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import { AnnotationHandler, BlockAnnotation, InnerLine } from 'codehike/code'
import { ChevronDownIcon } from 'lucide-react'

export const collapse: AnnotationHandler = {
  name: 'collapse',
  transform: (annotation: BlockAnnotation) => {
    const { fromLineNumber } = annotation
    return [
      annotation,
      {
        ...annotation,
        fromLineNumber: fromLineNumber,
        toLineNumber: fromLineNumber,
        name: 'CollapseTrigger',
      },
      {
        ...annotation,
        fromLineNumber: fromLineNumber + 1,
        name: 'CollapseContent',
      },
    ]
  },
  Block: ({ annotation, children }) => {
    return (
      <Collapsible defaultOpen={annotation.query !== 'collapsed'}>
        {children}
      </Collapsible>
    )
  },
}

const icon = (
  <ChevronDownIcon
    className="mb-0.5 inline-block select-none opacity-30 transition group-hover:!opacity-100 group-data-[state=closed]:-rotate-90 group-data-[state=closed]:opacity-80"
    size={15}
  />
)
export const collapseTrigger: AnnotationHandler = {
  name: 'CollapseTrigger',
  onlyIfAnnotated: true,
  AnnotatedLine: ({ annotation, ...props }) => (
    <CollapsibleTrigger
      className="group contents"
      aria-label="Toggle show code"
    >
      <InnerLine merge={props} data={{ icon }} />
    </CollapsibleTrigger>
  ),
  Line: (props) => {
    const icon = props.data?.icon as React.ReactNode
    return (
      <div className="table-row">
        <span className="table-cell min-w-5 text-center">{icon}</span>
        <div className="table-cell">
          <InnerLine merge={props} />
        </div>
      </div>
    )
  },
}

export const collapseContent: AnnotationHandler = {
  name: 'CollapseContent',
  Block: CollapsibleContent,
}
