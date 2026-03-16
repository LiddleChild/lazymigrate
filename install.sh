#!/usr/bin/env sh

APP_NAME="lazymigrate"

BRANCH=$1
BRANCH="${BRANCH:=main}"

# check for git
if ! command -v git &> /dev/null; then
  echo "git is required to install $APP_NAME"
  exit 1
fi

# check for go
if ! command -v go &> /dev/null; then
  echo "go1.26.0 is required to install $APP_NAME"
  exit 1
fi

# clone
git clone https://github.com/LiddleChild/$APP_NAME.git --depth 1 --branch $BRANCH /tmp/$APP_NAME

# build and put binary in path
cd /tmp/$APP_NAME; go install

# clean up
rm -rf /tmp/$APP_NAME

echo "$($APP_NAME -version) is ready to go"
