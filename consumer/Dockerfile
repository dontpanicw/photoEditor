#FROM golang:1.22.3-alpine as build
#
#WORKDIR /build
#
#COPY . .
#
#RUN go build -o bin/consumer consumer/main.go
#
#FROM alpine AS runner
#
#COPY --from=build /build/bin/consumer /consumer
#
#CMD ["/app/consumer"]


#FROM golang:1.23.4-alpine as build
#
#WORKDIR /app
#
#COPY go.mod go.sum ./
#
#RUN go mod download
#
#COPY consumer/ ./consumer
#
#RUN go build -o bin/consumer ./consumer
#
#FROM alpine AS runner
#
#COPY --from=build /app/bin/consumer /app/consumer
#
#CMD ["/app/consumer"]



FROM golang:1.23.4-alpine as build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY consumer/ ./consumer

RUN apk add --no-cache make

RUN go build -o bin/consumer ./consumer

FROM alpine AS runner

WORKDIR app

RUN apk add --no-cache curl

COPY --from=build /build/bin/consumer /build/consumer

CMD ["/build/consumer"]