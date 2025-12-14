FROM node:18-alpine

ARG PROJECT_VERSION=0.1.0
WORKDIR /home/app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy application files
COPY . .

# Expose port
EXPOSE 8300

# Set environment variables
ENV NODE_ENV=production
ENV PORT=8300

# Start the application
CMD ["node", "server.js"]
