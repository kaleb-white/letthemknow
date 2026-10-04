#!/bin/bash

# Description: Runs the message server tests in a container.

VERBOSE=1

log_always() { echo "[INFO] $1" >&1; }
log() { if [[ VERBOSE -eq 1 ]]; then echo "[INFO] $1" >&1; fi; }
warn() { if [[ VERBOSE -eq 1 ]]; then echo "[WARN] $1" >&1; fi; }
error() { if [[ VERBOSE -eq 1 ]]; then echo "[ERROR] $1" >&1; exit 1; fi; }

TAG=letthemknow/message-server-tests:latest

log "Starting build..."
docker build . -f docker/message-server.Dockerfile --target test -t "$TAG"

if [[ $? -ne 0 ]]; then
	error "Build failed, exiting."
else
	log "Build succeeded, continuing."
fi

log "Running tests..."
docker run --rm -it "$TAG" 
