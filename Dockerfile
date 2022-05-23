FROM golang:latest as build-stage
# Is used to indicate that this image contains test results
LABEL builder=true
# built binaries, test results and coverage report are saved here
ARG BUILD_DIR=/app/build

# Custom addition to enable compile windows64 binary
RUN apt-get update && apt-get install -y gcc-multilib && apt-get install -y gcc-mingw-w64

# Setup directories
WORKDIR /app
COPY ./ /app/
RUN mkdir -p $BUILD_DIR

# Lint code
RUN make vet
RUN make fmt-check
RUN make lint

# build, run tests and export test results ./build directory
RUN make

FROM alpine:latest as deploy-stage

ARG CITY_NAME=london
ARG BUILD_DIR=/app/build

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build-stage ${BUILD_DIR}/linux64/${CITY_NAME}-server ./server

EXPOSE 80

CMD ["./server", "-p", "80"]
