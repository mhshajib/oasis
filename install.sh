#! /bin/bash

exec_curl(){
  echo "Found oasis latest version: $VERSION"
  echo "Download may take few minutes depending on your internet speed"
  echo "Downloading oasis to $2"
  curl -L --silent --connect-timeout 30 --retry-delay 5 --retry 5 -o "$2" "$1"
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
    echo "$MESSAGE_START"
    chmod +x $TARGET
    echo "$MESSAGE_END"
    oasis
  fi
  
  if [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
    exec_curl $URL/releases/download/$VERSION/darwin_arm64 $TARGET
    echo "$MESSAGE_START"
    chmod +x $TARGET
    echo "$MESSAGE_END"
    oasis
  fi

elif [ "$OS" == "Linux" ]; then
  if [ "$ARCH" == "x86_64" ]; then
    exec_curl $URL/releases/download/$VERSION/linux_amd64 $TARGET
    echo "$MESSAGE_START"
    chmod +x $TARGET
    echo "$MESSAGE_END"
    oasis
  fi

  if [ "$ARCH" == "arm64" ] || [ "$ARCH" == "aarch64" ]; then
    exec_curl $URL/releases/download/$VERSION/linux_arm64 $TARGET
    echo "$MESSAGE_START"
    chmod +x $TARGET
    echo "$MESSAGE_END"
    oasis
  fi

  if [ "$ARCH" == "i368" ]; then
    exec_curl $URL/releases/download/$VERSION/linux_386 $TARGET
    chmod +x $TARGET
    oasis
  fi
fi