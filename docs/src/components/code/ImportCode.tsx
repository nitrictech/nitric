import React from 'react'
import { Code } from './Code'
import type { RawCode } from 'codehike/code'
import path from 'path'
import fs from 'fs/promises'

export const ImportCode: React.FC<{ file: string; meta: string }> = async ({
  file,
  meta = '',
}) => {
  try {
    const contents = await fs.readFile(path.join(process.cwd(), file), 'utf-8')

    const lang = path.extname(file).replace('.', '')

    const codeBlock: RawCode = {
      value: contents.trim(),
      lang,
      meta,
    }

    return <Code codeblock={codeBlock} />
  } catch (error) {
    console.error('Error reading file for ImportCode:', error)
    return <div>Error reading file {file} for ImportCode</div>
  }
}
