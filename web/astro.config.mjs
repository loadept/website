// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';

import preact from '@astrojs/preact';

import sitemap from '@astrojs/sitemap';

// https://astro.build/config
export default defineConfig({
  outDir: '../public',
  site: 'https://loadept.com',
  vite: {
    plugins: [tailwindcss()]
  },
  integrations: [
    preact({ compat: true }),
    sitemap()
  ],
});