# Stage 1: Modules
FROM golang:1.19-alpine3.16 as modules
WORKDIR /modules
COPY go.mod go.sum /modules/
RUN go mod download

# Stage 2: Builder
FROM golang:1.19-alpine3.16 as builder
RUN apk add --no-cache \
    libc6-compat=1.2.3-r3 \
    gcc=11.2.1_git20220219-r2 \
    g++=11.2.1_git20220219-r2

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app


# Stage 3: Runtime ENV
FROM golang:1.19-alpine3.16 AS runtime_env
WORKDIR /application
COPY --from=builder /app/credential.json ./credential.json
COPY --from=builder /app/config ./config
COPY --from=builder /app/.env .

# Stage 4: Build Service
FROM builder as build_core_service
RUN GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./api/.

# Stage 5: Runtime
FROM runtime_env AS core_service
COPY --from=build_core_service /bin/app ./bin/app
CMD ["./bin/app"]