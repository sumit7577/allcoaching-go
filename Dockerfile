FROM golang:alpine As build

RUN apk --no-cache add gcc g++ make git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN ls -la

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /app/bin /go/bin
COPY --from=build ./app/conf /usr/bin/conf

EXPOSE 8080
ENTRYPOINT /go/bin/web-app --port 8080