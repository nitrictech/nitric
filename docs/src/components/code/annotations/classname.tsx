import { AnnotationHandler } from 'codehike/code'

export const className: AnnotationHandler = {
  name: 'className',
  Block: ({ annotation, children }) => (
    <div className={annotation.query}>{children}</div>
  ),
  Inline: ({ annotation, children }) => (
    <span className={annotation.query}>{children}</span>
  ),
}
