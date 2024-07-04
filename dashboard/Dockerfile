FROM node:20 AS base
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the application
COPY . .

# setting mongodb connection string
ARG MONGODB_URI
ARG MONGO_URI
ARG REDIS_URL
ENV MONGODB_URI=$MONGODB_URI
ENV MONGO_URI=$MONGODB_URI
ENV REDIS_URL=$REDIS_URL

# Build the application
RUN npm run build

FROM node:20-alpine3.19 AS release
WORKDIR /app

# Copy necessary files from the base stage
COPY --from=base /app/node_modules ./node_modules
COPY --from=base /app/package.json ./package.json
COPY --from=base /app/.next ./.next
COPY --from=base /app/public ./public

EXPOSE 3000

# Start the application
CMD ["npm", "start"]