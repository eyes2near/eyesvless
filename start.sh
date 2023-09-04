#!/bin/bash
export PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
./glider -listen ws://:6781,vless://e52d7225-9450-3c9d-0b29-6dc1baea56dd@ &
./app