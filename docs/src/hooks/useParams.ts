import { usePathname, useSearchParams } from 'next/navigation'
import { useCallback } from 'react'

const useParams = () => {
  const searchParams = useSearchParams()
  const pathname = usePathname()

  const setParams = useCallback(
    (name: string, value: string | null) => {
      // Apparently this nonsense is necessary to update the URL.
      //  See: https://github.com/vercel/next.js/discussions/47583
      const currentParams = new URLSearchParams(
        Array.from(searchParams.entries()),
      )

      if (!value) {
        currentParams.delete(name)
      } else {
        currentParams.set(name, value)
      }

      const search = currentParams.toString()
      const url = search ? `/docs${pathname}?${search}` : `/docs${pathname}`

      window.history.pushState(null, '', url)
    },
    [searchParams, pathname],
  )

  return {
    setParams,
    searchParams,
  }
}

export default useParams
