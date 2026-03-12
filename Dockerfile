# --- STAGE 1: Build the binary ---
# We use the official etke.cc build image because it has the 
# necessary C libraries for Matrix encryption pre-installed.
FROM ghcr.io/etkecc/base/build:latest AS builder

WORKDIR /app

# Copy your modified source code
COPY . .

# IMPORTANT: Remove the vendor folder so the compiler 
# sees the new fields you added to the Config structs.
RUN rm -rf vendor

# Build a static binary with the required Matrix tags (goolm)
RUN go build -ldflags '-extldflags "-static"' -tags timetzdata,goolm -o postmoogle ./cmd/postmoogle

# --- STAGE 2: Create the final production image ---
FROM scratch

# 1. Copy CA Certificates so the bot can connect to Synapse via HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 2. Copy the binary we just built
COPY --from=builder /app/postmoogle /bin/postmoogle

# 3. Set the default DB path (Railway volume should be at /data)
ENV POSTMOOGLE_DB_DSN=/data/postmoogle.db

# 4. Run as root (UID 0) to ensure we have permission to write to the volume
USER 0

# 5. Start the bot
ENTRYPOINT ["/bin/postmoogle"]