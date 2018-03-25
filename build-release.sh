#!/bin/sh
set -e

PLATFORMS="windows linux freebsd darwin dragonfly"
APPS="docproc.fileinput docproc.proc"
FILES="LICENSE README.md"
FOLDERS=examples

VERSION=`cat VERSION`
LDFLAGS="-X main.version=$VERSION"
ARCH=`go env GOARCH`

echo "Creating release packages for version $VERSION..."

rm -rf dist doc/_build vendor
mkdir dist

echo "Creating documentation..."
make -C doc html

dep ensure -v || $GOPATH/bin/dep ensure -v

for os in $PLATFORMS; do
    suffix=""
    if [ "$os" = "windows" ]; then
        suffix=".exe"
    fi
    distname=docproc-$VERSION-$os-$ARCH
    destdir=dist/$distname
    echo "Building release for $os in $destdir..."
    GOOS=$os
    GOARCH=$ARCH
    cp -rf doc/_build/html $destdir
    echo "Building application..."
    for app in $APPS; do
        go build -tags "beanstalk nats nsq" -ldflags "$LDFLAGS" -o $destdir/$app$suffix ./$app
    done
    echo "Copying dist files..."
    for folder in $FOLDERS; do
        cp -rf $folder $destdir
    done
    for fname in $FILES; do
        cp -f $fname $destdir
    done
    echo "Creating package dist/$distname.zip...."
    zip -q -r9 dist/$distname.zip $destdir
    rm -rf $destdir
done
echo "All builds done..."

echo "Calculating hashes..."
for $fname in dist/docproc-*.zip; do
    md5sum $fname
done
echo "done"
