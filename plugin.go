// SPDX-FileCopyrightText: 2021 Alvar Penning
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"C"

	"crypto/ed25519"
	"fmt"
	"regexp"
	"strings"

	"github.com/oxzi/xochimilco"
)

var (
	identityKey ed25519.PrivateKey
	sessions    map[string]*xochimilco.Session
	msgAssembly map[string]string
)

//export xochimilco_init
func xochimilco_init() *C.char {
	var err error
	_, identityKey, err = ed25519.GenerateKey(nil)
	if err != nil {
		return C.CString(fmt.Sprintf("%v", err))
	}

	sessions = make(map[string]*xochimilco.Session)
	msgAssembly = make(map[string]string)
	return nil
}

//export xochimilco_start
func xochimilco_start(name_buff *C.char) (msg, err_msg *C.char) {
	name := C.GoString(name_buff)

	sessions[name] = &xochimilco.Session{
		IdentityKey: identityKey,
		VerifyPeer: func(peer ed25519.PublicKey) (valid bool) {
			// TODO
			return true
		},
	}

	if msg, err := sessions[name].Offer(); err != nil {
		return nil, C.CString(fmt.Sprintf("%v", err))
	} else {
		return C.CString(msg), nil
	}
}

//export xochimilco_recv
func xochimilco_recv(privmsg_in *C.char) (is_ack bool, msg_out, err_msg *C.char) {
	parts := strings.Split(C.GoString(privmsg_in), " ")
	if len(parts) < 4 {
		return false, nil, C.CString("expected at least 4 parts")
	} else if parts[1] != "PRIVMSG" {
		return false, nil, C.CString("expected a PRIVMSG")
	}

	name := parts[0]
	input := strings.Join(parts[3:], " ")[1:]

	re := regexp.MustCompile(`^:(\w+)!.*`)
	if matches := re.FindStringSubmatch(name); len(matches) != 2 {
		return false, nil, C.CString("cannot extract nick")
	} else {
		name = matches[1]
	}

	if prefix, ok := msgAssembly[name]; ok {
		input = prefix + input
		delete(msgAssembly, name)
	}

	if !strings.HasPrefix(input, xochimilco.Prefix) {
		return false, nil, nil
	}

	if !strings.HasSuffix(input, xochimilco.Suffix) {
		msgAssembly[name] = input
		return false, nil, nil
	}

	sess, ok := sessions[name]
	if !ok {
		sess = &xochimilco.Session{
			IdentityKey: identityKey,
			VerifyPeer: func(peer ed25519.PublicKey) (valid bool) {
				// TODO
				return true
			},
		}

		ack, err := sess.Acknowledge(input)
		if err != nil {
			return false, nil, C.CString(fmt.Sprintf("%v", err))
		} else {
			sessions[name] = sess
			return true, C.CString(ack), nil
		}
	}

	// TODO
	_, _, plaintext, err := sess.Receive(input)
	if err != nil {
		return false, nil, C.CString(fmt.Sprintf("%v", err))
	} else if plaintext != nil {
		out := fmt.Sprintf("%s PRIVMSG %s :%s", parts[0], parts[2], plaintext)
		return false, C.CString(out), nil
	}

	return false, nil, nil
}

//export xochimilco_send
func xochimilco_send(privmsg_in *C.char) (msg_out, err_msg *C.char) {
	parts := strings.Split(C.GoString(privmsg_in), " ")
	if len(parts) < 3 {
		return nil, C.CString("expected at least 3 parts")
	} else if parts[0] != "PRIVMSG" {
		return nil, C.CString("expected a PRIVMSG")
	}

	name := parts[1]
	input := strings.Join(parts[2:], " ")[1:]

	if strings.HasPrefix(input, xochimilco.Prefix) {
		return nil, nil
	}

	sess, ok := sessions[name]
	if !ok {
		return nil, nil
	}

	dataMsg, err := sess.Send([]byte(input))
	if err != nil {
		return nil, C.CString(fmt.Sprintf("%v", err))
	} else {
		return C.CString(fmt.Sprintf("PRIVMSG %s :%s", name, dataMsg)), nil
	}
}

//export xochimilco_stop
func xochimilco_stop(name_buff *C.char) (msg, err_msg *C.char) {
	name := C.GoString(name_buff)

	sess, ok := sessions[name]
	if !ok {
		return nil, C.CString("no such session")
	}
	defer delete(sessions, name)

	if msg, err := sess.Close(); err != nil {
		return nil, C.CString(fmt.Sprintf("%v", err))
	} else {
		return C.CString(msg), nil
	}
}

func main() {}
