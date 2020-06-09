bump_build

version=$(cat .release-version)

git clone https://getto-systems@gitlab.com/getto-systems-base/labo/project-example/k8s.git
cd k8s

bump_sync id/deployment.yaml 's|\(asia.gcr.io/getto-projects/example/id\):.*|\1:'$version'|'

curl $TRELLIS_CI_BUMP_VERSION/request.sh | bash -s -- $HOME/.message/update-id-image.sh

cd -
