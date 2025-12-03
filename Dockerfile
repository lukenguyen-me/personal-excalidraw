# Multi-stage build for Personal Excalidraw

# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

# Stage 2: Build backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

COPY backend/ ./backend/
RUN cd backend && go build -o server main.go

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
