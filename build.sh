#!/bin/bash

cd ~/go/pkg/mod/github.com/ingonyama-zk/icicle/v2@v2.0.3/wrappers/golang

chmod +x ./build.sh

sudo ./build.sh -curve=bn254 -g2
