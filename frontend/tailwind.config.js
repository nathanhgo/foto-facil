/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#121212',
        surface: '#1E1E1E',
        primary: '#9b51e0',
        primaryHover: '#b272ec',
        text: '#E0E0E0'
      }
    },
  },
  plugins: [],
}
