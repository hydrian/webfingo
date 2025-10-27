FROM golang:1.25-alpine AS builder
RUN apk add make
WORKDIR /builder
COPY ./ ./
RUN make build

FROM scratch
#RUN apk add jq
WORKDIR /app
#RUN mkdir -p bin config
COPY --from=builder /builder/bin/webfingo /app/bin/webfingo
ENTRYPOINT ["/app/bin/webfingo"]
CMD [ "-config", "/app/config/config.json" ]
EXPOSE 8080
