FROM  golang:1.21-alpine as builder

WORKDIR /project

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN go build -tags http -o /project/build/app .


FROM alpine:latest

# to fix timezone not loaded
RUN apk add --no-cache tzdata

COPY --from=builder /project/build/app /project/build/app

WORKDIR /project/build/

EXPOSE 8080
CMD [ "./app" ]
