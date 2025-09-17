import {
  ComputedFields,
  defineDocumentType,
  defineNestedType,
  FieldDefs,
  makeSource,
} from 'contentlayer2/source-files'
import { extractTocHeadings } from './src/mdx/remark-toc-headings.mjs'
import path from 'path'
import fs from 'fs'

const contentDirPath = 'docs'

const branch = process.env.NEXT_PUBLIC_GITHUB_BRANCH || 'main'

const baseFields: FieldDefs = {
  title_seo: {
    type: 'string',
    description:
      'The meta title of the doc, this will override the title extracted from the markdown and the nav title',
  },
  description: {
    type: 'string',
    description: 'The description of the doc',
  },
  image: {
    type: 'string',
    description: 'The image of the doc',
  },
  image_alt: {
    type: 'string',
    description: 'The image alt of the doc',
  },
  disable_edit: {
    type: 'boolean',
    description: 'Disable the github edit button',
  },
  canonical_url: {
    type: 'string',
    description: 'The canonical url of the doc, if different from the url',
  },
  noindex: {
    type: 'boolean',
    description: 'Prevent search engines from indexing this page',
  },
}

const computedFields: ComputedFields = {
  slug: {
    type: 'string',
    resolve: (doc) => doc._raw.flattenedPath,
  },
  toc: { type: 'json', resolve: (doc) => extractTocHeadings(doc.body.raw) },
  title: {
    type: 'string',
    resolve: async (doc) => {
      const headings = await extractTocHeadings(doc.body.raw, [1])

      return headings[0]?.value
    },
  },
  editUrl: {
    type: 'string',
    resolve: (doc) =>
      `https://github.com/nitrictech/nitric/edit/${branch}/docs/docs/${doc._raw.sourceFilePath}`,
  },
  lastModified: {
    type: 'date',
    resolve: (doc) => {
      // Get the full path to the markdown file
      const filePath = path.join(
        process.cwd(),
        contentDirPath,
        doc._raw.sourceFilePath,
      )
      // Extract and return the last modified date
      const stats = fs.statSync(filePath)
      return stats.mtime // This is the last modified date
    },
  },
}

const Doc = defineDocumentType(() => ({
  name: 'Doc',
  filePathPattern: '!**/guides/**/*.mdx',
  fields: baseFields,
  computedFields,
}))

const Featured = defineNestedType(() => ({
  name: 'Featured',
  fields: {
    image: {
      type: 'string',
      description:
        'The featured image of the post, not the same as og image. Use 1024x1024 with transparent background.',
      required: true,
    },
    image_alt: {
      type: 'string',
      description: 'The featured image alt of the post',
      required: true,
    },
  },
}))

const Guide = defineDocumentType(() => ({
  name: 'Guide',
  filePathPattern: '**/guides/**/*.mdx',
  fields: {
    ...baseFields,
    published_at: {
      type: 'date',
      description: 'The date the guide was published',
      required: true,
    },
    updated_at: {
      type: 'date',
      description:
        'The date the guide was last updated, will be set to published_at if not set',
    },
    featured: {
      type: 'nested',
      of: Featured,
    },
    tags: {
      type: 'list',
      of: {
        type: 'string',
      },
      description: 'The tags of the post',
      required: true,
    },
    languages: {
      type: 'list',
      of: {
        type: 'string',
      },
      description: 'The languages of the content',
    },
    start_steps: {
      type: 'markdown',
      description: 'The start steps of the doc',
    },
  },
  computedFields,
}))

export default makeSource({
  contentDirPath,
  documentTypes: [Doc, Guide],
})
