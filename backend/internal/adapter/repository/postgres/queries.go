package postgres

const (
	// queryCreateDrawing inserts a new drawing into the database
	queryCreateDrawing = `
		INSERT INTO drawings (id, name, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	// queryFindDrawingByID retrieves a drawing by its ID
	queryFindDrawingByID = `
		SELECT id, name, data, created_at, updated_at
		FROM drawings
		WHERE id = $1
	`

	// queryFindAllDrawings retrieves all drawings with pagination
	queryFindAllDrawings = `
		SELECT id, name, data, created_at, updated_at
		FROM drawings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	// queryUpdateDrawing updates an existing drawing
	queryUpdateDrawing = `
		UPDATE drawings
		SET name = $1, data = $2, updated_at = $3
		WHERE id = $4
	`

	// queryDeleteDrawing deletes a drawing by ID
	queryDeleteDrawing = `
		DELETE FROM drawings
		WHERE id = $1
	`

	// queryCountDrawings returns the total number of drawings
	queryCountDrawings = `
		SELECT COUNT(*)
		FROM drawings
	`
)
