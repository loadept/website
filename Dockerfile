FROM golang:1.26-alpine3.23 AS compile
WORKDIR /usr/src/app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -o server
RUN go build -trimpath -o migrate ./cmd/migrate
RUN go build -trimpath -o shortener ./cmd/shortener

FROM gcr.io/distroless/static-debian12:nonroot AS release
WORKDIR /app
COPY --from=compile /usr/src/app/server .
COPY --from=compile /usr/src/app/migrate /usr/src/app/shortener /

ENTRYPOINT [ "/app/server" ]
