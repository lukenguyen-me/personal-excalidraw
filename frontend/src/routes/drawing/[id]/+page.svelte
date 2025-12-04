<script lang="ts">
	import ExcalidrawWrapper from '$lib/components/ExcalidrawWrapper.svelte'
	import { page } from '$app/stores'

	const drawingId = $page.params.id

	// Placeholder for status - will be implemented later
	let status = 'Ready'
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
		<ExcalidrawWrapper />
	</div>
</div>
