FROM golang:1.11 as build

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum vendor /app/
RUN go mod download

COPY . /app/
RUN go mod verify && make

FROM gcr.io/distroless/base

COPY --from=build /app/itacho /usr/local/bin/itacho
CMD ["/usr/local/bin/itacho", "server"]
