FROM golang:alpine AS builder

WORKDIR /build

COPY ./ ./
RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -o ./snippet-war ./cmd/snippet-war

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /build/snippet-war ./

EXPOSE 8081

CMD ["./snippet-war"]