/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        taxi: '#FFD100', // Taksi sarısı
        dark: '#202020', // Koyu tema rengi
      }
    },
  },
  plugins: [],
}