#!/usr/bin/env bash
set -e
echo "Building container image"
[[ -z "${CIRCLE_TAG}" ]] && tag="$(echo $CIRCLE_SHA1 | cut -c -7)" || tag="${CIRCLE_TAG}"
echo "Computed tag: $tag"
docker build --no-cache -t quay.io/solarwinds/gitlic-check .
docker tag quay.io/solarwinds/gitlic-check quay.io/solarwinds/gitlic-check:$tag
docker build --no-cache -t quay.io/solarwinds/augit-server:$tag -f Dockerfile_augit .
echo "Login to quay"
docker login -u $DOCKER_USER -p $DOCKER_PASS quay.io
echo "Login succeeded. Pushing images"
docker push quay.io/solarwinds/gitlic-check
docker push quay.io/solarwinds/gitlic-check:$tag
docker push quay.io/solarwinds/augit-server:$tag
