# AGENTS.md
---

## Project scope

- Build a small, non-interactive Python CLI for generating and safely installing Dovecot/Pigeonhole Sieve vacation scripts.
- Froxlor is the current target environment, not a runtime dependency. Do not couple core behavior to Froxlor.
- Treat GitHub Issue #1 as the living requirements record. Update it when material requirements or design decisions change.

## Language and documentation

- Write project documentation, code comments, CLI help, errors, tests, and commit messages in English.
- Keep the README accurate as implementation and supported behavior evolve.
- Document user-visible options, defaults, prerequisites, and safety behavior.
- Do not document planned features as implemented; clearly label unfinished areas.

## Design principles

- Prefer the Python standard library and avoid unnecessary runtime dependencies.
- Keep the CLI non-interactive. Invalid input or operational failures must produce a useful error and a non-zero exit status.
- Keep responsibilities small and testable: argument/config handling, validation, Sieve generation, compiler validation, and installation.
- Prefer explicit, readable code over generic frameworks or speculative extensibility.
- Do not introduce Froxlor-specific logic where standard Dovecot/Pigeonhole integration is sufficient.

## Sieve and filesystem safety

- Treat all user input as untrusted, including email addresses, subject text, addresses, dates, and message contents.
- Escape values according to Sieve syntax; never interpolate untrusted strings without appropriate encoding.
- Validate the generated script with `sievec` before installation. Direct compiler output to a disposable temporary location and clean it up; never let validation create compiler artifacts in the managed destination.
- On validation or compiler failure, do not alter an existing destination.
- Publish new files atomically with no-replace semantics, so concurrent invocations cannot overwrite a script installed by another invocation. Do not rely on a check-then-rename sequence, because the absence check races. Prevent paths derived from email input from escaping the configured base directory.
- Never overwrite an existing vacation script by default. Any overwrite behavior must be explicit and documented.
- Manage only the `.sieve` source in the configured destination. Any `.svbin` output needed for `sievec` validation must be directed to a disposable temporary location and removed afterward. Do not create, replace, or delete `.svbin` files in the managed destination, and do not restart or reload Dovecot.
- Tests must use temporary directories and must not write to `/var/customers/sieve` or require a live mail system.

## Testing and quality

- Add tests for behavior and failure paths, including path derivation, validation, config precedence, Sieve escaping/output, compiler failure, and safe installation.
- Tests must be deterministic and runnable locally without a production server. Isolate tests that require `sievec` and document the prerequisite.
- Every implementation change must include appropriate automated tests; bug fixes should add regression tests where practical.
- On every pull request targeting `development` and every push to `development`, GitHub Actions must run the required baseline automated test suite and configured lint/quality checks. These required checks must not be skipped based on changed paths or a judgment that a change is unrelated. The pipeline must test the software, not merely lint it or check the environment.
- Use path filters or job-level change detection only for additional, specialized, or expensive jobs beyond the required baseline checks. Such jobs should run when relevant files change and normally skip documentation-only or unrelated changes. For example, a Sieve compiler compatibility matrix may be conditional on Sieve fixtures, generation logic, or its own workflow, while a broader compatibility sweep may be scheduled or manually triggered.
- The agent must actively manage the pull request's Draft status. Set the PR to **Draft** when starting or resuming implementation, experimentation, or debugging; keep it Draft while work is in progress. Once the change is complete, coherent, relevant local checks have been run, and known limitations are documented, the agent should mark the PR **Ready for review** without waiting for the user to ask. If follow-up work returns the PR to an implementation or debugging phase, set it back to Draft, then mark it Ready again when the follow-up is complete. Do not mark a PR Ready merely because no more work is planned if tests fail or a material issue remains unresolved.
- Avoid triggering review bots and optional expensive CI jobs for every intermediate commit when the PR is still being developed. Configure workflows to defer the full PR validation suite until the PR is marked ready for review, then run it on subsequent updates to that ready PR. Required baseline checks must still run for every ready-for-review PR update and every push to `development`; do not use changed-path filtering to skip them. Optional specialized jobs may use change detection or manual/scheduled triggers.
- Draft status is not a substitute for testing: run targeted local tests during development. Draft PRs may skip the full remote baseline suite until ready, but the required suite must run when the PR is marked ready and on subsequent updates. Check review-bot settings separately because Draft status does not guarantee that every bot will stay silent.
- Verify GitHub branch-protection and required-status-check behavior when implementing draft-aware workflows, so skipped draft jobs do not block the PR once it is ready and all mandatory checks are enforced.
- Required checks must pass before a PR is considered ready to merge. Do not make the test job optional or silently skip tests because the test directory is empty.
- Include a CI workflow with the initial implementation, so the first code PR introduces both meaningful tests and an enforced development pipeline. Document the local commands and CI behavior in the README.
- Run the relevant test suite and configured lint/format checks before proposing changes.
- Keep CI proportionate to the project and use it to run real tests, not just environment or image checks.
- Do not add abstractions, validators, dependencies, or CI jobs without a concrete requirement.

## Git and pull requests

- Keep changes focused and commit history reviewable. Do not commit exploratory/debug churn.
- Use short-lived branches for implementation, fixes, and documentation; open pull requests against `development`.
- `development` is the default integration branch.
- Reserve `master` for released and tagged versions; create it for the first release rather than using it for development.
- Promote releases from `development` to `master` through a release pull request and tag the release commit (for example, `v1.0.0`).
- Bring fixes made for released versions back into `development`.
- Update the README when behavior, setup, configuration, or supported usage changes.
- PR descriptions should state the purpose, important decisions, and checks performed.
