#!/bin/bash
set -euxo pipefail

cat <<EOF > /etc/modules-load.d/modules.conf
br_netfilter
EOF

cat <<EOF > /etc/sysctl.conf
net.ipv4.ip_forward=1
net.bridge.bridge-nf-call-iptables=1
EOF

apt-get install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl
