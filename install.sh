#!/bin/bash

# Function to show a loading animation
loading_animation() {
  local pid=$!
  local delay=0.2
  local spinstr='|/-\'
  local i=0
  echo -n 'Downloading'
  
  while ps -p $pid > /dev/null; do
    i=$(( (i+1) %4 ))
    printf "\r%s %s" 'Downloading' "${spinstr:$i:1}"
    sleep $delay
  done
  printf "\r%s %s\n" 'Downloading' 'done'
}

exec_curl(){
  echo "Found oasis latest version: $VERSION"
  echo "Download may take a few minutes depending on your internet speed"
  echo "Download URL: $1"  # Added to debug the download URL

  # Start the download in the background
  curl -L --silent --connect-timeout 30 --retry-delay 5 --retry 5 -o "$2" "$1" &
  loading_animation

  if [ $? -ne 0 ]; then
    echo "Error: Download failed. Please check your internet connection or the URL and try again."
    exit 1
  fi

  if [ ! -f "$2" ]; then
    echo "Error: Failed to download oasis to $2. File not found."
    exit 1
  fi
}

OS=`uname`
ARCH=`uname -m`
VERSION=$1
URL=https://github.com/mhshajib/oasis
TARGET=/usr/local/bin/oasis
MESSAGE_START="Installing oasis"
MESSAGE_END="Installation complete"

if [ "$VERSION" == "" ]; then
  LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' $URL/releases/latest)
  VERSION=$(echo $LATEST_RELEASE | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
fi

if [ "$OS" == "Darwin" ]; then
  if [ "$ARCH" == "x86_64" ]; then
    exec_curl $URL/releases/download/$VERSION/darwin_amd64 $TARGET
  elif [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
    exec_curl $URL/releases/download/$VERSION/darwin_arm64 $TARGET
  fi
elif [ "$OS" == "Linux" ]; then
  if [ "$ARCH" == "x86_64" ]; then
    exec_curl $URL/releases/download/$VERSION/linux_amd64 $TARGET
  elif [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
    exec_curl $URL/releases/download/$VERSION/linux_arm64 $TARGET
  elif [ "$ARCH" == "i368" ]; then
    exec_curl $URL/releases/download/$VERSION/linux_386 $TARGET
  fi
fi

if [ -f "$TARGET" ]; then
  echo "$MESSAGE_START"
  chmod +x $TARGET
  echo "$MESSAGE_END"
  $TARGET
else
  echo "Error: oasis was not downloaded or the file path is incorrect."
  exit 1
fi
