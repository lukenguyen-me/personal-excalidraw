import { writable } from 'svelte/store'

export type ToolType =
	| 'selection'
	| 'rectangle'
	| 'diamond'
	| 'ellipse'
	| 'arrow'
	| 'line'
	| 'freedraw'
	| 'text'
	| 'image'
	| 'eraser'
	| 'hand'

export interface UIState {
	sidebarOpen: boolean
	zoom: number
	toolbarVisible: boolean
	activeTool: ToolType | null
	viewMode: boolean
	zenMode: boolean
	gridMode: boolean
}

const defaultState: UIState = {
	sidebarOpen: true,
	zoom: 100,
	toolbarVisible: true,
	activeTool: 'selection',
	viewMode: false,
	zenMode: false,
	gridMode: false,
}

function createUIStore() {
	const { subscribe, set, update } = writable<UIState>(defaultState)

	return {
		subscribe,
		set,
		update,
		toggleSidebar: () => {
			update((state) => ({ ...state, sidebarOpen: !state.sidebarOpen }))
		},
		setZoom: (zoom: number) => {
			update((state) => ({ ...state, zoom: Math.max(10, Math.min(500, zoom)) }))
		},
		toggleToolbar: () => {
			update((state) => ({ ...state, toolbarVisible: !state.toolbarVisible }))
		},
		setActiveTool: (tool: ToolType) => {
			update((state) => ({ ...state, activeTool: tool }))
		},
		toggleViewMode: () => {
			update((state) => ({ ...state, viewMode: !state.viewMode }))
		},
		toggleZenMode: () => {
			update((state) => ({ ...state, zenMode: !state.zenMode }))
		},
		toggleGridMode: () => {
			update((state) => ({ ...state, gridMode: !state.gridMode }))
		},
	}
}

export const uiStore = createUIStore()
