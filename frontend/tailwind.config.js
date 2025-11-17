/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#FFF9F5',
          100: '#FFF0E6',
          200: '#FFE6D5',
          300: '#FFD4B8',
          400: '#FFC299',
          500: '#FFB08A',
          600: '#FF9466',
          700: '#FF7B47',
          800: '#FF6B35',
          900: '#E65A28',
        },
        orange: {
          DEFAULT: '#FFC299',
          light: '#FFD4B8',
          dark: '#FFB08A',
        }
      },
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', 'sans-serif'],
      },
    },
  },
  plugins: [],
}

