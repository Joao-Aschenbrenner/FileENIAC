/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        eniac: {
          50:  "#F0F4F8",
          100: "#E0F5FA",
          200: "#C4ECF5",
          300: "#89D8EC",
          400: "#22D3EE",
          500: "#06B6D4",
          600: "#0891B2",
          700: "#0E7490",
          800: "#155E75",
          900: "#162032",
          950: "#0B1220",
        },
      },
    },
  },
  plugins: [],
};
