import { writable } from 'svelte/store';

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
