# Stage 1: Build
FROM ghcr.io/etkecc/base/build AS builder

WORKDIR /app
COPY . .

# We remove the vendor folder just in case it was copied
RUN rm -rf vendor

# Build the binary directly to ensure it works
RUN go build -ldflags '-extldflags "-static"' -tags timetzdata,goolm -o postmoogle ./cmd/postmoogle

# Stage 2: Final Image
FROM scratch

# Set the DB path as a fallback
ENV POSTMOOGLE_DB_DSN=/data/postmoogle.db

# Copy certs and the binary we just built
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/postmoogle /bin/postmoogle

# We run as root (UID 0) to avoid "user not found" errors on Railway
USER 0

ENTRYPOINT ["/bin/postmoogle"]