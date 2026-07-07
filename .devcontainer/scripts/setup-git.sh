#!/bin/bash

# Git Configuration Script for Dev Container

set -e

echo "[ ] Setting up Git configuration..."

# Configure Git if credentials are available
if [ -n "$GIT_AUTHOR_EMAIL" ]; then
    git config --global user.email "$GIT_AUTHOR_EMAIL"
    git config --global user.name "$GIT_AUTHOR_NAME"
    echo "[x] Git configured for: $GIT_AUTHOR_NAME <$GIT_AUTHOR_EMAIL>"
else
    echo "[!] No Git credentials found. Set GIT_AUTHOR_EMAIL and GIT_AUTHOR_NAME in host environment"
fi

# Configure Git repository safety
git config --global --add safe.directory /workspaces/grpc-showcase

# Setup Git credential storage
git config --global credential.helper store

# Enable automatic remote setup for new branches
git config --global --type bool push.autoSetupRemote true

echo "[x] Git configuration complete!"