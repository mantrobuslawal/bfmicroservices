// Package grpcadapter contains the gRPC transport adapter for the bfstore
// Catalog Service.
//
// This package translates between generated Protobuf API types and the catalog
// domain/service layer. It is responsible for receiving gRPC requests, performing
// transport-level validation, calling the catalog application service, mapping
// domain results into Protobuf responses, and converting domain errors into
// gRPC status errors.
//
// It's the intention of this package to stay thin. Business rules live in the catalog package,
// generated Protobuf definitions belong to the API contract, and persistence concerns exist in
// repository implementations. grpcadaper only coordinates these boundaries.
//
// Typical repsonsibilities include:
//
//   - registering the CatalogServiceServer with a gRPC server;
//   - validating request shape, such as nil requests, missing required IDs,
//     invalid page sizes, and malformed page tokens;
//   - mapping Protobuf messages, enums, money values, products, categories,
//     variants, attributes, images, and pagination metadata to and from domain
//     types;
//   - translating catalog domain errors into appropriate gRPC status codes;
//   - keeping transport-specific concerns out of the catalog domain model.
//
// The package does not conatin SQL, database transactions, Kafka publishing,
// catalog business rules, or direct infrastructure orchestration. Those
// responsibilities belong to lower-level repositories, domain/application
// services, or platform packages.
package grpcadapter
