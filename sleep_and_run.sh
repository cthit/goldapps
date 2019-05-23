#!/bin/sh
echo Sleeping for "$WAIT" before running \"goldapps "$*"\"
sleep "$WAIT" && /app/goldapps $*
