#!/bin/bash
export PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
wget https://github.com/nadoo/glider/releases/download/v0.16.3/glider_0.16.3_linux_amd64.tar.gz
tar -xzvf glider_0.16.3_linux_amd64.tar.gz
cp ./glider_0.16.3_linux_amd64/glider . && rm -rf glider_0.16.3_linux_amd64
./glider -listen ws://:6781,vless://e52d7225-9450-3c9d-0b29-6dc1baea56dd@ &
./app