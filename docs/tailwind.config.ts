import typographyPlugin from '@tailwindcss/typography'
import { type Config } from 'tailwindcss'

const {
  default: flattenColorPalette,
} = require('tailwindcss/lib/util/flattenColorPalette')

import defaultTheme from 'tailwindcss/defaultTheme'

import typographyStyles from './typography'

const nitricColors = {
  blue: {
    DEFAULT: '#2C40F7',
    50: '#F3F4FF',
    100: '#DDE0FE',
    200: '#B1B8FC',
    300: '#8490FA',
    400: '#5868F9',
    500: '#2C40F7',
    600: '#0718B6',
    700: '#0718B6',
    800: '#030C59',
    900: '#010319',
  },
  purple: {
    DEFAULT: '#9E2CF7',
    50: '#F9F3FF',
    100: '#EFDDFE',
    200: '#DBB1FC',
    300: '#C27AFA',
    400: '#B258F9',
    500: '#9E2CF7',
    600: '#6907B6',
    700: '#580699',
    800: '#330359',
    900: '#0E0119',
  },
}

export default {
  content: ['./src/**/*.{js,mjs,jsx,ts,tsx,md,mdx}', './docs/**/*.{md,mdx}'],
  darkMode: ['class'],
  theme: {
    fontSize: {
      '2xs': ['0.75rem', { lineHeight: '1.25rem' }],
      xs: ['0.8125rem', { lineHeight: '1.5rem' }],
      sm: ['0.875rem', { lineHeight: '1.5rem' }],
      base: ['1rem', { lineHeight: '1.75rem' }],
      lg: ['1.125rem', { lineHeight: '1.75rem' }],
      xl: ['1.25rem', { lineHeight: '1.75rem' }],
      '2xl': ['1.5rem', { lineHeight: '2rem' }],
      '3xl': ['1.875rem', { lineHeight: '2.25rem' }],
      '4xl': ['2.25rem', { lineHeight: '2.5rem' }],
      '5xl': ['3rem', { lineHeight: '1' }],
      '6xl': ['3.75rem', { lineHeight: '1' }],
      '7xl': ['4.5rem', { lineHeight: '1' }],
      '8xl': ['6rem', { lineHeight: '1' }],
      '9xl': ['8rem', { lineHeight: '1' }],
    },
    typography: typographyStyles,
    extend: {
      fontFamily: {
        sans: ['var(--font-inter)'],
        display: ['var(--font-sora)'],
        mono: ['var(--font-jetbrains-mono)', ...defaultTheme.fontFamily.mono],
      },
      boxShadow: {
        glow: '0 0 4px rgb(0 0 0 / 0.1)',
      },
      maxWidth: {
        lg: '33rem',
        '2xl': '40rem',
        '3xl': '50rem',
        '5xl': '66rem',
      },
      opacity: {
        '1': '0.01',
        '15': '0.15',
        '2.5': '0.025',
        '7.5': '0.075',
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      colors: {
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        'foreground-light': 'var(--foreground-light)',
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        code: {
          DEFAULT: 'hsl(var(--code))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        primary: {
          ...nitricColors.blue,
          light: nitricColors.blue[300],
          DEFAULT: nitricColors.blue[500],
          dark: nitricColors.blue[600],
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          ...nitricColors.purple,
          light: nitricColors.purple[300],
          DEFAULT: nitricColors.purple[500],
          dark: nitricColors.purple[700],
          foreground: 'hsl(var(--secondary-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        chart: {
          '1': 'hsl(var(--chart-1))',
          '2': 'hsl(var(--chart-2))',
          '3': 'hsl(var(--chart-3))',
          '4': 'hsl(var(--chart-4))',
          '5': 'hsl(var(--chart-5))',
        },
        blue: {
          ...nitricColors.blue,
        },
        purple: {
          ...nitricColors.purple,
        },
      },
      animation: {
        orbit: 'orbit calc(var(--duration)*1s) linear infinite',
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'nitric-float': 'nitric-float 10s ease-in-out infinite alternate',
        'nitric-float-2': 'nitric-float-2 10s ease-in-out infinite alternate',
      },
      keyframes: {
        orbit: {
          '0%': {
            transform:
              'rotate(0deg) translateY(calc(var(--radius) * 1px)) rotate(0deg)',
          },
          '100%': {
            transform:
              'rotate(360deg) translateY(calc(var(--radius) * 1px)) rotate(-360deg)',
          },
        },
        'accordion-down': {
          from: {
            height: '0',
          },
          to: {
            height: 'var(--radix-accordion-content-height)',
          },
        },
        'accordion-up': {
          from: {
            height: 'var(--radix-accordion-content-height)',
          },
          to: {
            height: '0',
          },
        },
        'nitric-float': {
          '0%': { transform: 'translate(-2%, -2%)' },
          '20%': { transform: 'translate(20%, 22%)' },
          '40%': { transform: 'translate(25%, 28%)' },
          '60%': { transform: 'translate(20%, 20%)' },
          '80%': { transform: 'translate(12%, 10%)' },
          '100%': { transform: 'translate(-2%, -2%)' },
        },
        'nitric-float-2': {
          '0%': { transform: 'translate(2%, 2%)' },
          '20%': { transform: 'translate(-20%, -22%)' },
          '40%': { transform: 'translate(-25%, -28%)' },
          '60%': { transform: 'translate(-20%, -20%)' },
          '80%': { transform: 'translate(-12%, -10%)' },
          '100%': { transform: 'translate(2%, 2%)' },
        },
      },
    },
  },
  plugins: [
    typographyPlugin,
    require('tailwindcss-animate'),
    addVariablesForColors,
  ],
} satisfies Config

// This plugin adds each Tailwind color as a global CSS variable, e.g. var(--gray-200).
function addVariablesForColors({ addBase, theme }: any) {
  let allColors = flattenColorPalette(theme('colors'))
  let newVars = Object.fromEntries(
    Object.entries(allColors).map(([key, val]) => [`--${key}`, val]),
  )

  addBase({
    ':root': newVars,
  })
}
