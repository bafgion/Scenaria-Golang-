<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { createTranslator, locale } from './i18n'
  import { SearchSteps } from '../../wailsjs/go/wailsapp/App'
  import { asStepSearchQuery } from './stepSearch'
  import { debounce } from './uiScheduler'
  import type { StepCatalogEntry } from './stepTypes'

  export let onInsert: (template: string) => void = () => {}
  export let onClose: () => void = () => {}

  $: tr = createTranslator($locale)

  let query = ''
  let entries: StepCatalogEntry[] = []
  let loading = true

  async function runSearch(value: string) {
    loading = true
    try {
      entries = await SearchSteps(asStepSearchQuery(value))
    } catch {
      entries = []
    } finally {
      loading = false
    }
  }

  const debouncedSearch = debounce((value: string) => {
    void runSearch(value)
  }, 200)

  onMount(async () => {
    await runSearch('')
  })

  onDestroy(() => {
    debouncedSearch.cancel()
  })

  function onQueryInput() {
    debouncedSearch(query)
  }

  function pick(template: string) {
    onInsert(template)
    onClose()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.stopImmediatePropagation()
      onClose()
    }
  }
</script>