<script lang="ts">
	import { onMount } from 'svelte'
	import { goto } from '$app/navigation'
	import { page } from '$app/state'
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query'
	import { drawingsAPI } from '$lib/api'
	import ExcalidrawWrapper from '$lib/components/ExcalidrawWrapper.svelte'
	import type { DrawingState } from '$lib/stores/drawing'
	import type { ID } from '$lib/types'

	const queryClient = useQueryClient()
	const drawingId = page.params.id as ID

	// Query for drawing metadata
	const drawingQuery = createQuery(() => ({
		queryKey: ['drawing', drawingId],
		queryFn: () => drawingsAPI.get(drawingId),
		enabled: !!drawingId
	}))

	// Mutation for updating name
	const updateNameMutation = createMutation(() => ({
		mutationFn: (name: string) => drawingsAPI.update(drawingId, { name }),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['drawing', drawingId] })
			queryClient.invalidateQueries({ queryKey: ['drawings'] })
		}
	}))

	let initialData = $state<DrawingState | null>(null)

	// UI state for name editing (Svelte store)
	let isEditingName = $state(false)
	let editingName = $state('')

	onMount(async () => {
		if (!drawingId) {
			goto('/')
			return
		}

		// Load drawing content from cloud
		try {
			const drawing = await drawingsAPI.get(drawingId)
			initialData = {
				elements: (drawing.data?.elements || []) as unknown[],
				appState: (drawing.data?.appState || {}) as unknown,
				files: (drawing.data?.files || {}) as Record<string, unknown>
			}
		} catch (error) {
			console.error('Failed to load drawing from cloud:', error)
			// Initialize with empty data if loading fails
			initialData = {
				elements: [],
				appState: {},
				files: {}
			}
		}
	})

	function startEditingName() {
		if (!drawingQuery.data) return
		isEditingName = true
		editingName = drawingQuery.data.name
	}

	function saveNameEdit() {
		if (!editingName.trim()) return
		updateNameMutation.mutate(editingName.trim())
		isEditingName = false
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

		<div class="flex-1 flex justify-center px-4">
			{#if drawingQuery.isPending}
				<span class="loading loading-spinner loading-sm"></span>
			{:else if isEditingName}
				<input
					type="text"
					class="input input-bordered input-sm w-full max-w-md text-center"
					bind:value={editingName}
					onblur={saveNameEdit}
					onkeydown={handleNameKeydown}
				/>
			{:else if drawingQuery.data}
				<button
					class="text-lg font-semibold hover:text-primary hover:underline cursor-pointer transition-colors px-2"
					onclick={startEditingName}
					title="Click to edit name"
				>
					{drawingQuery.data.name}
				</button>
			{/if}
		</div>

		<div class="flex items-center gap-2">
			<span class="text-sm text-base-content/70">
				{drawingQuery.isPending ? 'Loading...' : 'Ready'}
			</span>
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
