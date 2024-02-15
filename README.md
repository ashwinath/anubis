# Anubis

Script to set up server/work laptop/PC. WIP.

## Joining clusters

Set `KUBEADM_JOIN_TOKEN` and `KUBEADM_JOIN_HASH` as environment variables before running.

## TODO

- sudo dnf --refresh upgrade and prompt a restart, os could be too old even if it's fresh from iso
- gracefully drain and shutdown non master when called via api, some other detection mechanism (time?)
- kubernetes upgrade, test and develop when v1.30 is out.
- mac setup

## Known issues

DNS resolution not working when node just joins:
- restart core dns
