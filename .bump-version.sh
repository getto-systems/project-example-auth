bump_build

version=$(cat .release-version)

git clone https://gitlab.com/getto-systems-base/labo/project-example/k8s.git
cd k8s

bump_sync id/deployment.yaml 's|\(asia.gcr.io/getto-projects/example/id\):.*|\1:'$version'|'

git commit -m "bump version: id"
git post "bump version: id"

cd -
