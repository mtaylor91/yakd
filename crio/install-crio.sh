#!/bin/bash
set -euxo pipefail
apt-get install -y cri-o cri-o-runc
apt-mark hold cri-o cri-o-runc
systemctl enable --now crio
