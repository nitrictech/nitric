'use server'

import { PrismaClient } from '@prisma/client'

const prisma = new PrismaClient()

export async function sendFeedback(prevState: any, formData: FormData) {
  const answer = formData.get('choice') || ''
  const comment = formData.get('comment') || ''
  const ua = formData.get('ua') || ''
  const url = formData.get('url') || ''

  console.log('sendFeedback', { answer, comment, ua, url })

  // disable on non prod
  if (process.env.NEXT_PUBLIC_VERCEL_ENV !== 'production') {
    return { message: 'Not available on production' }
  }

  // validate url and user agent
  if (!ua || !url.toString().startsWith('/docs')) {
    return { message: 'invalid' }
  }

  // validate answer
  if (!['yes', 'no', 'feedback'].includes(answer?.toString())) {
    return { message: 'invalid' }
  }

  try {
    await prisma.feedback.create({
      data: {
        url: url?.toString(),
        answer: answer?.toString(),
        comment: comment?.toString(),
        label: 'is docs page helpful',
        ua: ua?.toString(),
      },
    })

    return { message: 'Feedback sent!' }
  } catch (error) {
    console.error(error)
    return { message: 'failed to store feedback' }
  }
}
