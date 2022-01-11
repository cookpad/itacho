FROM public.ecr.aws/docker/library/golang:1.14 as build

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/
RUN go mod verify && make


FROM public.ecr.aws/ubuntu/ubuntu:18.04

RUN apt-get update && \
      apt-get install -y ca-certificates && \
      apt-get clean && \
      rm -rf /var/lib/apt/lists/*

COPY --from=build /app/itacho /usr/local/bin/itacho
CMD ["/usr/local/bin/itacho", "server"]
