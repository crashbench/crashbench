#!/usr/bin/env bash
set -e

echo "Installing CrashBench (The Open Chaos & Safety Benchmark for AI Agents)..."

REPO="crashbench/crashbench"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "Detected OS: $OS ($ARCH)"
echo "Compiling/installing via go install..."
if command -v go >/dev/null 2>&1; then
  go install github.com/crashbench/crashbench/cmd/crashbench@latest
  echo "✅ CrashBench installed successfully! Run 'crashbench help' to get started."
else
  echo "Go not found. Please install Go or download precompiled binaries from https://github.com/$REPO/releases"
  exit 1
fi
