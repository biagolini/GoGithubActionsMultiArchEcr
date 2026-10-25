# Build stage: compiles the server natively for the runner's architecture.
FROM public.ecr.aws/docker/library/golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

# Runtime stage: small image with a non-root user.
FROM public.ecr.aws/docker/library/alpine:3.22
RUN adduser -D -u 10001 app
COPY --from=build /out/server /usr/local/bin/server
USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]
