# WeeChat Xochimilco

[![CI](https://github.com/oxzi/weechat-xochimilco/actions/workflows/ci.yml/badge.svg)](https://github.com/oxzi/weechat-xochimilco/actions/workflows/ci.yml)
[![REUSE status](https://api.reuse.software/badge/github.com/oxzi/weechat-xochimilco)](https://api.reuse.software/info/github.com/oxzi/weechat-xochimilco)

A _proof-of-concept_ [WeeChat][weechat-main] plugin to encrypt `PRIVMSG`s with a variant of Signal's cryptography.
This [plugin][weechat-plugin] (not a script) is written in both Golang and C and is based on my [xochimilco][] library for encryption.
It is mostly an experiment.
So please do not use it for anything serious.

[![asciicast](https://asciinema.org/a/uc0FHddfKOXxczu8CVCGVvP3e.svg)](https://asciinema.org/a/uc0FHddfKOXxczu8CVCGVvP3e)


## Install

To build the plugin, at least a recent Go compiler and a GCC must be installed.
If you are using the Nix package manager, please enter a `nix-shell`.

```sh
# Create the plugin, xochimilco.so
make

# Now, you might wanna copy this shared object to your WeeChat's plugin dir.
# For the default location, ~/.weechat/, there is:
make install

# For testing, a temporary and preconfigured WeeChat can be launched:
make test-instance
```


[weechat-main]: https://weechat.org/
[weechat-plugin]: https://weechat.org/files/doc/stable/weechat_plugin_api.en.html
[xochimilco]: https://github.com/oxzi/xochimilco
