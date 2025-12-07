import { writable } from 'svelte/store'
import { browser } from '$app/environment'

interface AuthState {
	accessKey: string | null
	isAuthenticated: boolean
}

const STORAGE_KEY = 'excalidraw_access_key'

// Load initial state from localStorage (browser only)
function loadInitialState(): AuthState {
	if (browser) {
		const storedKey = localStorage.getItem(STORAGE_KEY)
		return {
			accessKey: storedKey,
			isAuthenticated: storedKey !== null
		}
	}
	return {
		accessKey: null,
		isAuthenticated: false
	}
}

function createAuthStore() {
	const { subscribe, set } = writable<AuthState>(loadInitialState())

	return {
		subscribe,

		getAccessKey(): string | null {
			let key: string | null = null
			subscribe((state) => {
				key = state.accessKey
			})()
			return key
		},

		hasAccessKey(): boolean {
			let hasKey = false
			subscribe((state) => {
				hasKey = state.isAuthenticated
			})()
			return hasKey
		},

		setAccessKey(key: string): void {
			if (browser) {
				localStorage.setItem(STORAGE_KEY, key)
			}
			set({ accessKey: key, isAuthenticated: true })
		},

		clearAccessKey(): void {
			if (browser) {
				localStorage.removeItem(STORAGE_KEY)
			}
			set({ accessKey: null, isAuthenticated: false })
		}
	}
}

export const authStore = createAuthStore()

// Listen for storage changes across tabs (browser only)
if (browser) {
	window.addEventListener('storage', (e) => {
		if (e.key === STORAGE_KEY) {
			if (e.newValue) {
				authStore.setAccessKey(e.newValue)
			} else {
				authStore.clearAccessKey()
			}
		}
	})
}
