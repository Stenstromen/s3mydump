FROM golang:alpine AS build
WORKDIR /src
ADD . /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o mydump/mydump ./mydump
RUN chmod +x /src/mydump/mydump

FROM scratch
COPY --from=build /src/mydump/mydump /mydump
USER 65534:65534
ENTRYPOINT ["/mydump"]
