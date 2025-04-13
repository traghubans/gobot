FROM node:20-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the source code
COPY . .

# Build the application
RUN npm run build

# Install serve to run the production build
RUN npm install -g serve

# Expose port 3000
EXPOSE 3000

# Command to serve the built application
CMD ["serve", "-s", "build", "-l", "3000"] 