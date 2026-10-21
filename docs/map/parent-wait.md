# Parent wait

Before touching the bundle, the updater polls once a second until the parent PID
is gone (`monitor.MonitorParentProcess`). `IsProcessAlive` asks the kernel with
`kill(pid, 0)`: no error or `EPERM` is alive, `ESRCH` is gone.

## Why

The swap renames the bundle the parent is running from, so "gone" must mean the
process is gone and nothing else. The earlier check ran `ps -p <pid>` and read
any failure as gone. Measured 2026-10-21 on darwin: `ps -p` exits 1 for a gone
PID, for a PID too large (`process id too large`) and for PID 0, and a `ps` that
cannot run at all (not on `PATH`) failed the same way. A failing `ps` started the
swap with the app still running, and every ordinary exit was logged as `ERROR`
with a stack trace (#3).

`kill(pid, 0)` measured from Go on the same machine: own PID → nil; PID 1
(root-owned) → `EPERM`; exited and reaped child → `ESRCH`; zombie → nil;
99999999 → `ESRCH`; **0 → nil**; -5 → `ESRCH`. `man 2 kill` lists only
`EINVAL`, `EPERM` and `ESRCH`.

## What it deliberately leaves out

- **PID 0 and below read as gone without a syscall.** To `kill` they name a
  process group; PID 0 is the caller's own group and always answers nil, which
  would wait forever. `ps -p 0` read as gone, so this keeps that answer.
- **Any other `kill` error reads as gone and is logged.** Only `EINVAL` remains,
  and signal 0 is valid, so the path is not expected to run; waiting forever
  would strand the user the way #1 did.
- **PID reuse is not detected.** A new process given the parent's PID reads as
  the parent. `ps -p` had the same blind spot.
- **A zombie reads as alive.** The parent's own parent (launchd) reaps it.

## Blast radius

- Both modes wait here before relaunching or swapping (`main.go`); see
  [relaunch-guarantee](relaunch-guarantee.md).
- `monitor/process_test.go` pins each row above against real processes.
