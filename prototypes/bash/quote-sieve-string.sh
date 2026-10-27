#!/usr/bin/env bash
# Prototype helper: quote one single-line value for a Sieve quoted-string.

set -euo pipefail

sieve_quote_string() {
    local value=${1-}

    # Newlines and control characters need context-specific handling. Reject
    # them here rather than emitting a potentially invalid or unsafe literal.
    if [[ $value == *[[:cntrl:]]* ]]; then
        printf 'sieve-quote: control characters are not supported\n' >&2
        return 2
    fi

    value=${value//\\/\\\\}
    value=${value//\"/\\\"}
    printf '"%s"\n' "$value"
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
    if [[ $# -ne 1 ]]; then
        printf 'Usage: %s VALUE\n' "$0" >&2
        exit 2
    fi
    sieve_quote_string "$1"
fi
