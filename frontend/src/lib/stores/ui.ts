import { writable } from 'svelte/store'

export interface UIState {
	sidebarOpen: boolean
	theme: 'light' | 'dark'
	zoom: number
}

const defaultState: UIState = {
	sidebarOpen: true,
	theme: 'light',
	zoom: 100,
}

function createUIStore() {
	const { subscribe, set, update } = writable<UIState>(defaultState)

	// Load theme from localStorage on initialization
	if (typeof window !== 'undefined') {
		const savedTheme = localStorage.getItem('excalidraw-theme') as 'light' | 'dark'
		if (savedTheme) {
			set({
				...defaultState,
				theme: savedTheme,
			})
		}
	}

	return {
		subscribe,
		set,
		update,
		toggleTheme: () => {
			update((state) => {
				const newTheme = state.theme === 'light' ? 'dark' : 'light'
				if (typeof window !== 'undefined') {
					localStorage.setItem('excalidraw-theme', newTheme)
				}
				return { ...state, theme: newTheme }
			})
		},
		toggleSidebar: () => {
			update((state) => ({ ...state, sidebarOpen: !state.sidebarOpen }))
		},
		setZoom: (zoom: number) => {
			update((state) => ({ ...state, zoom: Math.max(10, Math.min(500, zoom)) }))
		},
	}
}

export const uiStore = createUIStore()
