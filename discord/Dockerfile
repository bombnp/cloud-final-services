FROM node:16-alpine AS base

WORKDIR /app

COPY ["package.json", "yarn.lock", "./"]

# Install dependencies
RUN yarn --frozen-lockfile

COPY . .

# Starting scripts
CMD yarn start
