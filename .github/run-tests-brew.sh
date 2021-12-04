#!/bin/sh
# =============================================================================
#  Test Script for Homebrew Functionality
# =============================================================================

set -eu

git config --global init.defaultBranch main
brew install Qithub-BOT/QiiTask/qiitask

# Smoke test
qiitask --version
qiitask hello | grep Hi
