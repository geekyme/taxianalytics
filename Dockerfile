FROM golang:1.11.2

WORKDIR /taxianalytics

COPY . .

RUN go mod download

RUN go build -o taxianalytics

EXPOSE 8080

CMD ["./taxianalytics"]