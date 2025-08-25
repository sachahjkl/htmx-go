/**
 * This file only exists so that tailwind vs-code plugin detects it and
 * provides useful autocomplete for classnames
 * cf : https://tailwindcss.com/docs/configuration
 * >By default, Tailwind will look for an optional tailwind.config.js file
 * >at the root of your project where you can define any customizations.
 */
module.exports = {
  content: ["./**/*.{html,js,ts,svelte,jsx,tsx,md,mdx,}"],
  theme: {
    extend: {
      fontFamily: {
        comic: ["Comic Sans MS", "Comic Sans", "cursive"],
      },
      colors: {
        clifford: "#da373d",
      },
    },
  },
};
