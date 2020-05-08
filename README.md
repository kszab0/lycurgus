# Lycurgus
<img src="logo.png" alt="logo" width="128"/>

> Lycurgus is a simple ad blocker that works on all major operating systems.
It works as a proxy therefore won't be blocked or detected by any browser or web page.

## Features
 - HTTP(S) proxy for Windows, Linux and macOS
 - Download blocked hosts list from various sources (blocklists)
 - Parse hosts file format and plain text into proxy rules
 - Blacklist and whitelist with regexp rules
 - Run application automatically on startup
 - Configurable upstream proxy

## Build
You must have [Go](https://golang.org/) installed in order to build Lycurgus.

### Windows
```
Build.cmd
```

### Linux and macOS
```
make
```

Building the app requires gtk3, libappindicator3 and libwebkit2gtk-4.0-dev development headers to be installed. You can install them for Debian or Ubuntu using:

```
sudo apt-get install libgtk-3-dev libappindicator3-dev libwebkit2gtk-4.0-dev
```

## Usage
You must configure Lycurgus as your system's HTTP and HTTPS proxy. Here is how you do that in [Chrome](https://www.simplified.guide/google-chrome/change-proxy-setting) and in [Firefox](https://www.wikihow.com/Enter-Proxy-Settings-in-Firefox). For example, the URL you should use as proxy when running Lycurgus with default settings is `localhost:5678`. The the URL can be set with the `--address` command line flag.

### Blocklist
The default blocklist is coming from the advertising lists found on the awesome [firebog's adblock list](https://firebog.net/). This can be configured by creating a file named `blocklist` in the config directory and listing URLs pointing to a hosts file (eg.:[ad-wars](https://raw.githubusercontent.com/jdlingyu/ad-wars/master/hosts)) or simple text file (eg: [AdguardDNS.txt](https://v.firebog.net/hosts/AdguardDNS.txt)) one at a line. The blocklist file location can be set with the `--blocklist` command line flag.

### Blacklist
The blacklist can be created in the config directory with the name `blacklist`. You can specify custom regexp rules (one by line) for domains that you would like to block. The blacklist file location can be set with the `--blacklist` command line flag.

### Whitelist
The whitelist can be created in the config directory with the name `whitelist`. You can specify custom regexp rules (one by line) for domains that you would like to allow even if they are blocked by either the blocklist or blacklist. The whitelist file location can be set with the `--whitelist` command line flag.

### Config
The application can be configured with a yaml config file(named `lycurgus.yml`) in the config directory. All the flags can be used as keys in the config file. An example config can be found in the testdata folder. The flags will always have precedence over the values set in the config file. The settings not present in either the config file or flags will have their default values.

| Setting | Flag/Config key | Default Value |
| ------- | ---- | ------------- |
| http address to run blocker | address | :8080 |
| path to blocklist file | blocklist | <config_dir>/blocklist |
| path to blaclist file | blacklist | <config_dir>/blacklist |
| path to whitelist filer | whitelist | <config_dir>/whitelist |
| enable autostart | autostart | true |
| enable GUI | gui | true |
| enable logging | log | true |
| path to logfile | logfile | <log_dir>/lycurgus.log |
| upstream proxy address | proxy | no set |

#### Directories
|   | Windows | Linux/BSDs | macOS |
| - | ------- | ----- | ----- |
| **Config** | `%APPDATA%\lycurgus` (`C:\Users\%USERNAME%\AppData\Local\lycurgus`) | `$XDG_CONFIG_HOME/lycurgus` (`$HOME/.config/lycurgus`) | `$HOME/Library/Application Support/lycurgus` |
| **Log** | `%LOCALAPPDATA%\lycurgus` (`C:\Users\%USERNAME%\AppData\Roaming\lycurgus`) | `$XDG_CACHE_HOME/<name>/logs/lycurgus` | `$HOME/Library/Logs/lycurgus` |

## License
Authored by [Kristóf Szabó](mailto:kristofszabo@protonmail.com) and released under the MIT license.