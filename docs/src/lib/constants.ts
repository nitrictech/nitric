export const LANGUAGE_LABEL_MAP: Record<string, string> = {
  javascript: 'JavaScript',
  python: 'Python',
  typescript: 'TypeScript',
  go: 'Go',
  dart: 'Dart',
  java: 'JVM',
  kotlin: 'Kotlin',
  csharp: 'C#',
}

export const languages = [
  'javascript',
  'typescript',
  'python',
  'go',
  'dart',
] as const

export type Language = (typeof languages)[number]

export const isMobile =
  typeof navigator !== 'undefined' &&
  /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
    navigator.userAgent,
  )

export const discordChatUrl =
  process.env.NEXT_PUBLIC_VERCEL_ENV === 'production'
    ? 'https://nitric.io/chat'
    : process.env.NEXT_PUBLIC_VERCEL_URL
      ? `https://${process.env.NEXT_PUBLIC_VERCEL_URL}/chat`
      : `http://localhost:3000/chat`

export const BASE_URL =
  process.env.NEXT_PUBLIC_VERCEL_ENV === 'production'
    ? 'https://nitric.io'
    : process.env.NEXT_PUBLIC_VERCEL_URL
      ? `https://${process.env.NEXT_PUBLIC_VERCEL_URL}`
      : 'http://localhost:3000'
