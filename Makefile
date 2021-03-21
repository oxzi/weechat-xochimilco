# SPDX-FileCopyrightText: 2021 Alvar Penning
#
# SPDX-License-Identifier: GPL-3.0-or-later

CC ?= gcc
GO  = go

.PHONY: install test-instance clean

xochimilco.so: plugin.o plugin.a
	$(CC) -shared -fPIC -Wl,-Bsymbolic -o xochimilco.so plugin.o plugin.a

plugin.o: plugin.c plugin.h
	$(CC) -fPIC -c -pthread -Wl,-Bsymbolic -Wall -Werror plugin.c plugin.h

plugin.a plugin.h: plugin.go
	$(GO) build -buildmode=c-archive plugin.go

install: xochimilco.so
	cp xochimilco.so ~/.weechat/plugins/xochimilco.so

test-instance: xochimilco.so
	weechat -t -r "/plugin load $(shell pwd)/xochimilco.so;/server add hackint irc.hackint.org/6697 -ssl;/set irc.server.hackint.nicks $(shell pwgen 6 1);/connect hackint"

clean:
	$(RM) plugin.{a,h,h.gch,o} xochimilco.so
