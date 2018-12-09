#!/usr/bin/env bash
TESTDATA=./testdata/cache_files
TESTING_GO=internal/pkg/testfixtures
RUN_GO=cmd/ngxcpd/main.go
ZONES=2

cd "${BASH_SOURCE[0]%/*}/.."

rm -vrf $TESTDATA/zone*

export UID
docker-compose up --force-recreate --build -d

for i in $(seq 1 1000); do
    for z in $(seq 1 $ZONES); do
        curl http://localhost:8080/zone$z/ok/$RANDOM > /dev/null
    done
done

for i in $(seq 1 15); do
    for z in $(seq 1 $ZONES); do
        curl http://localhost:8080/zone$z/error/$RANDOM > /dev/null
        curl http://localhost:8080/zone$z/notfound/$RANDOM > /dev/null
    done
done

docker-compose stop

for z in $(seq 1 $ZONES); do
    TMP=$(mktemp)
    OUT=$TESTING_GO/zone${z}_files.go

    echo "// AUTOMAGICALLY GENERATED
// +build fixtures

package testfixtures

import (
	\"github.com/maetthu/ngxcpd/pkg/proxycache\"
	\"time\"
)

// TestdataCacheFilesZone${z} contains expected metadata of files in testdata/cache_files/zone${z}
var TestdataCacheFilesZone${z} = []*proxycache.Entry{
    " > $TMP
    find $TESTDATA/zone$z -type f | xargs go run $RUN_GO inspect -t | sed "s#$TESTDATA/zone$z/##" >> $TMP

    echo '}' >> $TMP
    mv $TMP $OUT
    go fmt $OUT
done