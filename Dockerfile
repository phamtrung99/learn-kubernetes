ARG GOLANG_TAG=1.16.3-alpine3.13
FROM golang:${GOLANG_TAG}

WORKDIR /app

COPY . .

EXPOSE 3000

CMD ["go", "run", "./web/cmd/main.go"]
