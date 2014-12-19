#!/bin/bash

set -e

function timecmd {
    { time -p "$@" ; } 2>&1 | tail -n 3 | grep real | cut -f2 -d' '
}

all_embedders="zhex zbase64 hex base64 quote"

embedders=$@
if [ "" == "$embedders" ]; then
    embedders=$all_embedders
fi

pushd cmd/goembed >/dev/null
go build
popd >/dev/null

pushd embedtesting >/dev/null
../cmd/goembed/goembed -package=embedtesting testdata
popd >/dev/null

for e in $embedders; do
    echo "# $e: go test"
    pushd ${e}embedder >/dev/null
    go test -bench=. -cpu 1,4 -benchtime 5s
    popd >/dev/null

    echo "# $e: goembed"
    pushd cmd/goembedtest >/dev/null
    wdir="$(mktemp -d)"
    # Amplify the size of the test data
    for i in `seq 0 63`; do
	cp -r testdata "$wdir/$i"
    done
    ../goembed/goembed -c=true -e $e "$wdir"
    echo -e "size\t$(du -h assets.generated.go | cut -f1)"
    gofmt assets.generated.go > assets.generated.go.gofmt
    diff -u assets.generated.go assets.generated.go.gofmt
    rm assets.generated.go.gofmt
    echo "# $e: goembedtest"
    go clean
    go build
    go clean
    t="$(timecmd go build)"
    echo -e "build\tgoembedtest\t${t}s"
    ./goembedtest "$wdir"
    echo -ne "time\t"
    for i in `seq 1 5`; do
	t="$(timecmd ./goembedtest -q '$wdir')"
	echo -n "${t}s  "
    done
    echo
    rm -rf "$wdir"
    popd >/dev/null

    echo
done

echo
echo "ALL GOOD"
echo
echo
