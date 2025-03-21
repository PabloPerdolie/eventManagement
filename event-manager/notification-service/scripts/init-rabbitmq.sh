#!/bin/sh

# Скрипт для инициализации очереди в RabbitMQ

# Ждем, когда RabbitMQ станет доступен
echo "Waiting for RabbitMQ to be ready..."
sleep 10

# Создаем очередь
echo "Creating 'notifications' queue..."
rabbitmqadmin declare queue name=notifications durable=true

echo "RabbitMQ initialized successfully"
