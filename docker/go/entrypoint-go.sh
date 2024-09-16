#!/bin/bash

USER_ID=${LOCAL_UID:-9001}
GROUP_ID=${LOCAL_GID:-9001}
USER_NAME=gopher

useradd -u $USER_ID -o -m $USER_NAME
groupmod -g $GROUP_ID $USER_NAME
export HOME=/home/$USER_NAME

# add permission of go directories to $USER_NAME
chown -R $USER_NAME:$USER_NAME /go/

echo "Started with UID : $USER_ID, GID: $GROUP_ID" && /bin/bash
