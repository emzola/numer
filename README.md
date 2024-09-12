# Numer

A microservices backend to support an invoice management system. Uses gRPC for inter-service communication while providing a RESTful API gateway for external clients to consume.

## Overview: 

### Go, PostgresQL, RabbitMQ, go-cron, Docker

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
## RESTful API Endpoints

### Invoices

- **Get all invoices**
  - `GET /invoices`
  - Description: Retrieve a list of all invoices for authenticated user.

- **Create a new invoice**
  - `POST /invoices`
  - Description: Create a new invoice.

- **Get a specific invoice by ID**
  - `GET /invoices/{id}`
  - Description: Retrieve a single invoice by its ID.

- **Update a specific invoice by ID**
  - `PATCH /invoices/{id}`
  - Description: Update an existing invoice by its ID.

- **Send a specific invoice**
  - `POST /invoices/send`
  - Description: Send an invoice.

### Stats

- **Get dashboard stats**
  - `GET /stats`
  - Description: Retrieve summary stats related to invoices for authenticated user.

### Activities

- **Get user activities**
  - `GET /activities/user`
  - Description: Retrieve activities related to authenticated user.

- **Get activities for a specific invoice by ID**
  - `GET /activities/invoice/{id}`
  - Description: Retrieve activities related to a specific invoice.

### Reminders

- **Create a new reminder**
  - `POST /reminders`
  - Description: Create a new reminder.

### Users

- **Create a new user**
  - `POST /users`
  - Description: Create a new user.

- **Get a specific user by ID**
  - `GET /users/{id}`
  - Description: Retrieve a single user by their ID.

- **Update a specific user by ID**
  - `PATCH /users/{id}`
  - Description: Update an existing user by their ID.

- **Delete a specific user by ID**
  - `DELETE /users/{id}`
  - Description: Delete a user by their ID.

### Customers

- **Create a new customer**
  - `POST /customers`
  - Description: Create a new customer.

- **Get a specific customer by ID**
  - `GET /customers/{id}`
  - Description: Retrieve a single customer by their ID.

- **Update a specific customer by ID**
  - `PATCH /customers/{id}`
  - Description: Update an existing customer by their ID.

- **Delete a specific customer by ID**
  - `DELETE /customers/{id}`
  - Description: Delete a customer by their ID.

