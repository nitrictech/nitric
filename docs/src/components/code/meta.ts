import { RawCode } from 'codehike/code'

// extract meta from code block meta
export const meta = (code: RawCode) => {
  const [base, title] = code.meta.trim().split('title:')

  return {
    base: base.trim(),
    title: title ? title.trim() : null,
  }
}
