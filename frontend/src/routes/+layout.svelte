<script lang="ts">
	import './layout.css'
	import ExcalidrawWrapper from '$lib/components/ExcalidrawWrapper.svelte'
	import Toolbar from '$lib/components/Toolbar.svelte'
	import favicon from '$lib/assets/favicon.svg'

	let { children } = $props()
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="h-screen w-screen flex flex-col bg-base-100" data-theme="light">
	<header class="bg-base-200 border-b border-base-300 p-4">
		<h1 class="text-2xl font-bold text-base-content">Personal Excalidraw</h1>
	</header>

	<Toolbar />

	<div class="flex-1 flex overflow-hidden">
		<ExcalidrawWrapper />
	</div>
</div>
