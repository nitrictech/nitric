const DEFAULT_STARGAZERS = 1250

export const getStarGazers = async (): Promise<number> => {
  try {
    // on next < 15.x fetch will force-cache by default
    const response = await fetch(
      'https://api.github.com/repos/nitrictech/nitric',
    )
    if (!response.ok) {
      console.error('Error fetching star count:', response.statusText)

      return DEFAULT_STARGAZERS
    }

    const repoData = await response.json()
    return repoData.stargazers_count
  } catch (e) {
    console.error('Error fetching star count:', e)
    return DEFAULT_STARGAZERS
  }
}
