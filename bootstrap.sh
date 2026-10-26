#!/bin/bash
set -e

# This script installs and sets up shokkupan with a few prerequisites:
# dokku must be installed
# a global domain must be set
# letsencrypt must be configured and enabled with global email set

DOKKU_STORAGE_DIR="/var/lib/dokku/data/storage/"
SHOKKUPAN_DATA_DIR="$DOKKU_STORAGE_DIR/shokkupan"
SHOKKUPAN_APP_DATA_MOUNT_PATH="$SHOKKUPAN_DATA_DIR:/data"
SHOKKUPAN_VERSION=${SHOKKUPAN_VERSION:-"latest"}
SHOKKUPAN_IMAGE="ghcr.io/randspace0/shokkupan:$SHOKKUPAN_VERSION"

SHOKKUPAN_DOKKU_USER="shokkupanadmin"
DISTROLESS_NONROOT_UID="65532"

clean-shokkupan() {
  echo "=> checking for existing resources"

  if dokku apps:exists shokkupan &>/dev/null; then
    echo "==> destroying old dokku app"
    dokku apps:destroy --force shokkupan &>/dev/null
  fi

  if dokku ssh-keys:list "$SHOKKUPAN_DOKKU_USER" &>/dev/null; then
    echo "==> removing existing dokku ssh key"
    dokku ssh-keys:remove $SHOKKUPAN_DOKKU_USER
  fi

  echo "==> done"
}

create-shokkupan-app() {
  clean-shokkupan

  echo "=> pulling image (version: $SHOKKUPAN_VERSION)"
  HOST_SSH_PORT=$(grep "Port " /etc/ssh/sshd_config | awk '{ print $2 }')
  docker pull "$SHOKKUPAN_IMAGE" &>/dev/null
  SHOKKUPAN_IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' "$SHOKKUPAN_IMAGE")

  echo "=> creating & configuring dokku app"
  dokku apps:create shokkupan &>/dev/null
  dokku docker-options:add shokkupan deploy \
    "--add-host=host.docker.internal:host-gateway" &>/dev/null
  dokku config:set shokkupan \
    DOKKU_SSH_HOST='host.docker.internal' \
    DOKKU_SSH_PORT="$HOST_SSH_PORT" &>/dev/null

  echo "==> creating storage"
  dokku storage:ensure-directory shokkupan --chown false &>/dev/null
  dokku storage:mount shokkupan "$SHOKKUPAN_APP_DATA_MOUNT_PATH" &>/dev/null
  chown -R "$DISTROLESS_NONROOT_UID":"$DISTROLESS_NONROOT_UID" "$SHOKKUPAN_DATA_DIR" &>/dev/null

  echo "==> bootstrapping"
  dokku config:set shokkupan DOKKU_SKIP_DEPLOY=true &>/dev/null
  dokku git:from-image shokkupan "$SHOKKUPAN_IMAGE_DIGEST" &>/dev/null

  shokkupan_ssh_key=$(dokku run shokkupan bootstrap) &>/dev/null
  echo "$shokkupan_ssh_key" | dokku ssh-keys:add "$SHOKKUPAN_DOKKU_USER" &>/dev/null

  echo "==> deploying"
  dokku config:unset shokkupan DOKKU_SKIP_DEPLOY &>/dev/null

  if [ -z "$SHOKKUPAN_CERT" ]; then
    echo "==> enabling letsencrypt"
    if [ -n "$SHOKKUPAN_LETSENCRYPT_EMAIL" ]; then
      dokku letsencrypt:set shokkupan email "$SHOKKUPAN_LETSENCRYPT_EMAIL" &>/dev/null
    fi
    # not silenced: a failure here (bad email, DNS, rate limit) must be visible
    dokku letsencrypt:enable shokkupan
  else
    echo "==> adding certificate from SHOKKUPAN_CERT"
    dokku certs:add shokkupan < "$SHOKKUPAN_CERT"
  fi
}

main() {
  if [[ "$(id -u)" != "0" ]]; then
    echo "This script must be run as root" 1>&2
    exit 1
  fi

  if ! command -v dokku &>/dev/null; then
    echo "Please install dokku first using the instructions at https://dokku.com" 1>&2
    exit 1
  fi

  if [ -z "$SHOKKUPAN_CERT" ]; then
    if ! dokku plugin:installed letsencrypt; then
      echo "Please setup letsencrypt using the instructions at https://dokku.com/docs/deployment/application-deployment/#setting-up-ssl" 1>&2
      exit 1
    fi
  fi

  if [ -z "$SHOKKUPAN_CERT" ] && [ -z "$SHOKKUPAN_LETSENCRYPT_EMAIL" ] \
    && [ -z "$(dokku letsencrypt:report --global --letsencrypt-global-email 2>/dev/null)" ]; then
    echo "No letsencrypt email: export SHOKKUPAN_LETSENCRYPT_EMAIL=you@example.com, or run 'dokku letsencrypt:set --global email you@example.com'" 1>&2
    exit 1
  fi

  for plugin in redis postgres mongo mysql; do
    if ! dokku plugin:installed $plugin; then
      echo "=> Installing plugin $plugin"
      dokku plugin:install https://github.com/dokku/dokku-$plugin.git $plugin &>/dev/null
    fi
  done

  create-shokkupan-app

  shokkupan_app_domain=$(dokku domains:report shokkupan --domains-app-vhosts)
  shokkupan_setup_key=$(dokku logs -q shokkupan | grep setup_key | jq ".setup_key")

  echo "---"
  echo "=> shokkupan installed and running "
  echo "--- proceed with setup using key $shokkupan_setup_key at https://$shokkupan_app_domain ---"
}

main
