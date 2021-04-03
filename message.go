package main

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	l.Printf("Message received. Type: %s\nContent: %v %s\n", m.Name, m.Payload, m.Payload)
	switch m.Name {
	case "getWallets":
		if payload, err = getWallets(); err != nil {
			payload = err.Error()
			return
		}
	case "getWallet":
		// Unmarshal payload
		var wallet getWalletRequest
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal([]byte(m.Payload), &wallet); err != nil {
				payload = err.Error()
				return
			}
		}

		if payload, err = getWallet(wallet.Owner, wallet.Passphrase); err != nil {
			payload = err.Error()
			return
		}
	case "createWallet": // Unmarshal payload
		var wallet createWalletRequest
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal([]byte(m.Payload), &wallet); err != nil {
				payload = err.Error()
				return
			}
		}

		if payload, err = createWallet(wallet.Owner, wallet.Passphrase); err != nil {
			payload = err.Error()
			return
		}
	case "addKeypair":
		var wallet createWalletRequest
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal([]byte(m.Payload), &wallet); err != nil {
				payload = err.Error()
				return
			}
		}

		if payload, err = createWallet(wallet.Owner, wallet.Passphrase); err != nil {
			payload = err.Error()
			return
		}
	case "taintkeypair":
		var kp taintKeypairRequest
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal([]byte(m.Payload), &kp); err != nil {
				payload = err.Error()
				return
			}
		}
		if payload, err = taintKeypair(kp.Owner, kp.Passphrase, kp.Pub); err != nil {
			payload = err.Error()
			return
		}
	case "signtransaction":
		var sign signTransactionRequest
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal([]byte(m.Payload), &sign); err != nil {
				payload = err.Error()
				return
			}
		}
		if payload, err = signTransaction(sign.Owner, sign.Passphrase, sign.Pub, sign.Message); err != nil {
			payload = err.Error()
			return
		}
	case "verifytransaction":
		var sign verifyTransactionRequest
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal([]byte(m.Payload), &sign); err != nil {
				payload = err.Error()
				return
			}
		}
		if payload, err = verifyTransaction(sign.Owner, sign.Passphrase, sign.Pub, sign.Message, sign.Signature); err != nil {
			payload = err.Error()
			return
		}
	case "getconfig":
		if payload, err = getConfig(); err != nil {
			payload = err.Error()
			return
		}
	case "startService":
		if payload, err = startService(); err != nil {
			payload = err.Error()
			return
		}
	case "stopService":
		if payload, err = stopService(); err != nil {
			payload = err.Error()
			return
		}
	}

	return
}
