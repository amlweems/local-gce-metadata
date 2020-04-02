#!/bin/sh
echo -n True > ~/.config/gcloud/gce
export GCE_METADATA_ROOT=127.0.0.1
export CLOUDSDK_CONFIG=`mktemp -d`