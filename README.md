# Numer

A microservices backend to support an invoice management system.

# Architecture: 
Uses gRPC for inter-service communication while providing a RESTful API gateway for external clients to consume.

1. Internal Microservices (Communicates with other services using gRPC):
    - Invoice Service: 
        - Handles all invoice-related operations.
    - User Service: 
        - Manages user data and authentication.
    - Activity Log Service: 
        - Logs all actions related to invoices and user activities.
    - Notification Service: 
        - Manages notifications like invoice reminders or payment confirmations.

2. API Gateway (RESTful):
    - Exposes a REST API for the web browser and other external clients.
    - Acts as a bridge between the frontend and the internal microservices.
    - Translates REST calls from external clients into gRPC calls to the respective microservices.
