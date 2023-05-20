#!/bin/bash
set -euxo pipefail
export RPI_MODEL=4
export DEBIAN_RELEASE=bullseye
curl -LfO https://raspi.debian.net/daily/raspi_${RPI_MODEL}_${DEBIAN_RELEASE}.img.xz
curl -LfO https://raspi.debian.net/daily/raspi_${RPI_MODEL}_${DEBIAN_RELEASE}.img.xz.sha256
sha256sum -c raspi_${RPI_MODEL}_${DEBIAN_RELEASE}.img.xz.sha256
