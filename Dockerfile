FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . .

RUN go mod download
RUN go build -o app ./cmd
RUN apk --no-cache add curl
EXPOSE ${TODOLIST_PORT}

CMD [ "./app" ]