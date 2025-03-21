# Notification Service

This service is responsible for sending notifications to users via email based on events received from RabbitMQ.

## Features

- Receives messages from RabbitMQ queue "notifications"
- Sends email notifications via SMTP
- Provides health check and service information endpoints
- Supports multiple notification types:
  - Event Creation
  - Task Assignment
  - Expense Addition

## Configuration

The service can be configured using environment variables:

### Server
- `PORT`: HTTP server port (default: 8082)

### RabbitMQ
- `RABBITMQ_HOST`: RabbitMQ host (default: localhost)
- `RABBITMQ_PORT`: RabbitMQ port (default: 5672)
- `RABBITMQ_USERNAME`: RabbitMQ username (default: guest)
- `RABBITMQ_PASSWORD`: RabbitMQ password (default: guest)
- `RABBITMQ_QUEUE`: RabbitMQ queue name (default: notifications)

### SMTP
- `SMTP_HOST`: SMTP server host (default: smtp.example.com)
- `SMTP_PORT`: SMTP server port (default: 587)
- `SMTP_USERNAME`: SMTP username (default: notifications@system.com)
- `SMTP_PASSWORD`: SMTP password (default: password)
- `SMTP_SENDER`: SMTP sender email (default: notifications@system.com)

## API Endpoints

- `GET /api/v1/health`: Health check endpoint
- `GET /api/v1/info`: Service information endpoint
- `GET /health`: Simple health check endpoint

## Message Format

Messages received from RabbitMQ should be in JSON format with the following structure:

```json
{
  "event": "event_created",
  "data": {
    "event_id": 1,
    "title": "Team Meeting",
    "user_email": "user@example.com"
  }
}
```

## Running the Service

### Using Go directly

```bash
go run main.go
```

### Building the Service

```bash
go build -o notification-service
```

## Docker Compose Setup

The service can be run with Docker Compose, which includes:

- RabbitMQ with management interface
- MailHog for email testing
- Notification Service

### Starting the Docker Compose environment

```bash
docker-compose up -d
```

### Accessing Services

- RabbitMQ Management: http://localhost:15672 (username: guest, password: guest)
- MailHog Interface: http://localhost:8025
- Notification Service: http://localhost:8082

### Testing with RabbitMQ Management Interface

1. Open RabbitMQ Management at http://localhost:15672
2. Navigate to "Queues" tab
3. Click on the "notifications" queue
4. Go to the "Publish message" section
5. Use the sample message format (see test-message.json) and click "Publish Message"
6. Check MailHog at http://localhost:8025 to see the delivered email

### Stopping the environment

```bash
docker-compose down
```
