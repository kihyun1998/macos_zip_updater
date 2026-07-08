# macos_zip_updater

A macOS zip-based application updater written in Go. Core packages: `installer/`,
`launcher/`, `monitor/`, `logger/`, with the entry point in `main.go`. Builds the
universal binary `ACRA_Point_Client_Updater` (plus per-arch `_AMD64` / `_ARM64`).

## Versioning & Release

Releases are pinned by acra_client (via its `accio.json` + `accio deps`) using git
tags in **semver** form `vMAJOR.MINOR.PATCH` (e.g. `v1.0.0`). Two-part tags like
`v1.0` are not valid Go module versions and will not stamp into the binary.

Release ritual:

1. Commit all changes — the working tree must be **clean**.
2. Tag the release commit and push it:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. Build with `make` — the build injects the version via
   `-ldflags "-X main.version=$(git describe --tags --always --dirty)"`.

4. Verify:

   ```bash
   bin/ACRA_Point_Client_Updater --version
   ```

   At a clean, tagged commit it prints the tag (e.g. `v1.0.0`). A
   `-<n>-g<sha>-dirty` suffix means the build was ahead of the tag and/or from
   uncommitted changes — do not ship it.

5. Commit the freshly built binary.

acra_client records the chosen tag in `accio.json` and clones this repo at that
tag. It never moves an existing tag; a re-pin is a new tag + an `accio.json` edit.
