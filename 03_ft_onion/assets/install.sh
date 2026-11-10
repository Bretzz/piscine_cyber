#!/bin/sh
set -e

#  Install packages
apk add --no-cache nginx tor su-exec openssh

# Generate ssh keys and set root password
ssh-keygen -A
source /tmp/assets/.env
if [ -z "$PASSWD" ]; then
    echo "Missing root password, exiting."
    exit 1
fi
echo "root:$PASSWD" | chpasswd

# Create the hidden service directory
mkdir -p /var/lib/tor/ft_onion/
# Tor requires 700 permissions and strict ownership by the tor user
chown -R tor:tor /var/lib/tor/ft_onion/
chmod 700 /var/lib/tor/ft_onion/

# Allow Nginx to interact with Tor's runtime group for sockets
# Add nginx user to the tor group so it can read/write to Tor-created sockets
addgroup nginx tor

# Replace default configs
mv /tmp/assets/sshd_config /etc/ssh/sshd_config
mv /tmp/assets/nginx.conf /etc/nginx/nginx.conf
mv /tmp/assets/torrc /etc/tor/torrc

# Move entrypoint and clean up tmp
mv /tmp/assets/entrypoint.sh /usr/local/bin/entrypoint.sh
chmod +x /usr/local/bin/entrypoint.sh
rm -rf /tmp/assets
