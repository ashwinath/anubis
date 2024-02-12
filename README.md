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
/proc/self/fd/13:2: command not found: compdef
/proc/self/fd/13:18: command not found: compdef
/home/ashwin/dotfiles/zsh/zshrc:source:33: no such file or directory: /home/ashwin/.passwords
/home/ashwin/dotfiles/zsh/zshrc:source:34: no such file or directory: /home/ashwin/.cargo/env
```
