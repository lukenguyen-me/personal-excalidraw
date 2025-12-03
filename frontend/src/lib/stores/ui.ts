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

export const uiStore = writable<UIState>(defaultState)

export function toggleSidebar() {
  uiStore.update((state) => ({
    ...state,
    sidebarOpen: !state.sidebarOpen,
  }))
}

export function setTheme(theme: 'light' | 'dark') {
  uiStore.update((state) => ({
    ...state,
    theme,
  }))
  localStorage.setItem('excalidraw-theme', theme)
}

export function setZoom(zoom: number) {
  uiStore.update((state) => ({
    ...state,
    zoom: Math.max(10, Math.min(500, zoom)),
  }))
}

export function loadUIState() {
  const savedTheme = (localStorage.getItem('excalidraw-theme') as 'light' | 'dark') || 'light'
  uiStore.set({
    ...defaultState,
    theme: savedTheme,
  })
}
