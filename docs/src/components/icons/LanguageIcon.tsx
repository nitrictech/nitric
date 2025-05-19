import React from 'react'
import JavaScriptLogoColour from './JavaScriptLogoColour'
import TypeScriptLogoColour from './TypeScriptLogoColour'
import PythonColorLogo from './PythonLogoColour'
import GoColorLogo from './GoLogoColour'
import DartLogoNoTextColour from './DartLogoNoTextColour'
import { cn } from '@/lib/utils'
import type { Language } from '@/lib/constants'

interface LanguageIconProps {
  name: Language
  className?: string
}

const icons: Record<Language, React.FC<{ className: string }>> = {
  javascript: JavaScriptLogoColour,
  typescript: TypeScriptLogoColour,
  python: PythonColorLogo,
  go: GoColorLogo,
  dart: DartLogoNoTextColour,
}

export const LanguageIcon: React.FC<LanguageIconProps> = ({
  name,
  className,
}) => {
  const IconComponent = icons[name]

  return <IconComponent className={cn(`size-8`, className)} />
}
