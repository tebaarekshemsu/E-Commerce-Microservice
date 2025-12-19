# 1. Use Node.js base image
FROM node:18-alpine

# 2. Set working directory inside container
WORKDIR /app

# 3. Copy package files first (for caching)
COPY package*.json ./

# 4. Install dependencies
RUN npm install

# 5. Copy the rest of the source code
COPY . .

# 6. Expose port (use the same port your app uses)
EXPOSE 4000

# 7. Start the app
CMD ["npm", "run"," dev"]
