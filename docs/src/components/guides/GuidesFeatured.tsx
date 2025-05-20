import React from 'react'
import { Heading } from '../ui/heading'
import Link from 'next/link'
import Image from 'next/image'
import { ArrowUpRightIcon } from '@heroicons/react/24/outline'
import { allGuides, Guide } from '@/content'

type RequiredFeaturedGuide = Guide & {
  featured: NonNullable<Guide['featured']>
}

const GuidesFeatured: React.FC = ({ take = 3 }: { take?: number }) => {
  const featuredGuides = allGuides
    .filter((guide): guide is RequiredFeaturedGuide => !!guide.featured)
    .sort((a, b) => {
      return a.published_at > b.published_at ? -1 : 1
    })
    .slice(0, take)

  return featuredGuides.length === 0 ? null : (
    <div>
      <Heading level={2} className="sr-only">
        Featured
      </Heading>
      <div className="mx-auto grid max-w-2xl auto-rows-fr grid-cols-1 gap-4 lg:mx-0 lg:max-w-none lg:grid-cols-3">
        {featuredGuides.map((guide) => (
          <article
            key={guide.slug}
            className="group relative isolate flex flex-col justify-end overflow-hidden rounded-lg bg-zinc-900 px-8 pb-8 pt-60"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 1162 896"
              fill="none"
              className="absolute inset-0 top-1/4 -z-10 blur-2xl"
            >
              <g className="opacity-30">
                <g filter="url(#filter0_f_684_803)">
                  <path
                    d="M649.084 233.018C704.89 329.678 671.772 453.277 575.112 509.084C478.452 564.89 354.854 531.772 299.047 435.112C243.24 338.452 276.358 214.854 373.018 159.047C469.678 103.24 593.277 136.358 649.084 233.018Z"
                    className="fill-secondary-400"
                  />
                </g>
                <g filter="url(#filter1_f_684_803)">
                  <path
                    d="M863.21 322.759C918.715 418.896 885.776 541.826 789.639 597.33C693.501 652.835 570.571 619.896 515.067 523.759C459.562 427.622 492.501 304.692 588.638 249.187C684.775 193.682 807.705 226.621 863.21 322.759Z"
                    fill="#2C40F7"
                    className="fill-primary-400"
                  />
                </g>
              </g>
            </svg>
            <Image
              alt={guide.featured.image_alt}
              src={guide.featured.image}
              priority
              fill
              sizes="(max-width: 768px) 100vw, (max-width: 1024px) 50vw, 33vw"
              className="absolute inset-0 -z-10 h-full object-contain"
            />
            <div className="absolute inset-0 -z-10 bg-gradient-to-t from-primary-400/60 via-secondary-400/30 to-primary-500/40 dark:from-primary-800/60 dark:via-secondary-800/20 dark:to-primary-900/50" />
            <div className="absolute inset-0 -z-10 rounded-2xl ring-1 ring-inset ring-primary-500/10 dark:ring-primary-900/10" />
            <time
              dateTime={guide.published_at}
              className="absolute left-0 top-0 m-4 text-2xs text-zinc-300 dark:text-muted-foreground"
            >
              {new Date(guide.published_at).toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
              })}
            </time>
            <h3 className="mt-4 text-lg/6 font-semibold tracking-wide text-white">
              <Link href={`/${guide.slug}`}>
                <span className="absolute inset-0" />
                {guide.title}
              </Link>
            </h3>
            <p className="mt-1 text-sm leading-5 text-white dark:text-muted-foreground">
              {guide.description}
            </p>
            <ArrowUpRightIcon
              aria-hidden="true"
              className="pointer-events-none absolute right-0 top-0 m-4 size-6 text-white transition-transform group-hover:-translate-y-1 group-hover:translate-x-1 dark:text-primary-light/70"
            />
          </article>
        ))}
      </div>
    </div>
  )
}

export default GuidesFeatured
