#!/bin/sh
DIR=$HOME/.gitdo
mkdir $DIR
echo "Copying plugins..."
echo "Copying hooks..."
cp -r ./plugins $DIR
cp -r ./hooks $DIR

PROCESSOR=`uname -p`
if [[ $PROCESSOR = "amd64" ]]
GOOS=`uname`
VERSIONTOCP="gitdo_${GOOS}_${GOARCH}"
echo "Copying $VERSIONTOCP to your /usr/local/bin/ ..."
cp $VERSIONTOCP /usr/local/bin/gitdo
echo "Done."
