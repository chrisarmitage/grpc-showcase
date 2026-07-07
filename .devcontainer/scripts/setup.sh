#!/bin/bash

# Dev Container Post-Create Setup Script

set -e  # Exit on any error

echo "[ ] Setting up development environment..."

# Make all scripts executable
chmod +x .devcontainer/scripts/*.sh

# # Run setup scripts
.devcontainer/scripts/setup-git.sh

# Add common aliases to .bashrc
echo "[ ] Configuring shell aliases..."
cat >> ~/.bashrc << 'EOF'

# Common aliases
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'
alias ..='cd ..'
alias ...='cd ../..'
alias grep='grep --color=auto'
EOF

echo "[x] Development environment setup complete!"