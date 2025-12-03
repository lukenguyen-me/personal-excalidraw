<script lang="ts">
	import { onMount } from 'svelte'
	import { drawingStore } from '$lib/stores/drawing'
	import { uiStore } from '$lib/stores/ui'

	let container: HTMLDivElement
	let excalidrawApp: any

	onMount(async () => {
		// Dynamic import to avoid issues with React in Svelte
		try {
			const { Excalidraw } = await import('@excalidraw/excalidraw')
			const React = await import('react')
			const ReactDOM = await import('react-dom/client')

			if (!container) return

			// Create a React root and render Excalidraw
			const root = ReactDOM.createRoot(container)

			const handleChange = (elements: any, appState: any, files: any) => {
				drawingStore.set({
					elements,
					appState,
					files,
				})
			}

			const ExcalidrawComponent = React.createElement(Excalidraw, {
				onChange: handleChange,
				initialData: {
					elements: [],
					appState: {},
					files: {},
				},
			})

			root.render(ExcalidrawComponent)
			excalidrawApp = root
		} catch (error) {
			console.error('Failed to load Excalidraw:', error)
		}
	})
</script>

<div bind:this={container} class="w-full h-full"></div>

<style>
	:global(.excalidraw) {
		font-family: 'Cascadia Code', monospace;
	}
</style>
