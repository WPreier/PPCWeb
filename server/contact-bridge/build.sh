#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
go mod tidy
go build -o ppc-contact .
echo "built ./ppc-contact"
./ppc-contact -h 2>/dev/null || true
