#!/bin/bash

REGISTRY="registry.xuelangyun.com"
NAMESPACE="shuzhi-amd64"
IMAGE_NAME="suanpan-go-pipeline"
IMAGE_URL="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}"
IMAGE_VERSION="$1"

docker build --pull \
-t ${IMAGE_URL}:${IMAGE_VERSION} \
-t ${IMAGE_URL}:latest \
. \
-f ./docker/Dockerfile

docker push ${IMAGE_URL}:${IMAGE_VERSION}
docker push ${IMAGE_URL}:latest
