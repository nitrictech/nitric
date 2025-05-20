// from: https://codehike.org/docs/code/token-transitions

import { AnnotationHandler, InnerToken } from 'codehike/code'
import { SmoothPre } from './token-transitions.client'

export const tokenTransitions: AnnotationHandler = {
  name: 'token-transitions',
  PreWithRef: SmoothPre,
  Token: (props) => (
    <InnerToken merge={props} style={{ display: 'inline-block' }} />
  ),
}
