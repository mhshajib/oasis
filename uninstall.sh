#! /bin/bash

TARGET=/usr/local/bin/oasis
MESSAGE_START="Removing oasis"
MESSAGE_END="oasis removed"

echo "$MESSAGE_START"
rm $TARGET
echo "$MESSAGE_END"