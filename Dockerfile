FROM golang AS build
WORKDIR /usr/src/app
ENV CGO_ENABLED=0
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN go build

FROM scratch
COPY --from=build /usr/src/app/local-gce-metadata /local-gce-metadata
ENTRYPOINT ["/local-gce-metadata"]