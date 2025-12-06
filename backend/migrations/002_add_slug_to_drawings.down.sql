-- Drop the index and column
DROP INDEX IF EXISTS idx_drawings_slug;
ALTER TABLE drawings DROP COLUMN IF EXISTS slug;
