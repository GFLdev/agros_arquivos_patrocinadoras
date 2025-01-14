/** @type {import('tailwindcss').Config} */
export default {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      backgroundImage: {
        'gradient-90': 'linear-gradient(90deg, var(--tw-gradient-stops))',
      },
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
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
