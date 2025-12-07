-- Add slug column to drawings table
ALTER TABLE drawings ADD COLUMN slug VARCHAR(50) NOT NULL DEFAULT '';

-- Create unique index on slug for fast lookups and uniqueness constraint
CREATE UNIQUE INDEX idx_drawings_slug ON drawings(slug) WHERE slug != '';
