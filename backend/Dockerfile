FROM golang:1.24-alpine

WORKDIR /app
COPY backend/ ./
RUN go build -o server .

EXPOSE 8080
CMD ["./server"]

