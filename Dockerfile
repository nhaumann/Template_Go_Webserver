FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Enable CGO for sqlite3
ENV CGO_ENABLED=1
# Install C compiler
RUN apk add --no-cache gcc musl-dev
# Build
RUN go build -o /app/Glossary

FROM alpine:latest

# Set default values for environment variables
ENV WEB_PORT=:8080
ENV STATIC_CONTENT_PREFIX=/static/
ENV WEB_PATH=/web
ENV TEMPLATES_PATH=/templates
ENV WEB_ROOT=/home
ENV NOT_FOUND_FILE_NAME_PATH=/templates/notefound.html


# Copy the binary to the root directory of the container
COPY --from=builder /app/Glossary /Glossary

# Copy the web directory to the root directory of the container
COPY --from=builder /app/web /web

# Copy the Glossary.db file to the data directory in the container
COPY --from=builder /app/data/Glossary.db /data/Glossary.db

# Copy the .env file to the root directory of the container
COPY --from=builder /app/.env /.env

# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# Set the working directory to the root directory of the container
WORKDIR /

#log directory files
# Run the binary at /Glossary when the container starts up
CMD ["/Glossary", "serve"]