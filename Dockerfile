FROM golang:1.22 AS base

WORKDIR /app

COPY go.* ./

ARG GH_ACCESS_TOKEN
ENV GH_ACCESS_TOKEN=${GH_ACCESS_TOKEN}

RUN git config --global --add url."https://${GH_ACCESS_TOKEN}@github.com/".insteadOf "https://github.com/"

ENV GOPRIVATE=github.com/adhfoundation

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
  swag init -g ./adapters/http/server.go --parseDependency --parseInternal

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 
  
ENV GOPATH /go

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ./application ./cmd/main.go 

FROM alpine:latest AS production

RUN apk update && apk add gcompat

EXPOSE 8080

RUN mkdir /app && \
  addgroup -S appgroup && \
  adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=base /app/application /app/
COPY --from=base /app/docs /app/docs
COPY --from=base /app/bootstrap/db /app/ 
COPY --from=base /go/bin/migrate /app/ 

USER appuser

CMD ["./application"]
