#!/bin/bash

REGISTRY="registry.xuelangyun.com"
NAMESPACE="shuzhi-amd64"
IMAGE_NAME="suanpan-go-pipeline"
IMAGE_URL="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}"
IMAGE_VERSION="preview"
VERSION="$1"

docker build --pull \
-t ${IMAGE_URL}:${IMAGE_VERSION}-${VERSION} \
-t ${IMAGE_URL}:${IMAGE_VERSION} \
. \
-f ./docker/Dockerfile

docker push ${IMAGE_URL}:${IMAGE_VERSION}-${VERSION}
docker push ${IMAGE_URL}:${IMAGE_VERSION}
