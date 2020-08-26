#!/bin/sh

deploy_main() {
  if [ ! -f "${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}" ]; then
    echo "key file : GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON is not exists"
    exit 1
  fi

  local host
  local project
  local image
  local version
  local tag

  host=asia.gcr.io

  project=getto-projects
  image=example/auth
  version=$(cat .release-version)

  tag=${host}/${project}/${image}:${version}

  echo "${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}" | gcloud auth activate-service-account --key-file=-
  gcloud run deploy example-auth --image="$tag" --platform=managed --region=asia-northeast1 --project=getto-projects
}

deploy_main
