import { recmaPlugins } from './recma.mjs'
import { rehypePlugins } from './rehype.mjs'
import { remarkPlugins } from './remark.mjs'
import { recmaCodeHike, remarkCodeHike } from 'codehike/mdx'

/** @type {import('codehike/mdx').CodeHikeConfig} */
const chConfig = {
  components: { code: 'Code' },
}

export const mdxOptions = {
  remarkPlugins: [...remarkPlugins, [remarkCodeHike, chConfig]],
  rehypePlugins,
  recmaPlugins: [...recmaPlugins, [recmaCodeHike, chConfig]],
  jsx: true,
}
