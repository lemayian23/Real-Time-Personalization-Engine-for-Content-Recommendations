# Multi-stage build for frontend
FROM node:18-alpine as builder

WORKDIR /app

# Copy package files first
COPY package*.json ./

# Clean install dependencies (don't use cache to avoid platform issues)
RUN npm ci --no-cache

# Copy source code
COPY . .

# Build the app
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built app to nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]