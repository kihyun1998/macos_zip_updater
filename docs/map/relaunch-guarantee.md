# Relaunch guarantee

Once the updater has its app path and parent PID, every path ends in a relaunch:
success, backup failure, missing extracted app, install failure, panic.

## Why

The caller does not wait for the updater. acra_client (read at `6dc471565`)
starts it detached and quits itself about a second later
(`update_service_provider.dart — runMacosUpdater` and its caller;
`permission_step_widget.dart` for `--restart`). An updater that exits without
relaunching leaves the user with no app running.

## Which bundle comes back

The relaunch opens the bundle at `app_path` (`open <app_path>`), never by app
name. acra_client passes the path of the bundle that was running, and its
installed and portable builds share one bundle name and ID
(`macos/Runner/Configs/AppInfo.xcconfig`: `ACRAPointClient`,
`kr.co.penta.ACRAPoint`).

Measured 2026-10-21 with two same-named probe bundles registered through
`lsregister`: `open -a <name>` resolves only from the standard app folders. With
a copy in `~/Applications` it opened that copy every time, even right after the
other copy had run; with no copy there it failed with *Unable to find
application* and opened nothing. `open <path>` opened exactly that bundle (#2).

`LaunchApp` waits for `open` and logs its output on failure. `open` returns once
the launch is requested, so the wait is short; before, a failed `open` was never
logged and the log always said *launched successfully*.

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
