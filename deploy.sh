#!/bin/sh
echo "#### Running script to deploy to Minikube..."

APP=taxianalytics
GCLOUD_KEY_FILE=key.json

echo "#### Loading env variables..."
if [ -z "$TAXI_TOPIC" ]
then
  echo "env TAXI_TOPIC must be set!"
  exit 1
else
  echo "using env TAXI_TOPIC=$TAXI_TOPIC"
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

echo "#### Checking code version..."
VERSION=$(git rev-parse --short HEAD)
echo "Current version: $VERSION"
export DOCKER_IMAGE_TAG=local/$APP:$VERSION

echo "#### Finding gcloud key file..."
if [ ! -f $GCLOUD_KEY_FILE ]; then
    echo "$GCLOUD_KEY_FILE not found!"
    exit 1
fi
echo "#### Loading gcloud key file..."
export GCLOUD_KEY=$(cat $GCLOUD_KEY_FILE)


echo "#### Switching to docker daemon in Minikube..."
eval $(minikube docker-env)

echo "#### Building image..."
docker build -t $DOCKER_IMAGE_TAG .

echo "#### Updating deployment in Minikube..."
cat kubernetes/$APP-deployment.yaml | ./envsubst | kubectl apply -f -

echo "#### Updating service in Minikube..."
export EXTERNAL_IP=$(minikube ip)
cat kubernetes/$APP-service.yaml | ./envsubst | kubectl apply -f -

echo "#### Getting app url..."
URL=$(minikube service $APP --url)
echo "$APP deployed at $URL"

eval $(minikube docker-env -u)
