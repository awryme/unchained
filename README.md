# unchained
Unchained is a vpn/proxy service that is very easy to setup (proto: trojan + reality)

# Goals
One app, one line setup

Sets up one proxy on one port with one [protocol](#Available-protocols)

Prints one connection URL and one qr code

That's it :)

# Available protocols
- Trojan + Reality

# Underlying implementation details (subject to change)
- Config is stored in a json file that can be edited if needed
    - All possible parameters are in this file
    - Only a few parameters a provided as flags/env for the 'run' command
- Public IP is retrieved using https://api.ipify.org/
- Actual proxy implementation is [sing-box](https://github.com/SagerNet/sing-box), big thank you