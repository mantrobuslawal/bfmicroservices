# Borough Online Furniture Store (bfstore)

**bfstore** is a cloud-native microservice backend for a fictional online furniture store operated by **ACME Ltd**.

The project is designed to demonstrate professional backend engineering, platform-aware application design, service boundaries, event-driven architecture, gRPC communication, protobuf contracts, Kafka messaging, service-owned databases, observability, testing, and deployment readiness.

This repository contains the application code, service contracts, application documentation, local development configuration, tests, and deployment manifests for the bfstore backend.

---

## Project Goals

The goal of this project is to build a realistic backend platform for an online furniture store while demonstrating production-style engineering practices.

bfstore is intended to show:

- Microservice architecture design
- Domain-led service boundaries
- gRPC-based internal service communication
- Protobuf-based API and event contracts
- Kafka-based asynchronous messaging
- Service-owned databases
- Structured logging, metrics, and tracing
- Contract, integration, end-to-end, and performance testing
- Containerized local development
- Kubernetes-ready deployment configuration
- Secure software delivery practices
- Clear documentation and architectural decision records

---

## Business Context

ACME Ltd operates an online furniture store where customers can browse furniture, view product details, manage a basket, place orders, reserve stock, make payments, and receive order notifications.

The backend is intentionally designed as a microservice system because the store contains several distinct business capabilities:

- Product catalogue management
- Inventory and stock reservation
- Basket management
- Customer management
- Order management
- Payment processing
- Notification delivery
- Authentication and authorization

Each capability is owned by a dedicated service with its own API, data model, and operational responsibilities.

---

## Architecture Overview

bfstore uses a hybrid communication model:

- **Synchronous communication** using gRPC
- **Asynchronous communication** using Kafka events
- **Protobuf** for service contracts and event payloads

The design follows this principle:

> Commands that need an immediate result use gRPC. Facts that have already happened are published as Kafka events.

Example:

```text
CreateOrder          -> gRPC command
ReserveStock         -> gRPC command
AuthorizePayment     -> gRPC command

OrderCreated         -> Kafka event
StockReserved        -> Kafka event
PaymentAuthorized    -> Kafka event
NotificationRequested -> Kafka event
```
## High-Level System Flow

The API Gateway exposes client-facing APIs and routes internal requests to backend services.

The backend services communicate through gRPC when they need a direct response, and publish Kafka events when other services need to react asynchronously.

```text
Customer / Frontend
        |
        v
API Gateway
        |
        +--> Auth Service
        +--> Customer Service
        +--> Catalog Service
        +--> Basket Service
        +--> Order Service
                  |
                  +--> Inventory Service
                  +--> Payment Service
                  |
                  v
                Kafka
                  |
                  +--> Notification Service
                  +--> Shipping Service
                  +--> Analytics Consumers
```
## Core Services

| Service | Responsibility |
|---|---|
| `api-gateway` | Public entry point for frontend clients |
| `auth-service` | Authentication, authorization, tokens, user sessions |
| `customer-service` | Customer profiles, addresses, preferences |
| `catalog-service` | Products, categories, furniture details, pricing |
| `inventory-service` | Stock levels, stock reservations, warehouse availability |
| `basket-service` | Customer basket and basket items |
| `order-service` | Order creation, order lifecycle, order history |
| `payment-service` | Payment authorization, capture, refunds |
| `shipping-service` | Delivery options, shipment creation, fulfilment status, tracking updates |
| `notification-service` | Email/SMS/event-driven customer notifications |
| `review-service` | Product reviews, ratings, moderation status |
| `search-service` | Product search, filtering, faceted search, search index updates |
| `recommendation-service` | Product recommendations, related items, personalised suggestions |

These services represent the target service landscape for bfstore. The initial implementation may focus on a smaller vertical slice first, such as catalogue, inventory, basket, order, payment, shipping, and notification, before expanding into reviews, search, and recommendations.

## Repository Layout

```text
bfstore/
├── README.md
├── Makefile
├── docker-compose.yml
├── buf.yaml
├── buf.gen.yaml
│
├── docs/
│   ├── requirements/
│   ├── architecture/
│   ├── api/
│   ├── events/
│   ├── data/
│   ├── testing/
│   ├── security/
│   ├── observability/
│   └── operations/
│
├── adr/
│   ├── README.md
│   ├── 0001-use-microservices.md
│   ├── 0002-use-grpc-for-service-communication.md
│   ├── 0003-use-kafka-for-events.md
│   └── 0004-use-service-owned-databases.md
│
├── proto/
│   └── acme/
│       ├── common/
│       ├── catalog/
│       ├── inventory/
│       ├── basket/
│       ├── order/
│       ├── payment/
│       ├── customer/
│       └── notification/
│
├── services/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── customer-service/
│   ├── catalog-service/
│   ├── inventory-service/
│   ├── basket-service/
│   ├── order-service/
│   ├── payment-service/
│   └── notification-service/
│
├── packages/
│   └── go/
│       ├── logger/
│       ├── config/
│       ├── grpc/
│       ├── kafka/
│       ├── telemetry/
│       ├── auth/
│       └── errors/
│
├── db/
│   ├── catalog/
│   ├── inventory/
│   ├── basket/
│   ├── order/
│   ├── payment/
│   ├── customer/
│   └── notification/
│
├── deploy/
│   ├── docker/
│   ├── kubernetes/
│   ├── helm/
│   └── kustomize/
│
├── tests/
│   ├── contract/
│   ├── integration/
│   ├── e2e/
│   ├── performance/
│   └── testdata/
│
├── tools/
├── scripts/
└── .github/
    ├── workflows/
    ├── CODEOWNERS
    ├── dependabot.yml
    └── pull_request_template.md
```
