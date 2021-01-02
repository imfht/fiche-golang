FROM golang:alpine as builder
RUN apk update && apk add ca-certificates
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -o fiche-go .
FROM alpine
ENV LANG en_US.UTF-8
ENV LC_ALL en_US.UTF-8
COPY --from=builder /build/fiche-go /app/
COPY --from=builder /build/data/ /app/data
WORKDIR /app
CMD ["./fiche-go"]
