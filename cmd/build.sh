#!/bin/bash

# Source environment variables
source ~/dev/lookinlabs/terraform/platform-tf/.env.production 

# Define the image name and tag
IMAGE_NAME="public.ecr.aws/c2w5h6c4/go-logger-middleware"
IMAGE_TAG="latest"

# Build the Docker image without using the cache
echo "Building Docker image without cache..."
docker build --no-cache -t ${IMAGE_NAME}:${IMAGE_TAG} .
if [ $? -ne 0 ]; then
    echo "Failed to build Docker image"
    exit 1
fi

# Reauthenticate with AWS ECR
echo "Authenticating with AWS ECR..."
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/c2w5h6c4
if [ $? -ne 0 ]; then
    echo "Failed to authenticate with AWS ECR"
    exit 1
fi

# Push the Docker image
echo "Pushing Docker image..."
docker push ${IMAGE_NAME}:${IMAGE_TAG}
if [ $? -ne 0 ]; then
    echo "Failed to push Docker image"
    exit 1
fi

echo "Docker image pushed successfully"