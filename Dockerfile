FROM oven/bun:1 AS bun
WORKDIR /usr/src/app

FROM bun AS statics
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile
COPY web .
RUN bun run build

FROM golang:1.25-alpine3.23 AS compile
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=statics /usr/src/app/static web/static
RUN go build -trimpath

FROM gcr.io/distroless/static-debian12:nonroot AS release
WORKDIR /app
COPY --from=compile /usr/src/app/loadept.com .

ENTRYPOINT [ "./loadept.com" ]
