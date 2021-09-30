#!/bin/bash

if [ "$REMOTE_DEBUG_BIN" == "" ]; then
  REMOTE_DEBUG_BIN="main-stage"
fi

if [ "$REMOTE_DEBUG" == "true" ]; then
  if [ "$REMOTE_DEBUG_MULTI" == "true" ]; then
    ./dlv --listen=:$REMOTE_DEBUG_PORT --headless=true --api-version=2 --log --accept-multiclient exec ./$REMOTE_DEBUG_BIN
  else
    ./dlv --listen=:$REMOTE_DEBUG_PORT --headless=true --api-version=2 --log exec ./$REMOTE_DEBUG_BIN
  fi
else
  ./$REMOTE_DEBUG_BIN
fi
