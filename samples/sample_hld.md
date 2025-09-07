# Notification System

## Problem Overview
Our current system sends user notifications via email only. This leads to low engagement because users increasingly prefer push notifications or in-app alerts.  
We need a scalable, multi-channel notification service that supports email, push, and SMS while maintaining low latency and reliability.

## Background
- The product team wants to send real-time alerts (e.g., fraud alerts, order updates).
- Current architecture couples notification logic inside the monolith API server.
- Spikes in notifications (e.g., during sales events) cause API slowdowns.
- We need to decouple notifications into a dedicated service that can scale independently.

## Options Considered
1. Continue with current monolith
    - ✅ No new infra
    - ❌ Risky: outages in monolith impact all notifications
    - ❌ Hard to add channels (push, SMS)

2. Use 3rd-party notification provider (e.g., Twilio, Firebase)
    - ✅ Easy to add channels, proven infra
    - ❌ Cost grows with scale
    - ❌ Less control over retries, routing logic

3. Build dedicated in-house notification service
    - ✅ Full control over retries, prioritization, future features (quiet hours, user preferences)
    - ✅ Independent scaling
    - ❌ Requires infra investment & ongoing maintenance

## High-Level Design (Chosen Approach: #3)
Pipeline:
API Gateway → Notification Service → Message Queue → Channel Workers

- API Gateway → Notification Service  
  - Accepts requests from product services (order system, fraud system).  
  - Each request specifies channel (email/push/SMS), user ID, and message payload.

- Message Queue (Kafka or Pub/Sub)  
  - Ensures decoupling and durability.  
  - Handles retry and backpressure.

- Channel Workers  
  - Independent workers per channel (email, push, SMS).  
  - Scale horizontally based on queue load.

- User Preference Store  
  - Database (Postgres) mapping users to allowed channels.  
  - Queried before sending notifications.

## Trade-offs
- ✅ Scales independently from monolith.
- ✅ Flexibility for new channels in future.
- ❌ Requires new infra (queues, DB).
- ❌ Higher operational burden compared to SaaS provider.