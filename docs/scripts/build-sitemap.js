const fs = require('fs/promises')

const excludedDirectories = ['[[...slug]]'] // ignore dynamic routes
const pages = []

const readDirRecursive = async (dir) => {
  const files = await fs.readdir(dir)

  for (const file of files) {
    const filePath = `${dir}/${file}`
    const stats = await fs.stat(filePath)
    if (stats.isDirectory() && !excludedDirectories.includes(file)) {
      await readDirRecursive(filePath)
    } else if (file.startsWith('page.')) {
      const loc = filePath
        .replace('src/app', '')
        .replace('.tsx', '')
        .replace('.mdx', '')
        .replace('page', '')
        .replace(/\/$/, '')

      pages.push(loc)
    }
  }
}

readDirRecursive('src/app').then(async () => {
  console.log('static pages for sitemap: ', pages)

  await fs.writeFile('./src/assets/sitemap.json', JSON.stringify(pages))
})
