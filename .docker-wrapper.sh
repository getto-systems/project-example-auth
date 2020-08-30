DOCKER_WRAPPER_IMAGE_golang(){ head -1 $APP_ROOT/Dockerfile | sed "s/.*golang:\(.*\) as.*/\1/"; }
DOCKER_WRAPPER_IMAGE_node(){ echo "14.8.0-buster-slim"; }
DOCKER_WRAPPER_IMAGE_google/cloud-sdk(){ echo "307.0.0-slim"; }

DOCKER_WRAPPER_SERVER_OPTS_http(){ echo "-p ${LABO_PORT_PREFIX}${APP_PORT}:8080"; }

DOCKER_WRAPPER_APP(){
  rm -rf $APP_ROOT/y_static/protocol_buffers && protoc -I=$APP_ROOT/auth-protocol_buffers --go_out=y_static $APP_ROOT/auth-protocol_buffers/*.proto
  go-server http $1 go run $(head -1 go.mod | cut -d' ' -f2)/x_http_server
}
DOCKER_WRAPPER_LOGS(){
  go-server http logs
}
