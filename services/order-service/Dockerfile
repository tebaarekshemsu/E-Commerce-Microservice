# Use official Node.js LTS alpine image
FROM node:18-alpine

# Set working directory
WORKDIR /home/app

# Copy only package files first (for caching)
COPY package*.json ./

# Install dependencies (production only)
RUN npm ci --only=production

# Copy the rest of the application
COPY . .

# Expose port
EXPOSE 8300

# Set environment variables
ENV NODE_ENV=production
ENV PORT=8300

# Start the application
CMD ["node", "server.js"]
