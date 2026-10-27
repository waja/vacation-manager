#!/usr/bin/env bash
# Prototype helper: atomically install a file only if the destination is absent.

set -euo pipefail

install_no_replace() {
    local source=${1-}
    local target=${2-}
    local directory temporary

    if [[ -z $source || -z $target || ! -f $source ]]; then
        printf 'install-no-replace: source file and target path are required\n' >&2
        return 2
    fi

    directory=${target%/*}
    if [[ $directory == "$target" ]]; then
        directory=.
    elif [[ -z $directory ]]; then
        directory=/
    fi

    if [[ ! -d $directory ]]; then
        printf 'install-no-replace: destination directory does not exist: %s\n' "$directory" >&2
        return 1
    fi

    # Keep the temporary file in the target directory so the hard link stays
    # on the same filesystem. ln creates the destination atomically and fails
    # if any entry (including a symlink) already exists there.
    temporary=$(mktemp -- "$directory/.vacation-manager.XXXXXXXX")
    trap 'rm -f -- "$temporary"' RETURN

    if ! cat -- "$source" > "$temporary"; then
        printf 'install-no-replace: could not prepare temporary file\n' >&2
        return 1
    fi
    if ! chmod --reference="$source" -- "$temporary"; then
        printf 'install-no-replace: could not preserve source mode\n' >&2
        return 1
    fi

    if ! ln -- "$temporary" "$target"; then
        printf 'install-no-replace: destination already exists or atomic link failed: %s\n' "$target" >&2
        return 1
    fi
    rm -f -- "$temporary"
    trap - RETURN
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
    if [[ $# -ne 2 ]]; then
        printf 'Usage: %s SOURCE TARGET\n' "$0" >&2
        exit 2
    fi
    install_no_replace "$1" "$2"
fi
