
FROM golang:1.24-alpine AS builder


RUN apk add --no-cache git


WORKDIR /app


COPY . .


RUN go mod tidy && CGO_ENABLED=0 go build -o app .


FROM gcr.io/distroless/static-debian12:debug-nonroot


WORKDIR /app


COPY --from=builder --chown=nonroot:nonroot /app/app /app/app


EXPOSE 8080


ENTRYPOINT ["/app/app"]


HEALTHCHECK --interval=30s --retries=3 --start-period=10s --timeout=10s \
  CMD curl -f http://localhost:8080/livez || exit 1