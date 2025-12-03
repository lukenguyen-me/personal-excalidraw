import { writable } from 'svelte/store'

export interface DrawingState {
	elements: unknown[]
	appState: unknown
	files: Record<string, unknown>
}

const defaultState: DrawingState = {
	elements: [],
	appState: {},
	files: {},
}

export const drawingStore = writable<DrawingState>(defaultState)

export function saveDrawing(state: DrawingState) {
	localStorage.setItem('excalidraw-drawing', JSON.stringify(state))
	drawingStore.set(state)
}

export function loadDrawing(): DrawingState {
	const saved = localStorage.getItem('excalidraw-drawing')
	const state = saved ? JSON.parse(saved) : defaultState
	drawingStore.set(state)
	return state
}

export function clearDrawing() {
	localStorage.removeItem('excalidraw-drawing')
	drawingStore.set(defaultState)
}
