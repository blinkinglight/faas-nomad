FROM golang:1.18-alpine as compiler

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY handlers ./handlers
COPY nomad ./nomad
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /service


FROM scratch
WORKDIR /
COPY --from=compiler /service /

ENTRYPOINT ["/service" ]

