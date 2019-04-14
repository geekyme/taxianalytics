#!/bin/sh
echo "#### Running script to deploy to Minikube..."

APP=taxianalytics
VERSION=$(git rev-parse --short HEAD)

echo "#### Current version: $VERSION"

export DOCKER_IMAGE_TAG=local/$APP:$VERSION

echo "#### Switching to docker daemon in Minikube"
eval $(minikube docker-env)

docker build -t $DOCKER_IMAGE_TAG .

echo "#### Updating deployment in Minikube..."

cat kubernetes/$APP-deployment.yaml | ./envsubst | kubectl apply -f -

kubectl expose deployment $APP --type=NodePort --name=$APP

URL=$(minikube service $APP --url)

echo "$APP deployed at $URL"

echo "#### Switching to back docker daemon in host"
eval $(minikube docker-env -u)
