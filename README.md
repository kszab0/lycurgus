# Lycurgus
Lycurgus is a simple ad blocker that works on all major operating systems.
It works as a proxy therefore won't be blocked or detected by any browser or web page.

## Features
 - HTTP(S) proxy for Windows, Linux and macOS
 - Download blocked hosts list from various sources (blocklists)
 - Parse hosts file format and plain text into proxy rules
 - Blacklist and whitelist with regexp rules

## Build
You must have Go installed in order to build Lycurgus.

### Windows
```
Build.cmd dependencies
Build.cmd
```

### Linux and macOS
```
make dependencies
make all
```

## Usage
You must configure Lycurgus as your system's HTTP and HTTPS proxy by setting them both to localhost:8080.

## License
Authored by Krisóf Szabó and released under the MIT license.