<script lang="ts">
	import { onMount } from 'svelte'
	import { goto } from '$app/navigation'
	import { page } from '$app/state'
	import ExcalidrawWrapper from '$lib/components/ExcalidrawWrapper.svelte'
	import { loadDrawingById } from '$lib/stores/drawing'
	import { drawingsStore } from '$lib/stores/mockDrawings'
	import type { DrawingState } from '$lib/stores/drawing'
	import type { ID } from '$lib/types'

	const drawingId = Number(page.params.id) as ID
	let initialData: DrawingState | null = null
	let status = 'Loading...'

	onMount(() => {
		if (!drawingId || isNaN(drawingId)) {
			status = 'Invalid drawing ID'
			goto('/')
			return
		}

		// Verify drawing exists in metadata
		const drawing = drawingsStore.getDrawing(drawingId)
		if (!drawing) {
			status = 'Drawing not found'
			goto('/')
			return
		}

		// Load drawing data
		try {
			initialData = loadDrawingById(drawingId)
			status = 'Ready'
		} catch (error) {
			console.error('Failed to load drawing:', error)
			status = 'Error loading drawing'
		}
	})
</script>

<div class="h-screen w-screen flex flex-col bg-base-100">
	<!-- Header -->
	<div class="flex items-center justify-between px-4 py-2 border-b border-base-300">
		<!-- Left: Back button -->
		<a href="/" class="btn btn-sm btn-ghost">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="size-4"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M10 19l-7-7m0 0l7-7m-7 7h18"
				/>
			</svg>
			Back
		</a>

		<!-- Right: Status indicator -->
		<div class="flex items-center gap-2">
			<span class="text-sm text-base-content/70">{status}</span>
		</div>
	</div>

	<!-- Drawing canvas -->
	<div class="flex-1 flex overflow-hidden">
		{#if initialData && drawingId}
			<ExcalidrawWrapper {drawingId} {initialData} />
		{:else}
			<div class="flex items-center justify-center w-full h-full">
				<span class="loading loading-spinner loading-lg"></span>
			</div>
		{/if}
	</div>
</div>
