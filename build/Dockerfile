FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

#RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/kubernetes-git-sync


FROM gcr.io/distroless/static:latest

USER nonroot
COPY --from=builder /app/kubernetes-git-sync .

CMD ["./kubernetes-git-sync"]