<script lang="ts">
	import { onMount } from 'svelte'
	import { drawingStore } from '$lib/stores/drawing'
	import { uiStore } from '$lib/stores/ui'
	import { excalidrawStore } from '$lib/stores/excalidraw'
	import type { ExcalidrawImperativeAPI } from '@excalidraw/excalidraw/types/types'

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

				// Sync active tool with UI store
				if (appState?.activeTool?.type) {
					uiStore.setActiveTool(appState.activeTool.type)
				}
			}

			const handleExcalidrawAPI = (api: ExcalidrawImperativeAPI) => {
				excalidrawStore.setAPI(api)
			}

			const ExcalidrawComponent = React.createElement(Excalidraw, {
				onChange: handleChange,
				excalidrawAPI: handleExcalidrawAPI,
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

	/* Hide the library button */
	:global(.excalidraw .sidebar-trigger) {
		display: none !important;
	}
</style>
