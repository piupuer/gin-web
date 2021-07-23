#!/bin/bash

if [ "$REMOTE_DEBUG" == "true" ]; then
  if [ "$REMOTE_DEBUG_MULTI" == "true" ]; then
    ./dlv --listen=:${REMOTE_DEBUG_PORT} --headless=true --api-version=2 --log --accept-multiclient exec ./main-stage
  else
    ./dlv --listen=:${REMOTE_DEBUG_PORT} --headless=true --api-version=2 --log exec ./main-stage
  fi
else
  ./main-stage
fi
