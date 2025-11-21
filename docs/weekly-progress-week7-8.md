# Weekly Progress Report
This section contains weekly project report for our Event based Ecommerce platform. for Distributed Systems Course.

**Team:** Group 3  
**Milestone:** Phase 1: System Design 

**Nov 21**  

**Submitted by (Lead):** Yohannes Tigistu

## Summary of Progress
We have successfully concluded the foundational design phase for the distributed e-commerce platform. Our primary focus during this sprint was translating the project requirements into a robust microservices architecture using a **Django (Backend)** and **Angular (Frontend)** technology stack.

### Work Completed
*   **Finalized System Architecture**
    *   Created component and deployment diagrams.
    *   Decided to utilize **Nginx** as our API Gateway to route traffic and manage load balancing.
    *   Implemented a **"database-per-service"** pattern using **PostgreSQL** to ensure loose coupling between the Order and Product services.
*   **Established Clear API Contracts**
    *   Drafted the **OpenAPI (Swagger)** specifications for the core services.
    *   Ensured that the interface between the Angular client and the Django backend is strictly defined before writing code.
    *   Defined standardized response formats for success and error states.
*   **Mapped Message Flows**
    *   Designed a sequence diagram detailing the order placement lifecycle.
    *   Adopted a hybrid communication strategy:
        *   **Synchronous HTTP calls** for immediate inventory checks.
        *   **Asynchronous messaging (via a Broker)** for payment processing and notifications.
    *   Prioritized data consistency and system resilience.
    *   Defined the event schemas (JSON) for the Pub/Sub topics to be used by the Message Broker.

## Pending Risks / Blockers
*   **Docker Orchestration:** Configuring the Docker Compose environment to reliably network multiple Django containers with the Message Broker (RabbitMQ/Redis) and Celery workers may present configuration challenges.
*   **Eventual Consistency:** Handling edge cases where the asynchronous payment fails after an order is created requires careful implementation of compensation logic (e.g., rolling back stock), which is complex to test.

## TA Feedback
*(To be filled by Instructor)*

## Action Items / Recommendations
*(To be filled by Instructor)*
