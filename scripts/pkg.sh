#!/usr/bin/env bash

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-"darwin linux windows"}

echo "==> Building..."
cd ..

IFS=' ' read -ra arch_list <<< "$XC_ARCH"
IFS=' ' read -ra os_list <<< "$XC_OS"
for arch in "${arch_list[@]}"; do
	for os in "${os_list[@]}"; do
		echo "--> $os/$arch"
		GOOS=${os} GOARCH=${arch} go build -tags="${BUILD_TAGS}" -o "build/pkg/${os}_${arch}/pieces/bcchain" ./cmd/...
	done
done