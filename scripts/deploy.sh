#!/bin/bash

# Description: Deploys nginx, react spa, go server for development or production.

VERBOSE=0

log_always() { echo "[INFO] $1" >&1; }
log() { if [[ $VERBOSE -eq 1 ]]; then echo "[INFO] $1" >&1; fi; }
warn() { if [[ $VERBOSE -eq 1 ]]; then echo "[WARN] $1" >&1; fi; }
error() { echo "[ERROR] $1" >&1; exit 1; }

# Options
HELP_STRING=""" 
	-b / --build: Build <react|api|nginx>, comma-separated.
	-h / --help: Print help string
	-v / --verbose: Print operation details to stdout"""

BUILD=0
BUILD_ARGS=

while [[ $# -gt 0 ]]; do
  case $1 in
    -b|--build)
      BUILD=1
		  BUILD_ARGS="$2"
      shift 
      shift
      ;;
    -v|--verbose)
      VERBOSE=1
      shift
      ;;
    -h|--help)
			printf "$HELP_STRING"
			exit 0
      ;;
    -*|--*)
      echo "Unknown option $1"
      exit 1
      ;;
    *)
      shift # past argument
      ;;
  esac
done

docker info > /dev/null 2>&1 
if [[ $? -eq 1 ]]; then
	error "Docker daemon is not running. Please start daemon and rerun."
	exit 1
fi

rg -q 'letthemknow$' <<< "$PWD"
if [[ $? -eq 1 ]]; then
	error "Call this script from the repo root, not from $PWD."
	exit 1
fi

BUILD_TARGETS=()
if [[ $BUILD -eq 1 ]]; then
	log "Parsing build args $BUILD_ARGS..."
	BUILD_TARGETS=(${BUILD_ARGS//,/ })
	for BUILD_TARGET in "${BUILD_TARGETS[@]}"; do
		case $BUILD_TARGET in
		  react)
				log_always "Build react not yet configured."
				;;
			api)
				log "Building api..."
				docker build -f docker/message-server.Dockerfile . -t $BUILD_TARGET &
				;;
			nginx)
				log "Building nginx..."
				docker build -f docker/nginx.Dockerfile . -t $BUILD_TARGET &
				;;
			*)
				warn "Didn't recognize build target $BUILD_TARGET."
				;;
		esac
	done
fi
wait

log "Sourcing .env..."
set -a && source .env && set +a

RUN_TARGETS=("react" "nginx" "api")
for RUN_TARGET in "${RUN_TARGETS[@]}"; do
	case $RUN_TARGET in
		react)
			log_always "Run react not yet configured."
			;;
		api)
			log "Running api..."
			if docker stop $RUN_TARGET > /dev/null 2>&1; then
				log "Stopped $RUN_TARGET, deleting it then restarting..."
				docker rm $RUN_TARGET > /dev/null 2>&1
			else
				log "$RUN_TARGET not found, starting..."
			fi
			
			docker run -p ${API_PORT}:${API_PORT} -d --name $RUN_TARGET $RUN_TARGET > /dev/null
			;;
		nginx)
			log "Running nginx..."
				if docker stop $RUN_TARGET > /dev/null 2>&1; then
				log "Stopped $RUN_TARGET, deleting it then restarting..."
				docker rm $RUN_TARGET > /dev/null 2>&1
			else
				log "$RUN_TARGET not found, starting..."
			fi
			docker run -p ${NGINX_PORT}:${NGINX_PORT} -d --name $RUN_TARGET $RUN_TARGET > /dev/null
			;;
		*)
			warn "Didn't recognize build target $BUILD_TARGET."
			;;
	esac
done
