import { writable } from 'svelte/store'
import type { ID } from '$lib/types'

export interface Drawing {
	id: ID
	name: string
	createdAt: Date
	updatedAt: Date
}

// Mock data - 8 sample drawings
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
	const { subscribe, set, update } = writable<Drawing[]>(mockDrawingsData)

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
			update((drawings) => [...drawings, drawing])
		},
		updateTimestamp: (id: ID) => {
			update((drawings) =>
				drawings.map((d) => (d.id === id ? { ...d, updatedAt: new Date() } : d))
			)
		},
		deleteDrawing: (id: ID) => {
			update((drawings) => drawings.filter((d) => d.id !== id))
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
