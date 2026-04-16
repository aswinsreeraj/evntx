# EVNTX – Event Management Platform

EVNTX is a scalable, multi-role event management platform that enables users to discover events, book tickets, manage wallets, and allows organizers and admins to manage events and platform operations

---

## Features

### User (Goer)

* Discover and explore events with categorized listings
* View event details, ticket types, and availability
* Reserve and book tickets with real-time inventory management
* Secure payment system integration (Razorpay) for ticket purchases
* Manage bookings and generated tickets (PDF / QR Code)
* Wallet system for easy refunds and seamless checkouts
* Calendar view for tracking upcoming events
* Passwordless authentication (Email OTP & Google OAuth)

### Organizer

* Create and manage events with rich media support
* Submit events for admin approval with state tracking
* **Engagement Analytics:** Track visitors, event views, ticket selections, and checkout starts
* **Sales Reports:** Detailed revenue and booking breakdowns
* Generate and manage tickets for attendees
* Organizer wallet for direct earnings and ledger-based transactions
* Request and track payout settlements

### Admin

* **Global Dashboard:** High-level platform metrics and engagement trends
* **Management Suites:** Comprehensive CRUD for Users, Organizers, and Events
* **Audit Logging:** Detailed trail of all administrative actions
* **Advanced Admin Management:** Secure admin creation and deletion with self-deletion protection
* **Platform Configuration:** Dynamic control over fees, payment providers, and system settings
* **Reporting:** Global revenue and engagement analytics with **CSV Export** capabilities

---

## System Architecture

### Backend – Clean Architecture

The backend strictly follows **Clean Architecture**, ensuring separation between business logic and frameworks.

#### Layers:

* **Domain** → Core business entities (Events, Bookings, Engagement, Audit)
* **Usecase** → Application logic and business rules
* **Repository (Interface)** → Contracts for data persistence
* **Infrastructure** → GORM implementations, email (SMTP), payments (Razorpay)
* **Delivery (HTTP)** → Gin-gonic handlers and routing

#### Core Systems:

* **Telemetry Pipeline:** Real-time tracking of visitor sessions and engagement events
* **Cron Scheduler:** Centralized job registry for expiring bookings, processing payouts, and event auto-completion
* **In-Memory Caching:** Local caching layer for high-demand endpoints (e.g., event details)
* **Auth System:** JWT-based session management with Refresh Tokens and RBAC
* **Audit System:** Automatic logging of sensitive administrative operations

---

## Tech Stack

### Backend

* **Language:** Go (Golang)
* **Framework:** Gin
* **ORM:** GORM
* **Database:** PostgreSQL
* **Authentication:** JWT + Refresh Tokens
* **Payments:** Razorpay API Integration
* **Hashing:** `crypto/sha256`
* **Logging:** Zerolog

### Frontend

* **Framework:** React (TypeScript)
* **State Management:** Zustand
* **Server State:** TanStack Query
* **HTTP Client:** Axios

---

## Project Structure

### Backend

```id="backend-structure"
backend/
├── cmd/
│   ├── server/          # Application entrypoint
│   ├── admin/           # Admin utilities
│   ├── seeder/          # DB seeding scripts
│   └── wallet_seeder/   # Wallet initializer
│
├── internal/
│   ├── cache/           # In-memory caching implementation
│   ├── domain/          # Core entities (engagement, audit, settings, etc.)
│   ├── usecase/         # Business logic layer
│   ├── repository/      # Repository interfaces
│   ├── infrastructure/
│   │   ├── database/    # DB connection management
│   │   ├── email/       # SMTP implementation
│   │   ├── payment/     # Razorpay integration
│   │   └── repository/  # GORM implementations
│   ├── delivery/
│   │   └── http/        # Gin handlers and API routes
│   └── middleware/      # Auth, Logging, Rate-limiting, RBAC
│
├── pkg/
│   ├── errors/          # Custom error types
│   ├── hash/            # Hashing utilities
│   ├── jwt/             # Token management
│   ├── logger/          # Zerolog configuration
│   ├── oauth/           # Google OAuth logic
│   ├── otp/             # OTP generation
│   ├── response/        # Standardize API responses
│   └── workers/         # Cron scheduler and background jobs
│
└── assets/              # Static files (images, event media)
```

---

### Frontend

```id="frontend-structure"
frontend/
├── src/
│   ├── app/             # App-level entry, routing, and providers
│   ├── modules/         # Feature-based modular architecture
│   │   ├── admin/       # Dashboard, User/Org management, Reports, Settings
│   │   ├── auth/        # Login, Register, OTP flows
│   │   ├── events/      # Event discovery and details
│   │   ├── home/        # Landing page
│   │   ├── notifications/ # Notification center
│   │   ├── organizer/   # Org dashboard, Event creation, Analytics
│   │   ├── payments/    # Razorpay checkout integration
│   │   └── user/        # User profile, Wallet, Bookings
│   │
│   ├── services/        # Centralized Axios API instances
│   └── shared/
│       ├── components/  # Layout, Navbar, Breadcrumbs
│       ├── hooks/       # Custom hooks and TanStack Query logic
│       ├── ui/          # Atomic UI components
│       └── utils/       # Helpers and formatters
│
└── assets/
```

---

## Authentication & Authorization

* Passwordless authentication
* Email OTP + Google OAuth
* JWT access tokens (short-lived)
* Refresh token session management
* Role-Based Access Control (RBAC)

Roles:

* `goer`
* `organizer`
* `admin`

---

## API Overview

Base URL:

```id="api-base"
/api/v1
```

Standard Response:

```json id="response-format"
{
  "success": true,
  "message": "Description",
  "data": {},
  "error": null
}
```

Full API Documentation: https://docs.google.com/document/d/1G3qYPxqshgR_f2C00_ebOxm6QKfmjdeXKKb1siq80aY/edit?usp=sharing

---

## Database Design

* PostgreSQL relational schema with curated indexing
* UUID v4 primary keys for enhanced security
* State-driven lifecycle modeling for Events and Bookings
* Ledger-based append-only wallet transaction system
* Global audit logging for administrative transparency
* Engagement telemetry tracking (Visitor sessions -> Daily aggregations)

Full Database Documentation: https://docs.google.com/document/d/1nFD9qYgVHtE5nrjjdYi2lhSgqYmTTEnDpu5OhtahSA0/edit?usp=sharing

---

## Lifecycle Models

### Event

```id="event-life"
draft → pending → approved → live → completed
draft → pending → rejected
live → cancellation_pending → cancelled
```

### Booking

```id="booking-life"
reserved → paid → cancelled
reserved → expired
```

### Payment

```id="payment-life"
initiated → success → refunded
initiated → failed
```

### Payout

```id="payout-life"
pending → approved → settled
pending → rejected
```

---

## Setup Instructions

### 1. Clone Repository

```bash id="clone"
git clone https://github.com/your-username/evntx.git
cd evntx
```

---

### 2. Backend Setup

```bash id="backend-run"
cd backend
go mod tidy
go run cmd/server/main.go
```

#### Environment Variables

```env id="env"
DB_HOST=hostname
DB_PORT=db_port
DB_USER=db_user
DB_PASSWORD=db_password
DB_NAME=db_name

JWT_SECRET=jwtsecretkey
JWT_REFRESH_SECRET=jwtrefreshsecret

GOOGLE_CLIENT_ID=googleclientid

ADMIN_EMAIL=adminemail@mail.com

RAZORPAY_KEY_ID=your_razorpay_key_id
RAZORPAY_KEY_SECRET=your_razorpay_key_secret
```

---

### 3. Frontend Setup

```bash id="frontend-run"
cd frontend
npm install
npm run dev
```

---

## Testing

* API testing via Postman / cURL
* Token validation and RBAC testing
* Edge case validation:

  * Concurrent ticket booking
  * Booking expiration
  * Wallet consistency

---

## Key Design Highlights

* **Concurrency-Safe Ticket Inventory:** Optimistic locking ensures zero double-booking.
* **Append-Only Ledger:** Financial integrity via immutable wallet transactions.
* **Engagement Pipeline:** Context-aware tracking of the user acquisition funnel.
* **Centralized Job Registry:** Reliable background processing with retry strategies.
* **In-Memory Caching:** Sub-millisecond response times for critical read paths.
* **Modular Frontend:** Domain-driven module structure for scalability.
* **Structured Logging:** Unified Zerolog implementation across all layers.

---

## Author

**Aswin Sreeraj**

---

## License

This project is intended for educational and development purposes.
