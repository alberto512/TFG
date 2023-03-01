FROM golang:1.20.1-alpine AS build
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /go/bin/app /
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]