import { writable } from 'svelte/store';
import type { ID } from '$lib/types';

export interface DrawingState {
	elements: unknown[];
	appState: unknown;
	files: Record<string, unknown>;
}

const defaultState: DrawingState = {
	elements: [],
	appState: {},
	files: {}
};

export const drawingStore = writable<DrawingState>(defaultState);

export function saveDrawing(state: DrawingState) {
	localStorage.setItem('excalidraw-drawing', JSON.stringify(state));
	drawingStore.set(state);
}

export function loadDrawing(): DrawingState {
	const saved = localStorage.getItem('excalidraw-drawing');
	const state = saved ? JSON.parse(saved) : defaultState;
	drawingStore.set(state);
	return state;
}

export function clearDrawing() {
	localStorage.removeItem('excalidraw-drawing');
	drawingStore.set(defaultState);
}

export function saveDrawingById(id: ID, state: DrawingState): boolean {
	const key = `excalidraw-drawing-${id}`;
	try {
		localStorage.setItem(key, JSON.stringify(state));
		drawingStore.set(state);

		// Sync metadata timestamp - will be imported after mockDrawings.ts is updated
		// to avoid circular dependency issues during development
		if (typeof window !== 'undefined' && window.__updateDrawingTimestamp) {
			window.__updateDrawingTimestamp(id);
		}

		return true;
	} catch (error: unknown) {
		console.error('Failed to save drawing:', error);
		return false;
	}
}

export function loadDrawingById(id: ID): DrawingState {
	const key = `excalidraw-drawing-${id}`;
	try {
		const saved = localStorage.getItem(key);
		if (!saved) {
			drawingStore.set(defaultState);
			return defaultState;
		}

		const state = JSON.parse(saved);

		// Validate structure
		if (!state.elements || !Array.isArray(state.elements)) {
			console.warn('Invalid drawing data, using default state');
			localStorage.removeItem(key);
			drawingStore.set(defaultState);
			return defaultState;
		}

		drawingStore.set(state);
		return state;
	} catch (error) {
		console.error('Failed to load drawing:', error);
		localStorage.removeItem(key);
		drawingStore.set(defaultState);
		return defaultState;
	}
}

export function deleteDrawingById(id: ID): void {
	const key = `excalidraw-drawing-${id}`;
	localStorage.removeItem(key);
}
