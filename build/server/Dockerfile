FROM golang:1.23 AS build
WORKDIR /go/src/server
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/server ./cmd/server

FROM gcr.io/distroless/static-debian12 
COPY --from=build /go/bin/server /
CMD ["/server"]
