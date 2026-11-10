#!/bin/sh
set -e

# Start ssh server
echo "Starting SSH..."
/usr/sbin/sshd -D -e &

# Fix runtime permissions just in case
chown -R tor:tor /var/lib/tor/ft_onion/
chmod 700 /var/lib/tor/ft_onion/

# Start Tor in the background
echo "Starting Tor..."
su-exec tor tor -f /etc/tor/torrc &

# Block and wait for the hostname file to be created by Tor
echo "Waiting for Tor to generate .onion keys..."
while [ ! -f /var/lib/tor/ft_onion/hostname ]; do
    sleep 1
done

# Update nginx.conf with the generated address
ONION_ADDR=$(cat /var/lib/tor/ft_onion/hostname)
echo ">"
echo "> Onion Address: $ONION_ADDR"
echo ">"

sed -i "s|REPLACE_WITH_GENERATED_ADDRES|${ONION_ADDR}|g" /etc/nginx/nginx.conf

# Ensure the socket created by Tor can be read/written by the 'tor' group (which Nginx belongs to)
if [ -S /var/lib/tor/ft_onion/onion.sock ]; then
    chmod 660 /var/lib/tor/ft_onion/onion.sock
fi

# Start Nginx in the foreground (keeps container alive)
echo "Starting Nginx..."
exec nginx -g "daemon off;"
