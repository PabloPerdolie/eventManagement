#!/bin/sh

echo "Waiting for RabbitMQ to be ready..."
sleep 10

echo "Creating 'notifications' queue..."
rabbitmqadmin declare queue name=notifications durable=true

echo "RabbitMQ initialized successfully"
