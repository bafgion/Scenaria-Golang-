<script lang="ts">
  import { HighlightFeature } from '../../wailsjs/go/wailsapp/App'
  import { gui } from '../../wailsjs/go/models'

  export let text = ''

  let spans: gui.HighlightSpan[] = []
  let timer: ReturnType<typeof setTimeout> | undefined

  $: scheduleHighlight(text)

  function scheduleHighlight(value: string) {
    clearTimeout(timer)
    timer = setTimeout(() => refreshHighlight(value), 120)
  }

  async function refreshHighlight(value: string) {
    try {
      spans = await HighlightFeature(value)
    } catch {
      spans = [{ text: value, kind: 'text' }]
    }
  }
</script>

<div class="feature-preview">
  <pre class="preview-code"><code>{#each spans as span}<span class={span.kind}>{span.text}</span>{/each}</code></pre>
</div>

<style>
  .feature-preview {
    height: 100%;
    overflow: auto;
    background: #12141a;
    border-left: 1px solid var(--color-border);
  }

  .preview-code {
    margin: 0;
    padding: 8px 12px;
    font-family: var(--font-mono);
    font-size: 13px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .comment {
    color: #6a737d;
    font-style: italic;
  }

  .tag {
    color: #79c0ff;
  }

  .gherkin {
    color: #d2a8ff;
    font-weight: 600;
  }

  .step {
    color: #ff7b72;
    font-weight: 600;
  }

  .block {
    color: #ffa657;
    font-weight: 600;
  }

  .string {
    color: #a5d6ff;
  }

  .testclient {
    color: #ffa657;
    font-weight: 600;
  }

  .error {
    color: var(--color-error);
    text-decoration: underline wavy var(--color-error);
  }

  .text {
    color: var(--color-text);
  }
</style>
