/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'wira-primary': '#8B4513',
        'wira-secondary': '#DAA520',
        'wira-accent': '#CD853F',
        'wira-background': '#1a1a1a',
        'wira-text': '#ffffff',
      },
      fontFamily: {
        'cinzel': ['Cinzel', 'serif'],
      },
    },
  },
  plugins: [],
}
