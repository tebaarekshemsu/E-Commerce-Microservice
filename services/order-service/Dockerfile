# Use official Node.js LTS Alpine image
FROM node:18-alpine

# Set working directory
WORKDIR /home/app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy all application files
COPY . .

# Expose port
EXPOSE 8300

# Set environment variables
ENV NODE_ENV=production
ENV PORT=8300

# Start the application
CMD ["node", "server.js"]
