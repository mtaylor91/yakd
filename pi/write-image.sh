#!/bin/bash
set -euxo pipefail
export RPI_MODEL=4
export DEBIAN_RELEASE=bullseye
export SD_CARD=$1
xzcat raspi_${RPI_MODEL}_${DEBIAN_RELEASE}.img.xz | \
  dd of=${SD_CARD} bs=64k oflag=dsync status=progress
