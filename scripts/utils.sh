#!/bin/bash

# Color Definitions
RESET='''\033[0m'''
RED='''\033[0;31m'''
GREEN='''\033[0;32m'''
YELLOW='''\033[0;33m'''
BLUE='''\033[0;34m'''
MAGENTA='''\033[0;35m'''
CYAN='''\033[0;36m'''

show_banner_from_variable() {
  if [ -n "$BANNER" ]; then
    echo -e "${MAGENTA}${BANNER}${RESET}"
    echo
  fi
}

log_step() {
  echo -e "${BLUE}==> $1${RESET}"
}

log_success() {
  echo -e "${GREEN}✓ $1${RESET}"
}

log_warning() {
  echo -e "${YELLOW}⚠ $1${RESET}"
}

log_error() {
  echo -e "${RED}✗ $1${RESET}" >&2
}

log_info() {
  echo -e "${CYAN}i $1${RESET}"
}
