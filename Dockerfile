FROM --platform=$BUILDPLATFORM node:alpine AS npm_builder
WORKDIR /app
COPY ui ./ui
RUN npm --prefix /app/ui install
RUN npm --prefix /app/ui run build

FROM --platform=$BUILDPLATFORM golang:alpine AS go_builder
ARG TARGETOS=linux
ARG TARGETARCH
ARG VERSION=dev
WORKDIR /app
RUN apk update && apk add --no-cache git
COPY cmd ./cmd
COPY go.mod ./go.mod
COPY internal ./internal
COPY --from=npm_builder /app/ui/dist ./cmd/shokkupan/dist
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s -X github.com/randspace0/shokkupan/internal/env.Version=${VERSION}" -o /go/bin/shokkupan ./cmd/shokkupan

FROM gcr.io/distroless/static:nonroot
COPY --from=go_builder /go/bin/shokkupan /go/bin/shokkupan
ENTRYPOINT ["/go/bin/shokkupan"]