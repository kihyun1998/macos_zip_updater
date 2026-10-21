# Relaunch guarantee

Once the updater has its app path and parent PID, every path ends in a relaunch:
success, backup failure, missing extracted app, install failure, panic.

## Why

The caller does not wait for the updater. acra_client (read at `6dc471565`)
starts it detached and quits itself about a second later
(`update_service_provider.dart — runMacosUpdater` and its caller;
`permission_step_widget.dart` for `--restart`). An updater that exits without
relaunching leaves the user with no app running.

## What it deliberately leaves out

- **The log file is not required.** `initLogger` warns on stderr and continues;
  the logger already tolerates a nil file. A log that cannot be opened is a lost
  diagnostic, not a reason to strand the user (#1).
- **Malformed argv still exits 1 before the wait.** With too few arguments or a
  non-numeric PID there is no app path, or nothing to wait on before relaunching
  — relaunching while the parent is still alive only activates it, and it quits a
  second later. The caller always passes a numeric `pid` and the full argument
  list, so these exits are a contract violation, not a runtime path.

## Blast radius

- A new early `os.Exit` in `main.go` between argument parsing and the parent
  wait reopens #1.
- `launchApp` in `main.go` is the seam `main_test.go` replaces; a relaunch that
  calls `launcher.LaunchApp` directly escapes the tests.
- Changing the caller's argv or its quit-after-spawn timing in acra_client moves
  this guarantee.
