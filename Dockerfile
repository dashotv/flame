############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

WORKDIR /go/src/app
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,target=. \
  go install

############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
WORKDIR /root/
COPY --from=builder /go/bin/flame .
COPY .env.vault .
CMD ["./flame", "server"]
