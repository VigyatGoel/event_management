# Use the official Node.js image as the base image
FROM node:22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Install openssl
RUN apk add --no-cache openssl

# Copy the rest of the application code to the working directory
COPY . .

# Build the application
RUN npm run build

# Install `serve` globally to serve the built application
RUN npm install -g serve

# Expose the port the app runs on
EXPOSE 3000

# Start the application using `serve`
CMD ["serve", "-s", "dist", "-l", "3000"]
