import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  build: {
    chunkSizeWarningLimit: 3500,
  },
  optimizeDeps: {
    include: ['monaco-editor', '@monaco-editor/loader'],
  },
  worker: {
    format: 'es',
  },
})
