#!/usr/bin/env bash
set -e
deploy="false"
echo "Branch name: $CIRCLE_BRANCH"
case $CIRCLE_BRANCH in
    "master")
    echo $MASTER_CONFIG | base64 -d > kubeconfig
    deploy="true"
    ;;
    "staging")
    echo $STAGING_CONFIG | base64 -d > kubeconfig
    deploy="true"
    ;;
esac
if [[ $deploy = "true" ]]; then
    echo "Proceeding with deployment"
    export KUBECONFIG=kubeconfig
    kubectl -n solarwindsio set image cronjob gitlic-check-cron gitlic-check=quay.io/solarwinds/gitlic-check:$tag
    kubectl -n solarwindsio set image cronjob augit-gh-report augit-gh-report=quay.io/solarwinds/augit-server:$tag
    kubectl -n solarwindsio set image cronjob augit-populator augit-populator=quay.io/solarwinds/augit-server:$tag
    kubectl -n solarwindsio set image deployment augit augit-server=quay.io/solarwinds/augit-server:$tag
    sleep 5
    response=`kubectl -n solarwindsio rollout status deployments/augit --watch=true`
    if [[ $response = *"error"* ]]; then
        echo "Deployment not successful with msg: '$response'. Rolling back. . . "
        kubectl rollout undo deployments/augit
        echo "Rolling back done . . . "
        exit 1
    fi
fi
echo "All done"
