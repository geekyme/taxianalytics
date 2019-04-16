#!/bin/sh
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "${GREEN}Running script to deploy to Minikube${NC}"


APP=taxianalytics
GCLOUD_KEY_FILE=key.json

echo "${GREEN}Loading env variables${NC}"
if [ -z "$TAXI_SUB_NAME" ]
then
  echo "env TAXI_SUB_NAME must be set!"
  exit 1
else
  echo "using env TAXI_SUB_NAME=$TAXI_SUB_NAME"
fi

if [ -z "$TAXI_PROJECT" ]
then
  echo "env TAXI_PROJECT must be set!"
  exit 1
else
  echo "using env TAXI_PROJECT=$TAXI_PROJECT"
fi

if [ -z "$DB_HOST" ]
then
  echo "env DB_HOST must be set!"
  exit 1
else
  echo "using env DB_HOST=$DB_HOST"
fi

echo "${GREEN}Checking code version${NC}"
VERSION=$(git rev-parse --short HEAD)
echo "Current version: $VERSION"
export DOCKER_IMAGE_TAG=local/$APP:$VERSION

echo "${GREEN}Finding gcloud key file${NC}"
if [ ! -f $GCLOUD_KEY_FILE ]; then
  echo "$GCLOUD_KEY_FILE not found!"
  exit 1
else
  echo "Found $GCLOUD_KEY_FILE"
fi

echo "${GREEN}Loading gcloud key file${NC}"
export GCLOUD_KEY=$(cat $GCLOUD_KEY_FILE)
echo "Loaded"

echo "${GREEN}Switching to docker daemon in Minikube${NC}"
eval $(minikube docker-env)
echo "Switched"

echo "${GREEN}Building image${NC}"
docker build -t $DOCKER_IMAGE_TAG .

echo "${GREEN}Updating deployment in Minikube${NC}"
cat kubernetes/$APP-deployment.yaml | ./envsubst | kubectl apply -f -

echo "${GREEN}Updating service in Minikube${NC}"
export EXTERNAL_IP=$(minikube ip)
cat kubernetes/$APP-service.yaml | ./envsubst | kubectl apply -f -

echo "${GREEN}Getting app url${NC}"
URL=$(minikube service $APP --url)
echo "$APP deployed at $URL"

eval $(minikube docker-env -u)
