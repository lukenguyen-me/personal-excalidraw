<script lang="ts">
  import './app.css';
</script>

<svelte:window />

<main>
  <slot />
</main>

<style global lang="postcss">
  :global(html, body) {
    @apply m-0 p-0 w-full h-full;
  }

  :global(body) {
    @apply bg-base-100 text-base-content;
  }
</style>
