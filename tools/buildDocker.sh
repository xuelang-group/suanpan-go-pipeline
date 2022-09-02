#!/bin/bash

REGISTRY="registry.xuelangyun.com"
NAMESPACE="shuzhi-amd64"
IMAGE_NAME="suanpan-go-pipeline"
IMAGE_URL="${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}"
IMAGE_VERSION="$1"

if [ -n "${PREVIEW}" ]; then
    docker build --pull \
    -t ${IMAGE_URL}:${PREVIEW}-${VERSION} \
    -t ${IMAGE_URL}:${PREVIEW} \
    . \
    -f ./docker/Dockerfile

    docker push ${IMAGE_URL}:${PREVIEW}-${VERSION}
    docker push ${IMAGE_URL}:${PREVIEW}
else
    docker build --pull \
    -t ${IMAGE_URL}:${VERSION} \
    -t ${IMAGE_URL}:latest \
    . \
    -f ./docker/Dockerfile

    docker push ${IMAGE_URL}:${VERSION}
    docker push ${IMAGE_URL}:latest
fi
