#!/usr/bin/env bash
echo "Building container image"
[[ -z "${CIRCLE_TAG}" ]] && tag="$(echo $CIRCLE_SHA1 | cut -c -7)" || tag="${CIRCLE_TAG}"
echo "Computed tag: $tag"
docker build --no-cache -t quay.io/solarwinds/gitlic-check . && \
docker tag quay.io/solarwinds/gitlic-check quay.io/solarwinds/gitlic-check:$tag && \
echo "Login to quay" && \
docker login -u $DOCKER_USER -p $DOCKER_PASS quay.io && \
echo "Login succeeded. Pushing images" && \
docker push quay.io/solarwinds/gitlic-check && \
docker push quay.io/solarwinds/gitlic-check:$tag && \
echo $MAIN_CONFIG | base64 -d > kubeconfig && \
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
chmod +x kubectl && \
export KUBECONFIG=kubeconfig && \
./kubectl -n solarwindsio set image cronjob gitlic-check-cron gitlic-check=quay.io/solarwinds/gitlic-check:$tag && \
sleep 5 && \
response=`./kubectl -n solarwindsio rollout status cronjob/gitlic-check-cron --watch=true` && \
if [[ $response = *"error"* ]]; then
    echo "Deployment not successful with msg: '$response'. Rolling back. . . "
    ./kubectl rollout undo cronjob/gitlic-check-cron
    echo "Rolling back done . . . "
    exit 1
fi
echo "All done"