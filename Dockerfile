FROM docker.io/library/debian:12.9-slim AS unzipper

SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]

RUN <<EOF
apt-get update
apt-get install -y --no-install-recommends \
  unzip
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF

WORKDIR /app

ARG VAULT_VERSION=1.17.6
ENV TARBALL=vault_${VAULT_VERSION}_linux_amd64.zip
ENV TARBALL_URL=https://releases.hashicorp.com/vault/${VAULT_VERSION}/${TARBALL}
ADD ${TARBALL_URL} /tmp/${TARBALL}
RUN <<EOF
unzip /tmp/${TARBALL} vault
EOF

FROM docker.io/library/debian:12.9-slim

SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]

RUN <<EOF
apt-get update
apt-get install -y --no-install-recommends \
  systemd=252.33-1~deb12u1
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF

COPY --from=unzipper /app/vault /bin/vault

RUN mkdir -p /run/systemd/system

ENV DBUS_SESSION_BUS_ADDRESS=unix:path=/run/dbus/system_bus_socket 
ENV SYSTEMCTL_FORCE_BUS=1

ENTRYPOINT [ "/bin/vault" ]
CMD ["agent", "-config=/etc/vault/agent.hcl"]