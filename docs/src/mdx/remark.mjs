import { mdxAnnotations } from 'mdx-annotations'
import remarkGfm from 'remark-gfm'
import mdxMermaid from 'mdx-mermaid'

export const remarkPlugins = [
  mdxAnnotations.remark,
  remarkGfm,
  [
    mdxMermaid,
    {
      output: 'svg',
      mermaid: {
        theme: 'base',
        // TODO: Relocate theme config
        themeVariables: {
          background: 'white',
          primaryColor: '#F9F3FF',
          primaryBorderColor: 'var(--secondary-300)',
          lineColor: '#000000',
          secondaryColor: '#ffffff',
          tertiaryColor: '#0000ff',
          primaryTextColor: '#000000',
          fontSize: '24px', // use with styles in mermaid.css, this zooms out the diagram
          fontFamily: 'var(--font-jetbrains-mono), monospace',
        },
      },
    },
  ],
]
