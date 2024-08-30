FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download

# Step 2: Builder
FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /note-keeper

# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
CMD ["/app"]