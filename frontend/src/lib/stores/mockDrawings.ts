import { writable } from 'svelte/store'
import type { ID } from '$lib/types'

export interface Drawing {
	id: ID
	name: string
	createdAt: Date
	updatedAt: Date
}

// localStorage key for persisting drawing metadata
const METADATA_STORAGE_KEY = 'excalidraw-drawings-metadata'

// Maximum length for drawing names
const MAX_NAME_LENGTH = 100

// Helper function to save drawings metadata to localStorage
function saveDrawingsMetadata(drawings: Drawing[]) {
	// Check if we're running in a browser environment
	if (typeof window === 'undefined') {
		return
	}

	try {
		localStorage.setItem(METADATA_STORAGE_KEY, JSON.stringify(drawings))
	} catch (error) {
		if (error instanceof DOMException && error.name === 'QuotaExceededError') {
			console.error('localStorage quota exceeded for metadata')
		} else {
			console.error('Failed to save drawings metadata:', error)
		}
	}
}

// Helper function to load drawings metadata from localStorage
function loadDrawingsMetadata(): Drawing[] {
	// Check if we're running in a browser environment
	if (typeof window === 'undefined') {
		return mockDrawingsData
	}

	try {
		const saved = localStorage.getItem(METADATA_STORAGE_KEY)
		if (!saved) return mockDrawingsData

		const parsed = JSON.parse(saved)
		// Convert date strings back to Date objects
		return parsed.map((d: any) => ({
			...d,
			createdAt: new Date(d.createdAt),
			updatedAt: new Date(d.updatedAt)
		}))
	} catch (error) {
		console.error('Failed to load drawings metadata:', error)
		return mockDrawingsData
	}
}

// Mock data - 8 sample drawings (used as fallback if localStorage is empty)
const mockDrawingsData: Drawing[] = [
	{
		id: 1,
		name: 'System Architecture Diagram',
		createdAt: new Date('2024-11-15T10:30:00'),
		updatedAt: new Date('2024-11-20T14:45:00')
	},
	{
		id: 2,
		name: 'Wireframe - Login Page',
		createdAt: new Date('2024-11-18T09:15:00'),
		updatedAt: new Date('2024-11-18T16:20:00')
	},
	{
		id: 3,
		name: 'Database Schema',
		createdAt: new Date('2024-11-22T11:00:00'),
		updatedAt: new Date('2024-11-25T10:30:00')
	},
	{
		id: 4,
		name: 'User Flow Chart',
		createdAt: new Date('2024-11-23T13:45:00'),
		updatedAt: new Date('2024-11-28T09:00:00')
	},
	{
		id: 5,
		name: 'API Design Mockup',
		createdAt: new Date('2024-11-25T08:30:00'),
		updatedAt: new Date('2024-11-29T15:10:00')
	},
	{
		id: 6,
		name: 'Component Hierarchy',
		createdAt: new Date('2024-11-27T14:20:00'),
		updatedAt: new Date('2024-12-01T11:30:00')
	},
	{
		id: 7,
		name: 'Network Topology',
		createdAt: new Date('2024-11-28T10:00:00'),
		updatedAt: new Date('2024-12-02T13:45:00')
	},
	{
		id: 8,
		name: 'Class Diagram - Payment Module',
		createdAt: new Date('2024-11-30T16:30:00'),
		updatedAt: new Date('2024-12-03T10:15:00')
	}
]

function createDrawingsStore() {
	const { subscribe, set, update } = writable<Drawing[]>(loadDrawingsMetadata())

	return {
		subscribe,
		getDrawings: () => {
			let drawings: Drawing[] = []
			const unsubscribe = subscribe((value) => {
				drawings = value
			})
			unsubscribe()
			return drawings
		},
		getDrawing: (id: ID) => {
			let drawings: Drawing[] = []
			const unsubscribe = subscribe((value) => {
				drawings = value
			})
			unsubscribe()
			return drawings.find((d) => d.id === id)
		},
		addDrawing: (drawing: Drawing) => {
			update((drawings) => {
				const updated = [...drawings, drawing]
				saveDrawingsMetadata(updated)
				return updated
			})
		},
		updateDrawingName: (id: ID, newName: string) => {
			// Validate and sanitize input
			let sanitizedName = newName.trim()

			// Prevent empty names
			if (!sanitizedName) {
				return false
			}

			// Enforce maximum length
			if (sanitizedName.length > MAX_NAME_LENGTH) {
				sanitizedName = sanitizedName.substring(0, MAX_NAME_LENGTH)
			}

			// Update the drawing name and timestamp
			update((drawings) => {
				const updated = drawings.map((d) =>
					d.id === id ? { ...d, name: sanitizedName, updatedAt: new Date() } : d
				)
				saveDrawingsMetadata(updated)
				return updated
			})

			return true
		},
		updateTimestamp: (id: ID) => {
			update((drawings) => {
				const updated = drawings.map((d) =>
					d.id === id ? { ...d, updatedAt: new Date() } : d
				)
				saveDrawingsMetadata(updated)
				return updated
			})
		},
		deleteDrawing: (id: ID) => {
			update((drawings) => {
				const updated = drawings.filter((d) => d.id !== id)
				saveDrawingsMetadata(updated)
				return updated
			})
			// Also delete the drawing content from localStorage
			const { deleteDrawingById } = require('./drawing')
			deleteDrawingById(id)
		},
		reset: () => set(mockDrawingsData)
	}
}

export const drawingsStore = createDrawingsStore()

// Expose updateTimestamp globally for use in drawing.ts to avoid circular dependency
if (typeof window !== 'undefined') {
	;(window as any).__updateDrawingTimestamp = drawingsStore.updateTimestamp
}
