#!/usr/bin/env python3
"""Small language spike for Sieve quoting and atomic no-replace installation."""

from __future__ import annotations

import argparse
import os
import stat
import tempfile
import unicodedata
from pathlib import Path


def quote_sieve_string(value: str) -> str:
    """Quote a single-line value for a Sieve quoted-string."""
    if any(unicodedata.category(character) == "Cc" for character in value):
        raise ValueError("control characters are not supported in quoted strings")
    return '"' + value.replace("\\", "\\\\").replace('"', '\\"') + '"'


def install_no_replace(source: Path, target: Path) -> None:
    """Publish source at target atomically, failing if target already exists."""
    source = Path(source)
    target = Path(target)

    if not source.is_file():
        raise ValueError("source must be an existing regular file")
    directory = target.parent
    if not directory.is_dir():
        raise FileNotFoundError(f"destination directory does not exist: {directory}")

    # A hard link is an atomic no-replace publication primitive on the same
    # filesystem. Place the temporary file beside the destination.
    fd, temporary_name = tempfile.mkstemp(prefix=".vacation-manager.", dir=directory)
    temporary = Path(temporary_name)
    try:
        with os.fdopen(fd, "wb") as output, source.open("rb") as input_file:
            while chunk := input_file.read(1024 * 1024):
                output.write(chunk)
        os.chmod(temporary, stat.S_IMODE(source.stat().st_mode))
        os.link(temporary, target)
    finally:
        try:
            temporary.unlink()
        except FileNotFoundError:
            pass


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    subparsers = parser.add_subparsers(dest="command", required=True)

    quote_parser = subparsers.add_parser("quote", help="quote a single-line Sieve string")
    quote_parser.add_argument("value")

    install_parser = subparsers.add_parser("install", help="install a file without replacing a target")
    install_parser.add_argument("source", type=Path)
    install_parser.add_argument("target", type=Path)

    args = parser.parse_args()
    try:
        if args.command == "quote":
            print(quote_sieve_string(args.value))
        else:
            install_no_replace(args.source, args.target)
    except (OSError, ValueError) as error:
        parser.exit(1, f"vacation-manager prototype: {error}\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
