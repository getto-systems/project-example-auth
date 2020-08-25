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

  cat $GOOGLE_CLOUD_SERVICE_ACCOUNT_KEY_JSON | docker login -u _json_key --password-stdin https://${host}

  project=getto-projects
  image=example/id
  version=$(cat .release-version)

  tag=${host}/${project}/${image}:${version}

  docker build -t $tag . &&
  docker push $tag
}

deploy_main
