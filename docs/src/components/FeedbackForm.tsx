'use client'

import { useState } from 'react'
import { Button } from './ui/button'
import { sendFeedback } from '@/actions/sendFeedback'
import { RadioGroup } from './ui/radio-group'
import { Label } from './ui/label'
import { RadioGroupItem } from '@radix-ui/react-radio-group'
import { useFormState } from 'react-dom'
import { CheckIcon } from './icons/CheckIcon'
import { usePathname } from 'next/navigation'

const choices = [
  {
    value: 'yes',
    label: 'It was helpful',
    emoji: '🤩',
  },
  {
    value: 'no',
    label: 'It was not helpful',
    emoji: '😓',
  },
  {
    value: 'feedback',
    label: 'I have feedback',
    emoji: '📣',
  },
]

const initialState = {
  message: '',
}

const FeedbackForm = () => {
  const [state, formAction] = useFormState(sendFeedback, initialState)
  const [selected, setSelected] = useState(false)
  const pathname = usePathname()

  return state.message ? (
    <div className="flex items-center gap-x-2 text-sm">
      <CheckIcon className="h-5 w-5 flex-none fill-green-500 stroke-white dark:fill-green-200/20 dark:stroke-green-200" />
      Thank you for your feedback! 🙌
    </div>
  ) : (
    <form
      action={formAction}
      className="flex flex-col items-start justify-start gap-6"
    >
      <input
        name="ua"
        value={typeof window !== 'undefined' ? window.navigator.userAgent : ''}
        className="hidden"
        readOnly
      />
      <input
        name="url"
        value={`/docs${pathname}`}
        className="hidden"
        readOnly
      />
      <Label htmlFor="choice" className="text-sm text-zinc-900 dark:text-white">
        What did you think of this content?
      </Label>
      <RadioGroup
        className="flex flex-col gap-2 md:flex-row"
        onChange={() => {
          if (!selected) setSelected(true)
        }}
        name="choice"
        id="choice"
      >
        {choices.map(({ label, value, emoji }) => (
          <div key={value} className="group flex items-center">
            <RadioGroupItem value={value} id={value} className="group flex">
              <span className="mr-2 grayscale transition-all group-checked:grayscale-0 group-hover:grayscale-0 group-data-[state=checked]:grayscale-0">
                {emoji}
              </span>
              <Label
                htmlFor={value}
                className="cursor-pointer text-xs text-muted-foreground transition-colors group-hover:text-zinc-900 group-data-[state=checked]:text-zinc-900 dark:group-hover:text-white dark:group-data-[state=checked]:text-white"
              >
                {label}
              </Label>
            </RadioGroupItem>
          </div>
        ))}
      </RadioGroup>
      {selected && (
        <div className="flex w-full max-w-[400px] flex-col gap-4 rounded-lg p-2 text-zinc-700 shadow-md ring-1 ring-zinc-300 dark:bg-white/2.5 dark:text-zinc-100 dark:ring-white/10 dark:hover:shadow-black/5">
          <label htmlFor="comment" className="sr-only">
            Comment
          </label>
          <div>
            <textarea
              rows={5}
              name="comment"
              id="comment"
              className="block w-full rounded-md border-0 bg-white p-2 py-1.5 text-zinc-900 shadow-sm ring-1 ring-inset ring-zinc-300 placeholder:text-zinc-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 dark:bg-zinc-700/70 dark:text-zinc-50 dark:ring-zinc-600 sm:text-sm sm:leading-6"
              placeholder="We'd love to hear your feedback!"
              defaultValue={''}
              autoFocus
            />
          </div>
          <div className="mt-2 flex justify-end">
            <Button type="submit">Send</Button>
          </div>
        </div>
      )}
    </form>
  )
}

export default FeedbackForm
