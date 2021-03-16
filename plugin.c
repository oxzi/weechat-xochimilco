// SPDX-FileCopyrightText: 2021 Alvar Penning
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <string.h>
#include "plugin.h"
#include "weechat-plugin.h"

WEECHAT_PLUGIN_NAME("xochimilco")
WEECHAT_PLUGIN_DESCRIPTION("E2E crypto for WeeChat")
WEECHAT_PLUGIN_AUTHOR("Alvar Penning <post@0x21.biz>")
WEECHAT_PLUGIN_VERSION("0.0.0")
WEECHAT_PLUGIN_LICENSE("GPL3")

struct t_weechat_plugin *weechat_plugin = NULL;

int hook_cmd(const void *pointer, void *data, struct t_gui_buffer *buffer,
             int argc, char **argv, char **argv_eol) {
  if (argc < 2) {
    weechat_printf(buffer, "%sxochimilco: Missing arguments.", weechat_prefix("error"));
    return WEECHAT_RC_ERROR;
  }

  if (strcmp("start", argv[1]) == 0) {
    weechat_printf(buffer, "start: %s", weechat_buffer_get_string(buffer, "name"));

    if (weechat_command(buffer, "hello") != WEECHAT_RC_OK) {
      return WEECHAT_RC_ERROR;
    }
  } else if (strcmp("stop", argv[1]) == 0) {
    weechat_printf(buffer, "stop: %s", weechat_buffer_get_string(buffer, "name"));

    if (weechat_command(buffer, "bye") != WEECHAT_RC_OK) {
      return WEECHAT_RC_ERROR;
    }
  } else {
    weechat_printf(buffer, "%sxochimilco: Unknown argument.", weechat_prefix("error"));
    return WEECHAT_RC_ERROR;
  }

  return cmdHandler(buffer, argc, argv);
}

/*
char *hook_privmsg_in(const void *pointer, void *data, const char *modifier,
                      const char *modifier_data, const char *string) {
  return NULL;
}

char *hook_privmsg_out(const void *pointer, void *data, const char *modifier,
                       const char *modifier_data, const char *string) {
  return NULL;
}
*/

int weechat_plugin_init(struct t_weechat_plugin *plugin, int argc,
                        char *argv[]) {
  weechat_plugin = plugin;

  weechat_hook_command(
      "xochimilco", "E2E crypto for WeeChat", "start | stop", "",
      "start %(nick) %(irc_servers) || stop %(nick) %(irc_servers)", &hook_cmd,
      NULL, NULL);

  // weechat_hook_modifier("irc_in2_privmsg", &hook_privmsg_in, NULL, NULL);
  // weechat_hook_modifier("irc_out1_privmsg", &hook_privmsg_out, NULL, NULL);

  return WEECHAT_RC_OK;
}

int weechat_plugin_end(struct t_weechat_plugin *plugin) {
  return WEECHAT_RC_OK;
}
