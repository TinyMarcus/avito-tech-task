# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app ./cmd/dynamic-user-segmentation-service

ENV PORT=8080
ENV HOST=postgres
ENV DBPORT=5432
ENV NAME="dynamic-user-segmentation"
ENV USER=postgres
ENV PASSWORD=postgres
ENV TYPE=postgres
EXPOSE ${PORT}

ENTRYPOINT ["./app"]