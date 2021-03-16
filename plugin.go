// SPDX-FileCopyrightText: 2021 Alvar Penning
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

/*
#include "weechat-plugin.h"
*/
import "C"

import (
	_ "github.com/oxzi/xochimilco"
)

//export cmdHandler
func cmdHandler(buffer *C.struct_t_gui_buffer, argc C.int, argv **C.char) C.int {
	return C.WEECHAT_RC_OK
}

func main() {}
