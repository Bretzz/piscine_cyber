FROM alpine:latest

# Copy all build/runtime scripts and assets
COPY ./assets/ /tmp/assets/

# Run the installation script
RUN chmod +x /tmp/assets/install.sh && /tmp/assets/install.sh

# Copy the application
COPY --chown=nginx:nginx ./src/ /var/www/html/

# Mount the secret file (named 'ssh_key') during this single RUN command
# and append it directly to authorized_keys
RUN --mount=type=secret,id=ssh_key \
    mkdir -p /root/.ssh && \
    cat /run/secrets/ssh_key >> /root/.ssh/authorized_keys && \
    chmod 700 /root/.ssh && \
    chmod 600 /root/.ssh/authorized_keys

# Expose port (9050 if you want to use this instance of Tor as entrypoint for browsing)
# EXPOSE 80 4242
EXPOSE 9050

# Run the entrypoint
ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
