# EVNTX – Event Management Platform

EVNTX is a scalable, multi-role event management platform that enables users to discover events, book tickets, manage wallets, and allows organizers and admins to manage events and platform operations

---

## Features

### User (Goer)

* Discover and explore events
* View event details and ticket types
* Reserve and book tickets
* Secure payment system integration (Razorpay) for ticket purchases
* Manage bookings and generated tickets (PDF / QR Code)
* Wallet system for easy refunds and seamless checkouts
* Calendar view for upcoming events
* Passwordless authentication (Email OTP & Google OAuth)

### Organizer

* Create and manage events
* Submit events for approval
* Track revenue and engagement analytics
* Generate and manage tickets for attendees
* Organizer wallet for receiving direct earnings
* Request payouts from wallet balance

### Admin

* Manage users and organizers
* Approve/reject events and organizers
* Configure platform settings
* Audit logs and reporting

---

## System Architecture

### Backend – Clean Architecture

The backend strictly follows **Clean Architecture**, ensuring separation between business logic and frameworks.

#### Layers:

* **Domain** → Core business entities
* **Usecase** → Application logic
* **Repository (Interface)** → Contracts
* **Infrastructure** → DB, email, external services
* **Delivery (HTTP)** → Gin handlers

#### Supporting Packages:

* JWT handling
* OTP system
* OAuth integration
* Structured logging (Zerolog)
* Background workers

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
│   ├── server/      # Application entrypoint
│   ├── admin/       # Admin utilities
│   ├── seeder/      # DB seeding scripts
│   └── wallet_seeder/ # Wallet initializer
│
├── internal/
│   ├── domain/      # Core entities
│   ├── usecase/     # Business logic
│   ├── repository/  # Interfaces
│   ├── infrastructure/
│   │   ├── database/
│   │   ├── email/
│   │   ├── payment/   # Payment gateway integration
│   │   └── repository/  # Implementations
│   ├── delivery/
│   │   └── http/    # Gin handlers
│   └── middleware/  # Auth, logging, etc.
│
├── pkg/
│   ├── errors/
│   ├── hash/
│   ├── jwt/
│   ├── logger/
│   ├── oauth/
│   ├── otp/
│   ├── response/
│   └── workers/
│
└── assets/          # Static files (images, event media)
```

---

### Frontend

```id="frontend-structure"
frontend/
├── src/
│   ├── app/             # App-level setup (routing, providers)
│   ├── modules/         # Feature-based modules
│   │   ├── admin/
│   │   ├── auth/
│   │   ├── events/
│   │   ├── home/
│   │   ├── notifications/
│   │   ├── organizer/
│   │   ├── payments/
│   │   └── user/
│   │
│   ├── services/        # API layer (Axios)
│   └── shared/
│       ├── components/
│       ├── hooks/       # TanStack Query hooks
│       ├── ui/
│       └── utils/
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

* PostgreSQL relational schema
* UUID primary keys
* State-driven lifecycle modeling
* Strong referential integrity
* Ledger-based wallet system

Full Database Documentation: https://docs.google.com/document/d/1nFD9qYgVHtE5nrjjdYi2lhSgqYmTTEnDpu5OhtahSA0/edit?usp=sharing

---

## Lifecycle Models

### Event

```id="event-life"
draft → pending → approved → live → completed
draft → pending → rejected
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

* Concurrency-safe ticket inventory
* Optimistic locking for ticket updates
* Append-only wallet transaction ledger
* Structured logging with Zerolog
* Modular feature-based frontend architecture

---

## Author

**Aswin Sreeraj**

---

## License

This project is intended for educational and development purposes.
