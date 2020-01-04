# Only edit this file if it's in the base folder

FROM golang:1.10 as build

ARG workshop

WORKDIR /go/src/workshops/$workshop
COPY . .
RUN make install


FROM gcr.io/distroless/static:latest
COPY --from=build /go/bin/hello /

EXPOSE 8080
CMD ["/hello"]
