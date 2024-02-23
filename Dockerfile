FROM golang:1.22-alpine as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/bin/app

FROM gcr.io/distroless/static
COPY --from=build /usr/bin/app /app
USER nobody:nobody
ENV PORT=8080
ENV GRPC=8081
EXPOSE ${PORT}
EXPOSE ${GRPC}
ENTRYPOINT [ "./app" ]