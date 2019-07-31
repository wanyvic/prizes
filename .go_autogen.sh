#!/usr/bin/env bash
GITCOMMIT=$(git rev-parse --short HEAD)
VERSION=$(git describe --abbrev=0 --tags)
BUILDTIME=$(date -u -d "@${SOURCE_DATE_EPOCH:-$(date +%s)}" --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
cat > prizesversion/version.go <<DVEOF
// +build autogen
// Package prizesversion is auto-generated at build-time
package prizesversion
// Default build-time variable for library-import.
// This file is overridden on build with build-time information.
const (
	GitCommit             string = "$GITCOMMIT"
	Version               string = "$VERSION"
	BuildTime             string = "$BUILDTIME"
)
DVEOF
