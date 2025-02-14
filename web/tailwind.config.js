/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      backgroundImage: {
        'gradient-90': 'linear-gradient(90deg, var(--tw-gradient-stops))',
      },
      fontFamily: {
        lato: ['Lato', 'sans-serif'],
      },
      colors: {
        blue: '#007bff',
        indigo: '#6610f2',
        purple: '#6f42c1',
        pink: '#e83e8c',
        red: '#dc3545',
        orange: '#fd7e14',
        yellow: '#ffc107',
        green: '#28a745',
        teal: '#20c997',
        cyan: '#17a2b8',
        white: '#fff',
        gray: '#6c757d',
        'gray-dark': '#343a40',
        primary: '#003b72',
        secondary: '#005CB2',
        success: '#28a745',
        info: '#17a2b8',
        warning: '#ffc107',
        danger: '#dc3545',
        light: '#ECF0F1',
        dark: '#343a40',
      },
      keyframes: {
        'enter-from-top': {
          '0%': {
            transform: 'translateY(-100%)',
            visibility: 'visible',
            opacity: 0,
          },
          '100%': {
            transform: 'translateY(0%)',
            opacity: 1,
          },
        },
        'exit-to-top': {
          '0%': {
            transform: 'translateY(0%)',
            visibility: 'visible',
            opacity: 1,
          },
          '100%': {
            transform: 'translateY(-100%)',
            visibility: 'hidden',
            opacity: 0,
          },
        },
        'open-with-fade': {
          '0%': {
            opacity: 0,
            visibility: 'hidden',
          },
          '0.01%': {
            visibility: 'visible',
          },
          '100%': {
            opacity: 1,
          },
        },
        'close-with-fade': {
          '0%': {
            opacity: 1,
            visibility: 'visible',
          },
          '100%': {
            opacity: 0,
            visibility: 'hidden',
          },
        },
      },
      animation: {
        'enter-from-top': 'enter-from-top 300ms ease-in-out forwards',
        'exit-to-top': 'exit-to-top 300ms ease-in-out forwards',
        'open-with-fade': 'open-with-fade 200ms ease-in-out forwards',
        'close-with-fade': 'close-with-fade 200ms ease-in-out forwards',
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
