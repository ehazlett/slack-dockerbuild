#!/bin/bash
URL=${URL:-}
CHANNEL=${CHANNEL:-#general}
/usr/local/bin/slack-dockerbuild -url $URL -channel $CHANNEL
