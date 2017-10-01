FROM node:boron-alpine

# Make sure git is available
RUN apk add --update git && \
  rm -rf /tmp/* /var/cache/apk/*

# Set up environment variables
ENV NODE_ENV production
ENV BUILD_VERSION 1

# Create app directory
RUN mkdir -p /app
WORKDIR /app

# Bundle app source
COPY . /app

# Install app dependencies, build production bundle, and clean up afterwards
RUN yarn && \
  yarn run build && \
  yarn run prune && \
  yarn cache clean

EXPOSE 3000
CMD [ "yarn", "run", "start" ]
