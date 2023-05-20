#!/bin/bash
set -euxo pipefail

CRIO_VERSION=1.24
OS=Debian_11

# Download kubernetes archive keyring:
curl -Lfo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
# Install kubernetes apt source:
cat <<EOF > /etc/apt/sources.list.d/kubernetes.list
deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main
EOF

# Download libcontainers archive keyring:
curl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/$OS/Release.key | \
	gpg --dearmor -o /usr/share/keyrings/libcontainers-archive-keyring.gpg
# Install libcontainers apt source:
cat <<EOF > /etc/apt/sources.list.d/devel\:kubic\:libcontainers\:stable.list
deb [signed-by=/usr/share/keyrings/libcontainers-archive-keyring.gpg] https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/Debian_11/ /
EOF

# Download libcontainers crio archive keyring:
curl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable:/cri-o:/$CRIO_VERSION/$OS/Release.key | \
	gpg --dearmor -o /usr/share/keyrings/libcontainers-crio-archive-keyring.gpg
# Install libcontainers crio apt source:
cat <<EOF > /etc/apt/sources.list.d/devel\:kubic\:libcontainers\:stable\:cri-o\:1.24.list
deb [signed-by=/usr/share/keyrings/libcontainers-crio-archive-keyring.gpg] https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable:/cri-o:/1.24/Debian_11/ /
EOF
