import { highlight as codehikeHighlight, RawCode } from 'codehike/code'
import CODE_THEME from './theme'

const cleanCode = (code: RawCode) => {
  // Replace tabs with two spaces
  code.value = code.value.replace(/\t/g, '  ')

  return code
}

export const highlight = (data: RawCode) =>
  codehikeHighlight(cleanCode(data), CODE_THEME)
