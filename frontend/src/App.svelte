<script lang="ts">
  import ExcalidrawWrapper from '$lib/components/ExcalidrawWrapper.svelte'
</script>

<div class="h-screen w-screen flex flex-col bg-base-100">
  <header class="bg-base-200 border-b border-base-300 p-4">
    <h1 class="text-2xl font-bold text-base-content">Personal Excalidraw</h1>
  </header>

  <div class="flex-1 flex overflow-hidden">
    <ExcalidrawWrapper />
  </div>
</div>

<style global lang="postcss">
  html, body {
    margin: 0;
    padding: 0;
    height: 100%;
    width: 100%;
    overflow: hidden;
  }
</style>
