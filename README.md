# Vacation Manager

A small, non-interactive command-line tool for managing Dovecot/Pigeonhole Sieve vacation responders.

The tool is intended for systems that use Dovecot and Pigeonhole Sieve. Froxlor is the current target environment, but **Froxlor is not a dependency**: the tool should work with other setups that provide compatible Dovecot/Pigeonhole tooling and filesystem permissions.

> **Project status:** Requirements and implementation are in progress. The command-line interface and installation instructions below will be finalized as implementation decisions are made.

## Goals

- Generate a Sieve vacation script for an email address and a configured absence period.
- Validate the generated script with `sievec` before installing it.
- Install only the `.sieve` source file. Do not manage the generated `.svbin` file or restart/reload Dovecot; Pigeonhole handles compilation as needed.
- Operate non-interactively and fail safely with actionable errors.
- Keep runtime dependencies minimal and avoid control-panel-specific behavior.

## Target layout

By default, the target base directory is `/var/customers/sieve`. For an address such as `waja@test.org`, the intended destination is:

```text
/var/customers/sieve/test.org/waja/vacation.sieve
```

The base directory will be configurable through an optional configuration file and a command-line override. The tool must validate the derived path and avoid writing outside the configured base directory.

## Planned inputs

The command-line interface is expected to accept:

- Email address
- Inclusive start and end dates in `YYYY-MM-DD` format
- Vacation subject and a message-body file (multiline)
- One or more addresses used by the Sieve `vacation :addresses` argument
- Optional time zone (default: `Europe/Berlin`)
- Optional vacation interval in days (default: `7`)

An optional configuration file will hold rarely changed defaults such as `base_path`, `timezone`, and `days`. Command-line options take precedence over configuration values, which take precedence over built-in defaults. The tool will not create or modify the configuration file automatically.

Exact option names and configuration-file conventions are not finalized yet. See [Issue #1](https://github.com/waja/vacation-manager/issues/1) for the living requirements and open decisions.

## Safety

The implementation is expected to:

- Validate input values, date ranges, time zones, and message-file readability.
- Escape user-supplied data safely for Sieve.
- Run `sievec` against the generated script before installation.
- Leave the current target untouched if validation or compilation fails.
- Avoid overwriting an existing script unless explicitly requested.
- Install the validated source atomically.
- Use temporary directories in tests rather than the production default path.

The exact Sieve conditions and required capabilities will be checked against Dovecot/Pigeonhole behavior and documented before release.

## Development

Repository development instructions are in [AGENTS.md](AGENTS.md).

The project uses these branch conventions:

- `development` is the default branch for ongoing work.
- Changes are developed on focused branches and proposed to `development` using pull requests.
- `master` is reserved for released versions and will be created for the first release.
- Releases are tagged (for example, `v1.0.0`) and promoted from `development` to `master` through a release pull request.

Testing, linting, installation, and usage commands will be documented here as the corresponding implementation is added.

## License

Vacation Manager is licensed under the **GNU General Public License, version 3 or (at your option) any later version** (SPDX: `GPL-3.0-or-later`). See [COPYING](COPYING) for the full license text.
