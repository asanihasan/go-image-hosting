# -------- Dockerfile ----------
FROM golang:1.24-alpine

WORKDIR /app

# copy module definition first and download deps (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# now copy the source
COPY . .

# make sure deps are still consistent after the full copy
RUN go mod tidy && go build -o server .

CMD ["./server"]
# --------------------------------
