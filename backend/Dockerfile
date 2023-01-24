FROM golang:latest AS build
WORKDIR /build
COPY . .
RUN make build


FROM ubuntu:latest
WORKDIR /app
COPY --from=build /build/bin/chess .
EXPOSE 8080
CMD ["/app/chess"]