````mermaid
graph TD
subgraph Клиентская_часть
Client[Клиент <br> React + JavaScript]
end

subgraph Серверная_часть
APIGateway[API Gateway] -->|REST| CoreService[Core Service <br> Go]
APIGateway -->|RabbitMQ| Service[Notification Service <br> Go]
APIGateway -->|RabbitMQ| CommunicationService[Communication Service <br> Go]
CoreService -->|RabbitMQ| Service[Notification Service <br> Go]

CoreService -->|SQL| Database[База данных <br> PostgreSQL]
CommunicationService -->|SQL| Database
end

Client -->|HTTP/HTTPS| APIGateway

subgraph Внешние_службы
Service -->|SMTP/Push| External[Внешние службы <br> Email, Push]
end
````