#!/bin/sh

deploy_main() {
  if [ ! -f "${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}" ]; then
    echo "key file : GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON is not exists"
    exit 1
  fi
  if [ ! -f "${GOOGLE_CLOUD_SERVICE_ACCOUNT_NAME}" ]; then
    echo "key file : GOOGLE_CLOUD_SERVICE_ACCOUNT_NAME is not exists"
    exit 1
  fi

  export GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}

  local host
  local project
  local image
  local version
  local tag

  host=asia.gcr.io

  cat $GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON | docker login -u _json_key --password-stdin https://${host}

  project=getto-projects
  image=example/id
  version=$(cat .release-version)

  tag=${host}/${project}/${image}:${version}

  gcloud run deploy --image="$tag" --platform=managed --service-account="${GOOGLE_CLOUD_SERVICE_ACCOUNT_NAME}"

  docker build -t $tag . &&
  docker push $tag
}

deploy_main
