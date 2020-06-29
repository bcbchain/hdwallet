#!/usr/bin/env bash

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-"darwin linux windows"}

echo "==> Building..."
cd ..
rm -rf build
mkdir -p build/dist

IFS=' ' read -ra arch_list <<< "$XC_ARCH"
IFS=' ' read -ra os_list <<< "$XC_OS"
for arch in "${arch_list[@]}"; do
	for os in "${os_list[@]}"; do
		echo "--> $os/$arch"
		mkdir -p build/bin/"${os}_${arch}"
		GOOS=${os} GOARCH=${arch} go build -tags="${BUILD_TAGS}" -o "build/bin/${os}_${arch}" ./cmd/...

		pushd "build/bin/${os}_${arch}" > /dev/null || exit 1
		pwd
    tar -zcvf ../../dist/$project_name\_$VERSION\_${os}_${arch}.tar.gz ./*
    popd > /dev/null || exit 1
	done
done