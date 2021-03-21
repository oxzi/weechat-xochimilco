// SPDX-FileCopyrightText: 2021 Alvar Penning
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"C"

	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/oxzi/xochimilco"
)

var (
	identityKey ed25519.PrivateKey
	sessions    map[string]*xochimilco.Session
	msgAssembly map[string]string
	peerPubkey  map[string]string
)

// privmsg formats an IRC PRIVMSG with WeeChat's /quote command.
func privmsg(nick, msg string) string {
	return fmt.Sprintf("/quote PRIVMSG %s :%s", nick, msg)
}

// delPeer removes state for a closed peer.
func delPeer(nick string) {
	delete(sessions, nick)
	delete(msgAssembly, nick)
	delete(peerPubkey, nick)
}

//export xochimilco_init
func xochimilco_init() *C.char {
	var err error
	_, identityKey, err = ed25519.GenerateKey(nil)
	if err != nil {
		return C.CString(fmt.Sprintf("%v", err))
	}

	sessions = make(map[string]*xochimilco.Session)
	msgAssembly = make(map[string]string)
	peerPubkey = make(map[string]string)

	return nil
}

//export xochimilco_start
func xochimilco_start(name_buff *C.char) (msg, err_msg *C.char) {
	name := C.GoString(name_buff)

	sessions[name] = &xochimilco.Session{
		IdentityKey: identityKey,
		VerifyPeer: func(peer ed25519.PublicKey) (valid bool) {
			peerPubkey[name] = base64.StdEncoding.EncodeToString(peer)
			return true
		},
	}

	if msg, err := sessions[name].Offer(); err != nil {
		return nil, C.CString(fmt.Sprintf("%v", err))
	} else {
		return C.CString(privmsg(name, msg)), nil
	}
}

//export xochimilco_recv
func xochimilco_recv(privmsg_in *C.char) (has_ack_msg, is_established, is_closed, is_fragment bool, msg_out, key_self, key_peer, err_msg *C.char) {
	parts := strings.Split(C.GoString(privmsg_in), " ")
	if len(parts) < 4 {
		err_msg = C.CString("expected at least 4 parts")
		return
	} else if parts[1] != "PRIVMSG" {
		err_msg = C.CString("expected a PRIVMSG")
		return
	}

	name := parts[0]
	input := strings.Join(parts[3:], " ")[1:]

	re := regexp.MustCompile(`^:(\w+)!.*`)
	if matches := re.FindStringSubmatch(name); len(matches) != 2 {
		err_msg = C.CString("cannot extract nick")
		return
	} else {
		name = matches[1]
	}

	if prefix, ok := msgAssembly[name]; ok {
		input = prefix + input
		delete(msgAssembly, name)
	}

	if !strings.HasPrefix(input, xochimilco.Prefix) {
		return
	}

	if !strings.HasSuffix(input, xochimilco.Suffix) {
		msgAssembly[name] = input
		is_fragment = true
		return
	}

	sess, ok := sessions[name]
	if !ok {
		sess = &xochimilco.Session{
			IdentityKey: identityKey,
			VerifyPeer: func(peer ed25519.PublicKey) (valid bool) {
				peerPubkey[name] = base64.StdEncoding.EncodeToString(peer)
				return true
			},
		}

		ack, err := sess.Acknowledge(input)
		if err != nil {
			err_msg = C.CString(fmt.Sprintf("%v", err))
			return
		} else {
			sessions[name] = sess
			has_ack_msg = true
			msg_out = C.CString(privmsg(name, ack))
			key_self = C.CString(base64.StdEncoding.EncodeToString(identityKey.Public().(ed25519.PublicKey)))
			key_peer = C.CString(peerPubkey[name])
			return
		}
	}

	isEstablished, isClosed, plaintext, err := sess.Receive(input)
	if err != nil {
		err_msg = C.CString(fmt.Sprintf("%v", err))
		return
	} else if isEstablished {
		is_established = true
		key_self = C.CString(base64.StdEncoding.EncodeToString(identityKey.Public().(ed25519.PublicKey)))
		key_peer = C.CString(peerPubkey[name])
		return
	} else if isClosed {
		defer delPeer(name)
		is_closed = true
		return
	} else {
		msg_out = C.CString(fmt.Sprintf("%s PRIVMSG %s :%s", parts[0], parts[2], plaintext))
		return
	}
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
		return C.CString(privmsg(name, dataMsg)), nil
	}
}

//export xochimilco_stop
func xochimilco_stop(name_buff *C.char) (msg, err_msg *C.char) {
	name := C.GoString(name_buff)

	sess, ok := sessions[name]
	if !ok {
		return nil, C.CString("no such session")
	}
	defer delPeer(name)

	if msg, err := sess.Close(); err != nil {
		return nil, C.CString(fmt.Sprintf("%v", err))
	} else {
		return C.CString(privmsg(name, msg)), nil
	}
}

func main() {}
