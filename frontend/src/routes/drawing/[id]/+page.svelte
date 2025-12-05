<script lang="ts">
	import { onMount } from 'svelte'
	import { goto } from '$app/navigation'
	import { page } from '$app/state'
	import ExcalidrawWrapper from '$lib/components/ExcalidrawWrapper.svelte'
	import { loadDrawingById } from '$lib/stores/drawing'
	import { drawingsStore } from '$lib/stores/mockDrawings'
	import type { DrawingState } from '$lib/stores/drawing'
	import type { Drawing } from '$lib/stores/mockDrawings'
	import type { ID } from '$lib/types'

	const drawingId = Number(page.params.id) as ID
	let initialData = $state<DrawingState | null>(null)
	let status = $state('Loading...')

	// State for inline name editing
	let isEditingName = $state(false)
	let editingName = $state('')
	let currentDrawing = $state<Drawing | undefined>(undefined)

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

		// Store current drawing reference
		currentDrawing = drawing

		// Load drawing data
		try {
			initialData = loadDrawingById(drawingId)
			status = 'Ready'
		} catch (error) {
			console.error('Failed to load drawing:', error)
			status = 'Error loading drawing'
		}
	})

	// Inline editing functions
	function startEditingName() {
		if (!currentDrawing) return
		isEditingName = true
		editingName = currentDrawing.name
	}

	function saveNameEdit() {
		if (!currentDrawing) return

		const success = drawingsStore.updateDrawingName(drawingId, editingName)

		if (success) {
			isEditingName = false
			// Update local reference
			currentDrawing = drawingsStore.getDrawing(drawingId)
		}
		// If validation fails (empty name), keep editing mode active
	}

	function cancelNameEdit() {
		isEditingName = false
		editingName = ''
	}

	function handleNameKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault()
			saveNameEdit()
		} else if (event.key === 'Escape') {
			cancelNameEdit()
		}
	}
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

		<!-- Center: Drawing Name -->
		<div class="flex-1 flex justify-center px-4">
			{#if isEditingName}
				<input
					type="text"
					class="input input-bordered input-sm w-full max-w-md text-center"
					bind:value={editingName}
					onblur={saveNameEdit}
					onkeydown={handleNameKeydown}
				/>
			{:else if currentDrawing}
				<button
					class="text-lg font-semibold hover:text-primary hover:underline cursor-pointer transition-colors px-2"
					onclick={startEditingName}
					title="Click to edit name"
				>
					{currentDrawing.name}
				</button>
			{/if}
		</div>

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
