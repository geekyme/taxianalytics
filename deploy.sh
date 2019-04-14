#!/bin/sh
echo "#### Running script to deploy to Minikube..."

APP=taxianalytics
VERSION=$(git rev-parse --short HEAD)
GCLOUD_KEY_FILE=key.json

if [ -z "$TAXI_TOPIC" ]
then
  echo "env TAXI_TOPIC must be set!"
  exit 1
fi

if [ -z "$TAXI_PROJECT" ]
then
  echo "env TAXI_PROJECT must be set!"
  exit 1
fi

echo "#### Current version: $VERSION"
export DOCKER_IMAGE_TAG=local/$APP:$VERSION

echo "#### Finding gcloud key file..."
if [ ! -f $GCLOUD_KEY_FILE ]; then
    echo "$GCLOUD_KEY_FILE not found!"
    exit 1
fi
echo "#### Loading gcloud key file..."
export GCLOUD_KEY=$(cat $GCLOUD_KEY_FILE)


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
