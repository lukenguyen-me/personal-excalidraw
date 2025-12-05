<script lang="ts">
	import { goto } from '$app/navigation'
	import { drawingsStore } from '$lib/stores/mockDrawings'
	import type { ID } from '$lib/types'

	// State for inline name editing
	let editingId = $state<ID | null>(null)
	let editingValue = $state('')

	function formatDate(date: Date): string {
		return new Date(date).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		})
	}

	function handleNewDrawing() {
		// Generate new ID
		const newId = Date.now() as ID

		// Create drawing entry
		drawingsStore.addDrawing({
			id: newId,
			name: 'Untitled Drawing',
			createdAt: new Date(),
			updatedAt: new Date()
		})

		// Navigate to new drawing
		goto(`/drawing/${newId}`)
	}

	function handleDelete(id: ID) {
		drawingsStore.deleteDrawing(id)
	}

	// Inline editing functions
	function startEditing(id: ID, currentName: string) {
		editingId = id
		editingValue = currentName
	}

	function saveEdit() {
		if (editingId === null) return

		const success = drawingsStore.updateDrawingName(editingId, editingValue)

		if (success) {
			editingId = null
			editingValue = ''
		}
		// If validation fails (empty name), keep editing mode active
	}

	function cancelEdit() {
		editingId = null
		editingValue = ''
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault()
			saveEdit()
		} else if (event.key === 'Escape') {
			cancelEdit()
		}
	}
</script>

<div class="min-h-screen bg-base-100 p-8">
	<div class="max-w-6xl mx-auto">
		<div class="flex justify-between items-center mb-8">
			<h1 class="text-3xl font-bold">My Drawings</h1>
			<button onclick={handleNewDrawing} class="btn btn-primary">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				New
			</button>
		</div>

		<div class="overflow-x-auto">
			<table class="table table-zebra w-full">
				<thead>
					<tr>
						<th>Name</th>
						<th>Created At</th>
						<th>Updated At</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each $drawingsStore as drawing}
						<tr>
							<td class="font-medium">
								{#if editingId === drawing.id}
									<input
										type="text"
										class="input input-bordered input-sm w-full max-w-xs"
										bind:value={editingValue}
										onblur={saveEdit}
										onkeydown={handleKeydown}
									/>
								{:else}
									<button
										class="text-left hover:text-primary hover:underline cursor-pointer transition-colors"
										onclick={() => startEditing(drawing.id, drawing.name)}
										title="Click to edit name"
									>
										{drawing.name}
									</button>
								{/if}
							</td>
							<td>{formatDate(drawing.createdAt)}</td>
							<td>{formatDate(drawing.updatedAt)}</td>
							<td class="text-right">
								<div class="flex gap-2 justify-end">
									<a
										href="/drawing/{drawing.id}"
										class="btn btn-ghost btn-info btn-sm"
										aria-label="Edit drawing"
									>
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
												d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
											/>
										</svg>
									</a>
									<button
										class="btn btn-ghost btn-error btn-sm"
										onclick={() => handleDelete(drawing.id)}
										aria-label="Delete drawing"
									>
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
												d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
											/>
										</svg>
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
