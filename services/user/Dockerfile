ARG RELEASE_MODE=prod

FROM golang:1.23.3-bookworm AS base


WORKDIR /app

COPY ./services/user .
COPY ./common/rpc /app/rpc

RUN go mod tidy

RUN go build -o main .

FROM base AS dev
RUN go install github.com/air-verse/air@latest
ENTRYPOINT [ "air" ]

FROM debian:bookworm-slim AS prod

WORKDIR /usr/app

COPY --from=base /app/main /usr/app/main

RUN chmod +x /usr/app/main

ENTRYPOINT ["/usr/app/main"]

FROM ${RELEASE_MODE}