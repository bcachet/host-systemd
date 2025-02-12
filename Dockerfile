FROM docker.io/library/debian:12.9

SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]

# ARG VAULT_VERSION=1.17.6
# ENV TARBALL=vault_$(VAULT_VERSION)_linux_amd64.zip
# ENV TARBALL_URL=https://releases.hashicorp.com/vault/$(VAULT_VERSION)/$(TARBALL)

RUN <<EOF
apt-get update
apt-get install -y --no-install-recommends \
  wget \
  ca-certificates \
  lsb-release \
  gpg
wget -O - https://apt.releases.hashicorp.com/gpg | gpg --dearmor | tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
gpg --no-default-keyring --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg --fingerprint
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list
apt-get update
apt-get install -y --no-install-recommends \
  systemd=252.33-1~deb12u1 \
  vault=1.17.6-1
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF

RUN mkdir /run/systemd/system

ENV DBUS_SESSION_BUS_ADDRESS=unix:path=/run/dbus/system_bus_socket 
ENV SYSTEMCTL_FORCE_BUS=1

