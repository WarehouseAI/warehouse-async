FROM golang:1.22.2-bullseye as build-deps

WORKDIR /usr/src/backend

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .
RUN go build /usr/src/backend/cmd/warehouse/main.go

FROM alpine:3.19.1
WORKDIR /usr/src/app
ARG env

COPY --from=build-deps /usr/src/backend/run.sh run.sh
COPY --from=build-deps /usr/src/backend/main main
COPY --from=build-deps /usr/src/backend/configs/$env configs/
RUN chmod +x run.sh
RUN apk add --no-cache bash
RUN apk add --no-cache libc6-compat

ARG module 
ENV LOG_PATH=/logs/$module.log

ENTRYPOINT ["./run.sh"]