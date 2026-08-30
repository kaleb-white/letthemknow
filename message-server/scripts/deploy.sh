#!/bin/bash

# Description: 
# Usage: 

VERBOSE=0

log() { if [ VERBOSE ]; then echo "[INFO] $1" >&1; fi; }
warn() { if [ VERBOSE ]; then echo "[WARN] $1" >&1; fi; }
error() { if [ VERBOSE ]; then echo "[ERROR] $1" >&1; exit 1; fi; }

# 1. 
