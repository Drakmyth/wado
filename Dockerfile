# Base Stage
FROM golang:1.22.3 as base
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify


# Dev Stage
FROM base as dev
WORKDIR /app

RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

COPY .air.toml ./
CMD ["air", "-c", ".air.toml"]


# Production Build Stage
FROM base as build
WORKDIR /app

RUN useradd -u 1001 appuser

COPY . ./
RUN go build -ldflags="-linkmode external -extldflags -static" -o ./bin/go-template


# Production Release Stage
FROM scratch
WORKDIR /app

ENV GIN_MODE=release

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /app/bin/go-template ./go-template
COPY --from=build /app/public/ ./public/
COPY --from=build /app/templates/ ./templates/

USER appuser
EXPOSE 8080

CMD ["./go-template"]