import { Footer } from '@/components/Footer'
import { Header } from '@/components/layout/Header'
import { Navigation } from '@/components/nav/Navigation'
import { getStarGazers } from '@/lib/stargazers'

export async function BaseLayout({ children }: React.PropsWithChildren) {
  const defaultStarCount = await getStarGazers()

  return (
    <div className="w-full lg:ml-72 xl:ml-80">
      <header className="contents lg:pointer-events-none lg:fixed lg:inset-0 lg:z-40 lg:flex">
        <div className="contents lg:pointer-events-auto lg:mt-[calc(var(--header-height))] lg:block lg:w-72 lg:overflow-y-auto lg:border-r lg:border-zinc-900/10 lg:px-6 lg:pb-8 lg:pt-4 lg:dark:border-white/10 xl:w-80">
          <Header defaultStarCount={defaultStarCount} />

          <Navigation className="hidden lg:block" />
        </div>
      </header>
      <div className="relative px-2 pt-14 sm:px-6 lg:px-8">
        {children}
        <Footer />
      </div>
    </div>
  )
}
