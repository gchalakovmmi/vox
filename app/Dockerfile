FROM alpine:latest
# RUN apk add --no-cache sqlite
WORKDIR /app
COPY cmd/app/app .
COPY web/static/ ./web/static
# COPY data/app.db ./data/app.db
EXPOSE 8080
ENTRYPOINT ["/app/app"]
