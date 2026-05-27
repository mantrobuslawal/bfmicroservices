# syntax=docker/dockerfile:1

# Catalogue Service container image.
#
# This Dockerfile is intentionally simple for the first implementation stage.
# It can be hardened later with:
# - non-root user
# - distroless runtime image
# - SBOM generation
# - vulnerability scanning
# - build provenance

FROM golang:1.22-bookworm AS build

WORKDIR /src

COPY services/catalog-service/go.mod services/catalog-service/go.sum* ./
RUN go mod download

COPY services/catalog-service/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/catalog-service ./cmd/catalog-service

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=build /out/catalog-service /catalog-service

EXPOSE 50051

ENTRYPOINT ["/catalog-service"]
