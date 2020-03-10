# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM golang:1.14 as builder
RUN             apt update && apt install -y git gcc musl-dev make libportmidi-dev
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/music-paint
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM            debian:buster
RUN             apt update && apt install -y libportmidi-dev
LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="music-paint" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/music-paint/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/music-paint" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/music-paint" \
                org.label-schema.help="docker exec -it $CONTAINER music-paint --help"
COPY            --from=builder /go/bin/music-paint /bin/
ENTRYPOINT      ["/bin/music-paint"]
#CMD             []
