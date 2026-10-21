# thegraph build

## What this project is

The macOS zip-based self-updater that acra_client launches to swap in a new app bundle and restart it.

## References

| Source | Informs | Reached by | Binding |
|---|---|---|---|
| acra_client (the caller that launches the updater and pins it by tag) | how it works | its source tree, `../../` from this clone — raw | binding |
| Apple codesign / hardened runtime / notarization docs | how it works | the published docs — **summarized** | binding |
