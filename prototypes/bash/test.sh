#!/usr/bin/env bash
set -euo pipefail

root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
quote="$root/quote-sieve-string.sh"
install="$root/install-no-replace.sh"
temporary_directory=$(mktemp -d)
trap 'rm -rf -- "$temporary_directory"' EXIT

assert_equal() {
    local expected=$1 actual=$2 description=$3
    if [[ $actual != "$expected" ]]; then
        printf 'not ok - %s\nexpected: <%s>\nactual:   <%s>\n' "$description" "$expected" "$actual" >&2
        exit 1
    fi
    printf 'ok - %s\n' "$description"
}

quoted=$("$quote" 'Hello "world" \ mailbox')
expected=$'"Hello \\"world\\" \\\\ mailbox"'
assert_equal "$expected" "$quoted" 'quote Sieve string escapes quotes and backslashes'

if "$quote" $'line one\nline two' >/dev/null 2>&1; then
    printf 'not ok - reject multiline quoted string\n' >&2
    exit 1
fi
printf 'ok - reject multiline quoted string\n'

source_file="$temporary_directory/source.sieve"
target_file="$temporary_directory/vacation.sieve"
printf 'require ["vacation"];\n' > "$source_file"
chmod 0640 "$source_file"

"$install" "$source_file" "$target_file"
assert_equal "$(cat -- "$source_file")" "$(cat -- "$target_file")" 'install file content'
assert_equal '640' "$(stat -c '%a' "$target_file")" 'preserve source mode'

printf 'existing target must survive\n' > "$target_file"
if "$install" "$source_file" "$target_file" >/dev/null 2>&1; then
    printf 'not ok - refuse existing destination\n' >&2
    exit 1
fi
assert_equal 'existing target must survive' "$(cat -- "$target_file")" 'existing destination remains unchanged'

if "$install" "$source_file" "$temporary_directory/missing/vacation.sieve" >/dev/null 2>&1; then
    printf 'not ok - reject missing destination directory\n' >&2
    exit 1
fi
printf 'ok - reject missing destination directory\n'

printf 'All Bash prototype tests passed.\n'
