import Image from 'next/image'

import logoNode from '@/images/logos/node.svg'
import logoPython from '@/images/logos/python.svg'
import logoCsharp from '@/images/logos/csharp.svg'
import logoGo from '@/images/logos/go.svg'
import logoJava from '@/images/logos/java.svg'
import logoDart from '@/images/logos/dart.svg'
import Link from 'next/link'
import { ArrowRightIcon } from '@heroicons/react/24/outline'
import { Heading } from './ui/heading'

const libraries = [
  {
    href: '/reference/nodejs',
    name: 'Node.js',
    description: 'View full API for Node.js',
    logo: logoNode,
  },
  {
    href: '/reference/python',
    name: 'Python',
    description: 'View full API for Python',
    logo: logoPython,
  },
  {
    href: '/reference/go',
    name: 'Go',
    description: 'View full API for Go',
    logo: logoGo,
  },
  {
    href: '/reference/dart',
    name: 'Dart',
    description: 'View full API for Dart',
    logo: logoDart,
  },
  {
    href: '/reference/csharp/v0',
    name: 'C# .NET',
    description: 'View full API for C# .NET',
    logo: logoCsharp,
  },
  {
    href: '/reference/jvm/v0',
    name: 'JVM',
    description: 'View full API for JVM',
    logo: logoJava,
  },
]

export interface LibrariesProps {
  minimal?: boolean
}

export function Libraries({ minimal = false }: LibrariesProps) {
  if (minimal) {
    return (
      <div className="flex h-fit w-fit items-center gap-2">
        {libraries.map((library) => (
          <Link
            href={library.href}
            key={library.name}
            className="opacity-90 grayscale transition-opacity hover:opacity-100 hover:grayscale-0"
            target="_blank"
          >
            <Image
              src={library.logo}
              alt={library.name + ' Logo'}
              className="h-12 w-12"
              unoptimized
            />
          </Link>
        ))}
      </div>
    )
  }

  return (
    <div className="my-16 xl:max-w-none">
      <Heading level={2} id="libraries">
        Libraries
      </Heading>
      <div className="not-prose mt-4 grid grid-cols-1 gap-x-6 gap-y-10 border-t border-zinc-900/5 pt-10 dark:border-white/5 sm:grid-cols-2 xl:max-w-none xl:grid-cols-3">
        {libraries.map((library) => (
          <Link
            href={library.href}
            key={library.name}
            className="group relative flex items-center gap-4 rounded-2xl bg-zinc-50 p-4 transition-shadow hover:shadow-md hover:shadow-zinc-900/5 dark:bg-white/2.5 dark:hover:shadow-black/5"
          >
            <div className="absolute inset-0 rounded-2xl ring-1 ring-inset ring-zinc-900/7.5 group-hover:ring-zinc-900/10 dark:ring-white/10 dark:group-hover:ring-white/20" />

            <Image
              src={library.logo}
              alt={library.name + ' Logo'}
              className="h-12 w-12"
              unoptimized
            />
            <div className="flex-auto">
              <h3 className="text-lg font-semibold text-zinc-900 dark:text-white">
                {library.name}
              </h3>
              <p className="mt-1 text-sm text-zinc-600 dark:text-zinc-400">
                {library.description}
              </p>
            </div>
            <ArrowRightIcon className="h-5 w-5 -translate-x-1 transition-transform group-hover:translate-x-0" />
          </Link>
        ))}
      </div>
    </div>
  )
}
