<script lang="ts">
	import { onMount, onDestroy } from 'svelte'
	import { get } from 'svelte/store'
	import { drawingStore, saveDrawingById } from '$lib/stores/drawing'
	import { uiStore } from '$lib/stores/ui'
	import { excalidrawStore } from '$lib/stores/excalidraw'
	import type { ExcalidrawImperativeAPI } from '@excalidraw/excalidraw/types/types'
	import type { DrawingState } from '$lib/stores/drawing'
	import type { ID } from '$lib/types'

	// Accept props from parent
	export let drawingId: ID
	export let initialData: DrawingState

	let container: HTMLDivElement
	let excalidrawApp: any
	let saveTimeout: ReturnType<typeof setTimeout> | null = null
	const AUTOSAVE_DELAY = 1000 // 1 second debounce

	// Debounced save function
	function debouncedSave(state: DrawingState) {
		if (saveTimeout) {
			clearTimeout(saveTimeout)
		}

		saveTimeout = setTimeout(() => {
			const success = saveDrawingById(drawingId, state)
			if (success) {
				console.log('Drawing auto-saved:', drawingId)
			} else {
				console.error('Failed to save drawing')
			}
			saveTimeout = null
		}, AUTOSAVE_DELAY)
	}

	onMount(() => {
		// Dynamic import to avoid issues with React in Svelte
		const initExcalidraw = async () => {
			try {
				const { Excalidraw } = await import('@excalidraw/excalidraw')
				const React = await import('react')
				const ReactDOM = await import('react-dom/client')

				if (!container) return

				// Create a React root and render Excalidraw
				const root = ReactDOM.createRoot(container)

				const handleChange = (elements: any, appState: any, files: any) => {
					const state = { elements, appState, files }

					// Update store immediately for reactivity
					drawingStore.set(state)

					// Save to localStorage with debounce
					debouncedSave(state)

					// Sync UI state
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
						elements: (initialData.elements || []) as any,
						appState: {
							...(initialData.appState || {}),
							collaborators: new Map(),
						},
						files: (initialData.files || {}) as any,
					},
				})

				root.render(ExcalidrawComponent)
				excalidrawApp = root
			} catch (error) {
				console.error('Failed to load Excalidraw:', error)
			}
		}

		initExcalidraw()

		// Return cleanup function for onMount
		return () => {
			if (excalidrawApp) {
				excalidrawApp.unmount()
			}
		}
	})

	// Cleanup: flush any pending saves on unmount
	onDestroy(() => {
		if (saveTimeout) {
			clearTimeout(saveTimeout)
			// Immediately save any pending changes
			const currentState = get(drawingStore)
			saveDrawingById(drawingId, currentState)
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
