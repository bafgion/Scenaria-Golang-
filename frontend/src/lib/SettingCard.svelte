<script lang="ts">
  export let title = ''
  export let description = ''
</script>

<div class="setting-card">
  <div class="setting-card-text">
    <div class="setting-card-title">{title}</div>
    {#if description}
      <div class="setting-card-desc">{description}</div>
    {/if}
  </div>
  <div class="setting-card-control">
    <slot />
  </div>
</div>
