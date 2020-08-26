#!/bin/sh

deploy_main() {
  if [ ! -f "${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}" ]; then
    echo "key file : GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON is not exists"
    exit 1
  fi

  local host
  local region
  local project
  local image
  local version
  local tag
  local account

  host=asia.gcr.io
  region=asia-northeast1

  project=getto-projects
  image=example/auth
  version=$(cat .release-version)

  tag=${host}/${project}/${image}:${version}

  export HOME=$(pwd)

  account=$(cat "${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}" | grep "client_email" | cut -d'"' -f4)

  echo "${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}" | gcloud auth activate-service-account "$account" --key-file=-
  gcloud run deploy example-auth --image="$tag" --platform=managed --region="$region" --project="$project"
}

deploy_main
