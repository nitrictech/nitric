import { cn } from '@/lib/utils'

const classes: Record<string, string> = {
  javascript: 'group-data-[current-lang="javascript"]/body:block',
  typescript: 'group-data-[current-lang="typescript"]/body:block',
  dart: 'group-data-[current-lang="dart"]/body:block',
  python: 'group-data-[current-lang="python"]/body:block',
  go: 'group-data-[current-lang="go"]/body:block',
}

export function ShowIfLang({
  children,
  lang,
}: {
  children: React.ReactNode
  lang: string
}) {
  if (!classes[lang]) {
    console.error(`No class found for language in ShowIfLang: ${lang}`)
    return null
  }

  return <div className={cn('hidden', classes[lang])}>{children}</div>
}
