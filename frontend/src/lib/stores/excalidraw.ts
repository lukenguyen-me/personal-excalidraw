import { writable } from 'svelte/store';
import type { ExcalidrawImperativeAPI } from '@excalidraw/excalidraw/types/types';

interface ExcalidrawStore {
	api: ExcalidrawImperativeAPI | null;
	ready: boolean;
}

const createExcalidrawStore = () => {
	const { subscribe, set, update } = writable<ExcalidrawStore>({
		api: null,
		ready: false
	});

	return {
		subscribe,
		setAPI: (api: ExcalidrawImperativeAPI) => {
			update((state) => ({
				...state,
				api,
				ready: true
			}));
		},
		clearAPI: () => {
			set({
				api: null,
				ready: false
			});
		},
		reset: () => {
			set({
				api: null,
				ready: false
			});
		}
	};
};

export const excalidrawStore = createExcalidrawStore();
