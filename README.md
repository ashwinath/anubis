# Anubis

Script to set up server/work laptop/PC. WIP.

## Joining clusters

Set `KUBEADM_JOIN_TOKEN` and `KUBEADM_JOIN_HASH` as environment variables before running.

## TODO

- sudo dnf --refresh upgrade and prompt a restart, os could be too old even if it's fresh from iso
- gracefully drain and shutdown non master when called via api, some other detection mechanism (time?)
- kubernetes upgrade, test and develop when v1.30 is out.
- understand why kubernetes cni has some problem and requires a restart
- get config from github instead of downloading together or use systemd
- mac setup

errors:

```
ERROR: 2024/02/12 14:29:30 error installing binaries, error: could not exec as user, error: output: zsh:1: permission denied: /tmp/anubis/runbinaries/rustup-init
, error: exit status 126
```
