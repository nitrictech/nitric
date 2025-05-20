import React from 'react'

export const HomeHeader = ({
  title,
  description,
}: {
  title: string
  description: string
}) => {
  return (
    <>
      <h1 className="mb-3 text-2xl sm:text-3xl">{title}</h1>
      <p className="max-w-[700px] text-lg">{description}</p>
    </>
  )
}
