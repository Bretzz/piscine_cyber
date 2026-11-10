FROM alpine:latest

# Copy all build/runtime scripts and assets
COPY ./assets/ /tmp/assets/

# Run the installation script
RUN chmod +x /tmp/assets/install.sh && /tmp/assets/install.sh

# Copy the application
COPY --chown=nginx:nginx ./src/ /var/www/html/

# Expose port (9050 if you want to use this instance of Tor as entrypoint for browsing)
EXPOSE 80 4242

# Run the entrypoint
ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
