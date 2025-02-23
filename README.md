# unchained
Unchained is a vpn/proxy service that is very easy to setup (proto: trojan/vless + reality)

# Goals
One app, one line setup

Sets up one proxy on one port with one [protocol](#Available-protocols)

Prints one connection URL and one qr code

That's it :)

# Available protocols
- Trojan + Reality (default)
- Vless + VISION + Reality

# Run with docker
`docker run -d --name unchained --network host -e TAGS=myserver ghcr.io/awryme/unchained:latest`

Example `docker-compose.yml` provided in repo

# Underlying implementation details (subject to change)
- Config is stored in a json file that can be edited if needed
    - All possible parameters are in this file
    - Only a few parameters a provided as flags/env for the 'run' command
- Public IP is retrieved using https://api.ipify.org/
- IPv6 support is detected by interfaces and by dialing google.com (ipv6)
- Actual proxy implementation is [sing-box](https://github.com/SagerNet/sing-box), big thanks to them
- Name for the proxy is formed by `$ID_$TAGS_$PROTO`

# Help output
`unchained -h`

```
Usage: unchained <command> [flags]

unchained is a vpn/proxy application that sets up everything for you

Flags:
  -h, --help                         Show context-sensitive help.
  -c, --config="./unchained.json"    file to store generated/edited config file ($CONFIG)

Commands:
  run [flags]
    run vpn server, generates config if it doesn't exist (default command if no other provided)

  print [flags]
    print connection info for client

  generate [flags]
    generate config without running the server

  reset [flags]
    cleans up configs/files used by this command

Run "unchained <command> --help" for more information on a command.
```

`unchained run -h` (`generate` flags are the same, except `--no-config`)

```
Usage: unchained run [flags]

run vpn server, generates config if it doesn't exist

Flags:
  -h, --help                         Show context-sensitive help.
  -c, --config="./unchained.json"    file to store generated/edited config file ($CONFIG)

      --log-level="warn"             sing-box log level ($LOG_LEVEL)
      --dns="tls://1.1.1.1"          dns address, in sing-box format ($DNS)
  -p, --proto="trojan"               set used protocol: trojan,vless ($PROTO)
      --id=STRING                    proxy id (used to identify proxy in client apps), random by default ($ID)
      --tags=TAGS,...                proxy tags (used to identify proxy in client apps) ($TAGS)
      --no-config                    only generate config, ignore existing ($NO_CONFIG)
```