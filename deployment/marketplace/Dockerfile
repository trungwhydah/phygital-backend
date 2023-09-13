# Stage 1: Modules
FROM golang:1.20.7-alpine3.18 as modules
WORKDIR /modules
COPY go.mod go.sum /modules/
RUN go mod download

# Stage 2: Builder
FROM golang:1.20.7-alpine3.18 as builder
RUN apk add --no-cache \
    libc6-compat=1.2.4-r1 \
    gcc=12.2.1_git20220924-r10 \
    g++=12.2.1_git20220924-r10 \
    tzdata=2023c-r1

ENV TZ=Asia/Bangkok
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

# Stage 3: Runtime ENV
FROM golang:1.20.7-alpine3.18 AS runtime_env
WORKDIR /application
COPY --from=builder /app/service-account-file.json ./service-account-file.json
COPY --from=builder /app/config ./config
COPY --from=builder /app/translation.en.yaml .
COPY --from=builder /app/translation.vi.yaml .

# Stage 4: Build Service
FROM builder as build_core_service
RUN GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/marketplace

# Stage 5: Runtime
FROM runtime_env AS core_service
COPY --from=build_core_service /bin/app ./bin/app
CMD ["./bin/app"]
