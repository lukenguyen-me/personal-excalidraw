# Multi-stage build for Personal Excalidraw

# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app/frontend

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY frontend/ .
RUN pnpm run build

# Stage 2: Build backend
FROM golang:1.25-alpine3.22 AS backend-builder

RUN apk add --no-cache git

WORKDIR /app

COPY backend/ ./backend/

# Build the backend binary
RUN cd backend && go build -o server cmd/server/main.go

# Stage 3: Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy backend binary
COPY --from=backend-builder /app/backend/server ./

EXPOSE 8080

CMD ["./server"]
