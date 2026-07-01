import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  build: {
    chunkSizeWarningLimit: 3500,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules/monaco-editor')) {
            return 'monaco'
          }
          if (id.includes('node_modules/@monaco-editor')) {
            return 'monaco'
          }
          if (id.includes('node_modules')) {
            return 'vendor'
          }
        },
      },
    },
  },
  optimizeDeps: {
    include: ['monaco-editor', '@monaco-editor/loader'],
  },
  worker: {
    format: 'es',
  },
})
