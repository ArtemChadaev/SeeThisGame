FROM n8nio/n8n:latest

USER root

# Install useful tools for "Execute Command" nodes
RUN apk add --no-cache git curl bash

USER node

# Install any additional packages if needed
# RUN apk add --no-cache ...

USER node

# Persistence directory
VOLUME /home/node/.n8n

EXPOSE 5678
