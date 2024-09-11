# Numer

A microservices backend to support an invoice management system. Uses gRPC for inter-service communication while providing a RESTful API gateway for external clients to consume.

## Overview: 

### Go, PostgresQL, RabbitMQ, go-cron

1. Internal Microservices (Communicates with other services using gRPC):
    - Invoice Service: 
        - Handles all invoice-related operations.
    - User Service: 
        - Manages user and customer data.
    - Activity Log Service: 
        - Logs all actions related to invoices and user activities.
    - Notification Service: 
        - Manages notifications like invoice reminders or payment confirmations.
    - Stats Service:
        - Aggregate and summarize key metrics related to invoices (e.g., total paid, overdue, draft, and unpaid invoices).
    - Reminders Service:
        - Handles reminders and the scheduling of reminders (e.g., 7 days or 3 days before due date).

2. API Gateway (RESTful):
    - Exposes a REST API for the web browser and other external clients.
    - Acts as a bridge between the frontend and the internal microservices.
    - Translates REST calls from external clients into gRPC calls to the respective microservices.

# Running the Application Locally

Follow the steps below to run the application locally.

## Prerequisites

- Ensure you have Docker and `make` installed on your machine.
- Clone this repository.

## Instructions

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone <repository-url>
cd <repository-directory>
```

### 2. Ensure .env files exist at the root of each service

### 3. Build the Docker Images
```bash
make build
```

### 4. Start the services
```bash
make up
```

### 5. Run tests
```bash
make test
```

### 6. Stop the services
```bash
make down
```

