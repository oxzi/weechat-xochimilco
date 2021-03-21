// SPDX-FileCopyrightText: 2021 Alvar Penning
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include <stdlib.h>
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

  const char *buffer_name = weechat_buffer_get_string(buffer, "short_name");

  if (strcmp("start", argv[1]) == 0) {
    struct xochimilco_start_return offer = xochimilco_start((char*) buffer_name);
    if (offer.r1 != NULL) {
      weechat_printf(buffer, "%sxochimilco: Invalid Offer Message, %s", weechat_prefix("error"), offer.r1);
      return WEECHAT_RC_ERROR;
    }

    weechat_printf(buffer, "%sxochimilco: Sending Offer", weechat_prefix("action"));
    weechat_command(NULL, offer.r0);
  } else if (strcmp("stop", argv[1]) == 0) {
    struct xochimilco_stop_return close = xochimilco_stop((char*) buffer_name);
    if (close.r1 != NULL) {
      weechat_printf(buffer, "%sxochimilco: Invalid Close Message, %s", weechat_prefix("error"), close.r1);
      return WEECHAT_RC_ERROR;
    }

    weechat_printf(buffer, "%sxochimilco: Sending Close", weechat_prefix("action"));
    weechat_command(NULL, close.r0);
  } else {
    weechat_printf(buffer, "%sxochimilco: Unknown argument.", weechat_prefix("error"));
    return WEECHAT_RC_ERROR;
  }

  return WEECHAT_RC_OK;
}

char *hook_privmsg_in(const void *pointer, void *data, const char *modifier,
                      const char *modifier_data, const char *string) {
  struct t_gui_buffer *buffer = weechat_current_buffer();

  char *res = malloc(strlen(string)+1);
  strcpy(res, string);

  struct xochimilco_recv_return recv = xochimilco_recv((char*) string);
  if (recv.r2 != NULL) {
    weechat_printf(buffer, "%sxochimilco: Receiving error, %s", weechat_prefix("error"), recv.r2);
    return res;
  } else if (recv.r0) {
    weechat_printf(buffer, "%sxochimilco: Acknowledge conversation", weechat_prefix("action"));
    weechat_command(NULL, recv.r1);
    return NULL;
  } else if (!recv.r0 && recv.r1) {
    return recv.r1;
  } else {
    return res;
  }
}

char *hook_privmsg_out(const void *pointer, void *data, const char *modifier,
                       const char *modifier_data, const char *string) {
  struct t_gui_buffer *buffer = weechat_current_buffer();

  struct xochimilco_send_return send = xochimilco_send((char*) string);
  if (send.r1 != NULL) {
    weechat_printf(buffer, "%sxochimilco: Sending error, %s", weechat_prefix("error"), send.r1);
    return NULL;
  } else if (send.r0 != NULL) {
    weechat_command(NULL, send.r0);

    weechat_printf(buffer, "%s\t%s",
        weechat_string_eval_expression("${nick}", NULL, NULL, NULL),
        strchr(string, ':')+1);

    // returning NULL results in a retransmission
    char *res = malloc(1);
    strcpy(res, "");
    return res;
  } else {
    char *res = malloc(strlen(string)+1);
    strcpy(res, string);
    return res;
  }
}

int weechat_plugin_init(struct t_weechat_plugin *plugin, int argc, char *argv[]) {
  char *err_msg;
  if ((err_msg = xochimilco_init()) != NULL) {
    weechat_printf(NULL, "%sxochimilco: Initializing failed, %s", weechat_prefix("error"), err_msg);
    return WEECHAT_RC_ERROR;
  }

  weechat_plugin = plugin;

  weechat_hook_command(
      "xochimilco", "E2E crypto for WeeChat", "start | stop", "",
      "start %(nick) %(irc_servers) || stop %(nick) %(irc_servers)", &hook_cmd,
      NULL, NULL);

  weechat_hook_modifier("irc_in2_privmsg", &hook_privmsg_in, NULL, NULL);
  weechat_hook_modifier("irc_out1_privmsg", &hook_privmsg_out, NULL, NULL);

  return WEECHAT_RC_OK;
}

int weechat_plugin_end(struct t_weechat_plugin *plugin) {
  return WEECHAT_RC_OK;
}
