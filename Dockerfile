FROM golang:1.19

ENV GOPATH=/
COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh


RUN go mod download
RUN go build -o ordering-bot-migrate ./cmd/migrate/main.go
RUN go build -o ordering-bot ./cmd/vkbot/main.go
CMD ["./ordering-bot-migrate"]
CMD ["./ordering-bot"]