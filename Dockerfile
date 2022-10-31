# Build the manager binary
FROM golang:1.19 as builder

ENV PORT=8080

WORKDIR /app
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY server/ server/
COPY .github_token .github_token

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o api .

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine:latest
WORKDIR /
RUN apk update && apk upgrade && apk add curl && apk add tar && apk add gzip
RUN curl -L --output - https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.4.7/argocd-autopilot-linux-amd64.tar.gz | tar zx
RUN mv ./argocd-autopilot-* /usr/local/bin/argocd-autopilot
RUN argocd-autopilot version

COPY --from=builder /app/api .
COPY --from=builder /app/.github_token .
USER 65532:65532

EXPOSE $PORT
ENTRYPOINT ["./api"]
