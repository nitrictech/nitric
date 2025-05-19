import { visit } from 'unist-util-visit'
import { toString } from 'mdast-util-to-string'
import { remark } from 'remark'
import { slugifyWithCounter } from '@sindresorhus/slugify'

/**
 * Extracts TOC headings from markdown file and adds it to the file's data object.
 */
export const remarkTocHeadings =
  (includedDepths = [2]) =>
  () => {
    const slugify = slugifyWithCounter()

    return (tree, file) => {
      const toc = []
      visit(tree, 'heading', (node) => {
        if (!includedDepths.includes(node.depth)) return // ignore h1 and h3+

        const textContent = toString(node)
        toc.push({
          value: textContent,
          url:
            '#' +
            slugify(textContent, {
              decamelize: false,
              customReplacements: [['’', '']],
            }),
          depth: node.depth,
        })
      })
      file.data.toc = toc
    }
  }

/**
 * Passes markdown file through remark to extract TOC headings
 *
 * @param {string} markdown
 * @return {*}  {Promise<Toc>}
 */
export async function extractTocHeadings(markdown, includedDepths = [2]) {
  const vfile = await remark()
    .use(remarkTocHeadings(includedDepths))
    .process(markdown)
  // @ts-ignore
  return vfile.data.toc
}
