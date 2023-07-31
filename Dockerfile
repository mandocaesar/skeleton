# Build stage
FROM golang:1.18-alpine AS builder

# Set ARGs
ARG APP_NAME=$APP_NAME

# Set workdir
WORKDIR /app

# Copy all project code
ADD . .

RUN apk update && apk add git && apk add openssh

RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Set Go Private
RUN go env -w GOPRIVATE=github.com/machtwatch/catalystdk*

# Download dependencies
RUN --mount=type=ssh go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o /tmp/app main.go

# Final stage
FROM alpine:latest AS production


# Copy output binary file from build stage
COPY --from=builder /tmp/app .
COPY --from=builder /app/flag.yml .
COPY --from=builder /app/files ./files

# Run the app
CMD ["./app"]
