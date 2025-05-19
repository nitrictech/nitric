import React from 'react'
import * as mdxComponents from '@/components/mdx'
import { compile, run } from '@mdx-js/mdx'
import * as runtime from 'react/jsx-runtime'
import { MDXComponents } from 'mdx/types'
import { mdxOptions } from '@/mdx/mdx-options.mjs'

interface MDXContentProps {
  mdx: string
}

const MDXContent = async ({ mdx }: MDXContentProps) => {
  /// Compile the MDX source code to a function body
  const code = String(
    await compile(mdx, {
      outputFormat: 'function-body',
      ...mdxOptions,
      jsx: false, // Disable JSX transformation, as it's already transformed by the runtime
    }),
  )

  // Run the compiled code with the runtime and get the default export
  const { default: MDXContent } = await run(code, {
    ...runtime,
    baseUrl: import.meta.url,
  } as any)

  return <MDXContent components={mdxComponents as MDXComponents} />
}

export default MDXContent
