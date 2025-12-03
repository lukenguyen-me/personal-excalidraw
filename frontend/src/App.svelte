<script lang="ts">
  import './app.css'
  import ExcalidrawWrapper from '$lib/components/ExcalidrawWrapper.svelte'
</script>

<svelte:window />

<div class="h-screen w-screen flex flex-col bg-base-100">
  <header class="bg-base-200 border-b border-base-300 p-4">
    <h1 class="text-2xl font-bold text-base-content">Personal Excalidraw</h1>
  </header>

  <div class="flex-1 flex overflow-hidden">
    <ExcalidrawWrapper />
  </div>
</div>

<style global lang="postcss">
  :global(html, body) {
    @apply m-0 p-0 w-full h-full overflow-hidden;
  }

  :global(body) {
    @apply bg-base-100 text-base-content;
  }
</style>
