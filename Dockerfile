FROM golang:1.21-alpine as build
WORKDIR /app
ADD go.mod ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /usr/bin/app

FROM gcr.io/distroless/static
COPY --from=build /usr/bin/app /app
USER nobody:nobody
ENV TLS_CERT=""
ENV TLS_KEY=""
EXPOSE 8080
ENTRYPOINT [ "./app" ]