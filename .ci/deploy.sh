#!/bin/sh

deploy_main() {
  if [ ! -f "${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}" ]; then
    echo "key file : GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON is not exists"
    exit 1
  fi

  export GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON}

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

  echo gcloud run deploy --image="$tag"
}

deploy_main
