#!/usr/bin/env bash
echo "Building container image"
docker build --no-cache -t quay.io/solarwinds/gitlic-check .
[[ -z "${CIRCLE_TAG}" ]] && tag="$(echo $CIRCLE_SHA1 | cut -c -7)" || tag="${CIRCLE_TAG}"
echo "Computed tag: $tag"
docker tag quay.io/solarwinds/gitlic-check quay.io/solarwinds/gitlic-check:$tag
echo "Login to quay"
docker login -u $DOCKER_USER -p $DOCKER_PASS quay.io
echo "Login succeeded. Pushing images"
docker push quay.io/solarwinds/gitlic-check
docker push quay.io/solarwinds/gitlic-check:$tag
echo "All done"