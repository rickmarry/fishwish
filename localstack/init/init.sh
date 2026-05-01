#!/bin/bash
set -e

echo "Creating S3 bucket..."
awslocal s3 mb s3://fishwish-photos

echo "Creating DynamoDB tables..."
awslocal dynamodb create-table \
  --table-name catch_logs \
  --attribute-definitions AttributeName=user_id,AttributeType=S AttributeName=created_at,AttributeType=S \
  --key-schema AttributeName=user_id,KeyType=HASH AttributeName=created_at,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1

awslocal dynamodb create-table \
  --table-name user_sessions \
  --attribute-definitions AttributeName=user_id,AttributeType=S \
  --key-schema AttributeName=user_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1

echo "LocalStack initialization complete."
