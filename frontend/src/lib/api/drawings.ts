// API Types (matching backend)
export interface DrawingDTO {
	id: string
	name: string
	data: Record<string, unknown>
	created_at: string
	updated_at: string
}

export interface DrawingListResponse {
	drawings: DrawingDTO[]
	total: number
	limit: number
	offset: number
}

export interface CreateDrawingRequest {
	name: string
	data: Record<string, unknown>
}

export interface UpdateDrawingRequest {
	name?: string
	data?: Record<string, unknown>
}

// API Client
class DrawingsAPI {
	private baseURL = '/api'

	async list(params?: { limit?: number; offset?: number }): Promise<DrawingListResponse> {
		const query = new URLSearchParams({
			limit: String(params?.limit || 10),
			offset: String(params?.offset || 0)
		})
		const response = await fetch(`${this.baseURL}/drawings?${query}`)
		if (!response.ok) throw new Error('Failed to fetch drawings')
		return response.json()
	}

	async get(id: string): Promise<DrawingDTO> {
		const response = await fetch(`${this.baseURL}/drawings/${id}`)
		if (!response.ok) throw new Error('Failed to fetch drawing')
		return response.json()
	}

	async create(data: CreateDrawingRequest): Promise<DrawingDTO> {
		const response = await fetch(`${this.baseURL}/drawings`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(data)
		})
		if (!response.ok) throw new Error('Failed to create drawing')
		return response.json()
	}

	async update(id: string, data: UpdateDrawingRequest): Promise<DrawingDTO> {
		const response = await fetch(`${this.baseURL}/drawings/${id}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(data)
		})
		if (!response.ok) throw new Error('Failed to update drawing')
		return response.json()
	}

	async delete(id: string): Promise<void> {
		const response = await fetch(`${this.baseURL}/drawings/${id}`, {
			method: 'DELETE'
		})
		if (!response.ok) throw new Error('Failed to delete drawing')
	}
}

export const drawingsAPI = new DrawingsAPI()
