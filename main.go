//Original version is here: https://git.dero.io/DeroProject/derosuite_webwallet/src/branch/master/cmd/webwallet/main.go

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"
	"time"
	"runtime/debug"
	"strings"
)
import "github.com/romana/rlog"
import "github.com/deroproject/derosuite/walletapi"
import "github.com/deroproject/derosuite/globals"
import "github.com/deroproject/derosuite/config"
import "github.com/deroproject/derosuite/address"
import "github.com/deroproject/derosuite/transaction"
import "github.com/deroproject/derosuite/crypto"

var miner_tx bool = false

var Local_wallet_instance *walletapi.Wallet

func register_wallet_callbacks() {

	js.Global().Set("DERO_JS_Create_New_Wallet", js.NewCallback(js_Create_New_Wallet))

	js.Global().Set("DERO_JS_Create_Encrypted_Wallet_From_Recovery_Words", js.NewCallback(js_Create_Encrypted_Wallet_From_Recovery_Words))
	js.Global().Set("DERO_JS_Open_Encrypted_Wallet", js.NewCallback(js_Open_Encrypted_Wallet))
	js.Global().Set("DERO_JS_Create_Wallet", js.NewCallback(js_Create_Wallet))
	js.Global().Set("DERO_JS_Create_Encrypted_Wallet_ViewOnly", js.NewCallback(js_Create_Encrypted_Wallet_ViewOnly))
	js.Global().Set("DERO_JS_GenerateIAddress", js.NewCallback(js_GenerateIAddress))

	js_GetSeedinLanguage := func(params []js.Value) {
		seed := "Some error occurred"
		if Local_wallet_instance != nil && len(params) == 1 {
			seed = Local_wallet_instance.GetSeedinLanguage(params[0].String())
		}
		js.Global().Set("wallet_seed", seed)
	}
	js.Global().Set("DERO_JS_GetSeedinLanguage", js.NewCallback(js_GetSeedinLanguage))

	js_TX_history := func(params []js.Value) {
		go func() {
			error_message := "Wallet is Closed"
			var buffer []byte
			var err error

			defer func() {
				js.Global().Set("tx_history", string(buffer))
				js.Global().Set("error_message", error_message)
			}()

			if Local_wallet_instance != nil {

				min_height, _ := strconv.ParseUint(params[6].String(), 0, 64)
				max_height, _ := strconv.ParseUint(params[7].String(), 0, 64)

				entries := Local_wallet_instance.Show_Transfers(params[0].Bool(), params[1].Bool(), params[2].Bool(), params[3].Bool(), params[4].Bool(), params[5].Bool(), min_height, max_height)

				if len(entries) == 0 {
					return
				}
				buffer, err = json.Marshal(entries)
				if err != nil {
					error_message = err.Error()
					return
				}
			}

		}()
	}
	js.Global().Set("DERO_JS_TX_History", js.NewCallback(js_TX_history))

	js_Transfer2 := func(params []js.Value) {
		transfer_error := "error"
		var transfer_txid, transfer_txhex, transfer_fee, transfer_amount, transfer_inputs_sum, transfer_change string

		defer func() {
			rlog.Warnf("setting values of tranfer variables")
			js.Global().Set("transfer_txid", transfer_txid)
			js.Global().Set("transfer_txhex", transfer_txhex)
			js.Global().Set("transfer_amount", transfer_amount)
			js.Global().Set("transfer_fee", transfer_fee)
			js.Global().Set("transfer_inputs_sum", transfer_inputs_sum)
			js.Global().Set("transfer_change", transfer_change)
			js.Global().Set("transfer_error", transfer_error)
			rlog.Warnf("setting values of tranfesr variables %s ", transfer_error)
		}()

		var address_list []address.Address
		var amount_list []uint64

		if params[0].Length() != params[1].Length() {

			return
		}

		for i := 0; i < params[0].Length(); i++ { // convert string address to our native form
			a, err := globals.ParseValidateAddress(params[0].Index(i).String())
			if err != nil {
				rlog.Warnf("Parsing address failed %s err %s\n", params[0].Index(i).String(), err)
				transfer_error = err.Error()
				return
			}
			address_list = append(address_list, *a)
		}

		for i := 0; i < params[1].Length(); i++ { // convert string address to our native form
			amount, err := globals.ParseAmount(params[1].Index(i).String())
			if err != nil {
				rlog.Warnf("Parsing address failed %s err %s\n", params[0].Index(i).String(), err)
				transfer_error = err.Error()
				return
				//return nil, jsonrpc.ErrInvalidParams()
			}

			amount_list = append(amount_list, amount)
		}

		payment_id := params[2].String()

		if len(payment_id) > 0 && !(len(payment_id) == 64 || len(payment_id) == 16) {
			transfer_error = "Invalid payment ID"
			return // we should give invalid payment ID
		}
		if _, err := hex.DecodeString(payment_id); err != nil {
			transfer_error = "Invalid payment ID"
			return // we should give invalid payment ID
		}

		unlock_time := uint64(0)
		fees_per_kb := uint64(0)
		mixin := uint64(0)

		tx, inputs, input_sum, change, err := Local_wallet_instance.Transfer(address_list, amount_list, unlock_time, payment_id, fees_per_kb, mixin)
		_ = inputs
		if err != nil {
			rlog.Warnf("Error while building Transaction err %s\n", err)
			transfer_error = err.Error()
			return
			//return nil, &jsonrpc.Error{Code: -2, Message: fmt.Sprintf("Error while building Transaction err %s", err)}

		}

		rlog.Infof("Inputs Selected for %s \n", globals.FormatMoney(input_sum))
		amount := uint64(0)
		for i := range amount_list {
			amount += amount_list[i]
		}
		rlog.Infof("Transfering total amount %s \n", globals.FormatMoney(amount))
		rlog.Infof("change amount ( will come back ) %s \n", globals.FormatMoney(change))
		rlog.Infof("fees %s \n", globals.FormatMoney(tx.RctSignature.Get_TX_Fee()))

		rlog.Infof(" size of tx %d", len(hex.EncodeToString(tx.Serialize())))

		transfer_fee = globals.FormatMoney12(tx.RctSignature.Get_TX_Fee())
		transfer_amount = globals.FormatMoney12(amount)
		transfer_change = globals.FormatMoney12(change)
		transfer_inputs_sum = globals.FormatMoney12(input_sum)
		transfer_txid = tx.GetHash().String()
		transfer_txhex = hex.EncodeToString(tx.Serialize())
		transfer_error = "success"
	}

	js_Transfer := func(params []js.Value) {
		go js_Transfer2(params)
	}
	js.Global().Set("DERO_JS_Transfer", js.NewCallback(js_Transfer))

	js_Transfer_Everything2 := func(params []js.Value) {
		transfer_error := "error"
		var transfer_txid, transfer_txhex, transfer_fee, transfer_amount, transfer_inputs_sum, transfer_change string

		defer func() {
			rlog.Warnf("setting values of tranfer variables")
			js.Global().Set("transfer_txid", transfer_txid)
			js.Global().Set("transfer_txhex", transfer_txhex)
			js.Global().Set("transfer_amount", transfer_amount)
			js.Global().Set("transfer_fee", transfer_fee)
			js.Global().Set("transfer_inputs_sum", transfer_inputs_sum)
			js.Global().Set("transfer_change", transfer_change)
			js.Global().Set("transfer_error", transfer_error)
			rlog.Warnf("setting values of tranfesr variables %s ", transfer_error)
		}()

		var address_list []address.Address
		var amount_list []uint64

		if params[0].Length() != 1 {
			return
		}

		for i := 0; i < params[0].Length(); i++ { // convert string address to our native form
			a, err := globals.ParseValidateAddress(params[0].Index(i).String())
			if err != nil {
				rlog.Warnf("Parsing address failed %s err %s\n", params[0].Index(i).String(), err)
				transfer_error = err.Error()
				return
				//return nil, jsonrpc.ErrInvalidParams()
			}
			address_list = append(address_list, *a)
		}

		payment_id := params[1].String()

		if len(payment_id) > 0 && !(len(payment_id) == 64 || len(payment_id) == 16) {
			transfer_error = "Invalid payment ID"
			return // we should give invalid payment ID
		}
		if _, err := hex.DecodeString(payment_id); err != nil {
			transfer_error = "Invalid payment ID"
			return // we should give invalid payment ID
		}

		//unlock_time := uint64(0)
		fees_per_kb := uint64(0)
		mixin := uint64(0)

		tx, inputs, input_sum, err := Local_wallet_instance.Transfer_Everything(address_list[0], payment_id, 0, fees_per_kb, mixin)
		_ = inputs
		if err != nil {
			rlog.Warnf("Error while building Everything Transaction err %s\n", err)
			transfer_error = err.Error()
			return
			//return nil, &jsonrpc.Error{Code: -2, Message: fmt.Sprintf("Error while building Transaction err %s", err)}

		}

		rlog.Infof("Inputs Selected for %s \n", globals.FormatMoney(input_sum))
		amount := uint64(0)
		for i := range amount_list {
			amount += amount_list[i]
		}
		amount = uint64(input_sum - tx.RctSignature.Get_TX_Fee())
		change := uint64(0)
		rlog.Infof("Transfering everything total amount %s \n", globals.FormatMoney(amount))
		rlog.Infof("change amount ( will come back ) %s \n", globals.FormatMoney(change))
		rlog.Infof("fees %s \n", globals.FormatMoney(tx.RctSignature.Get_TX_Fee()))

		rlog.Infof(" size of tx %d", len(hex.EncodeToString(tx.Serialize())))

		transfer_fee = globals.FormatMoney12(tx.RctSignature.Get_TX_Fee())
		transfer_amount = globals.FormatMoney12(amount)
		transfer_change = globals.FormatMoney12(change)
		transfer_inputs_sum = globals.FormatMoney12(input_sum)
		transfer_txid = tx.GetHash().String()
		transfer_txhex = hex.EncodeToString(tx.Serialize())
		transfer_error = "success"
	}

	js_Transfer_Everything := func(params []js.Value) {
		go js_Transfer_Everything2(params)
	}
	js.Global().Set("DERO_JS_Transfer_Everything", js.NewCallback(js_Transfer_Everything))

	js_Relay_TX2 := func(params []js.Value) {
		hex_tx := strings.TrimSpace(params[0].String())
		rlog.Warnf("tx decoding  hex")
		tx_bytes, err := hex.DecodeString(hex_tx)
		rlog.Warnf("tx decoding hex err %s", err)

		if err != nil {
			js.Global().Set("relay_error", fmt.Sprintf("Transaction Could NOT be hex decoded err %s", err))
			return
		}

		var tx transaction.Transaction

		err = tx.DeserializeHeader(tx_bytes)
		rlog.Warnf("tx decoding tx err %s", err)

		if err != nil {
			js.Global().Set("relay_error", fmt.Sprintf("Transaction Could NOT be deserialized err %s", err))
			return
		}

		err = Local_wallet_instance.SendTransaction(&tx) // relay tx to daemon/network
		rlog.Infof("tx relaying tx err %s", err)

		if err != nil {
			js.Global().Set("relay_error", fmt.Sprintf("Transaction sending failed txid = %s, err %s", tx.GetHash(), err))
			return
		}
		js.Global().Set("relay_error", "success")
	}

	js_Relay_TX := func(params []js.Value) {
		go js_Relay_TX2(params)
	}
	js.Global().Set("DERO_JS_Relay_TX", js.NewCallback(js_Relay_TX))

	js_Close_Encrypted_Wallet := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.Close_Encrypted_Wallet()
			Local_wallet_instance = nil

			fmt.Printf("Wallet has been closed\n")
		}

	}
	js.Global().Set("DERO_JS_Close_Encrypted_Wallet", js.NewCallback(js_Close_Encrypted_Wallet))

	// these function does NOT report back anything
	js_Rescan_Blockchain := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.Clean()               // clean existing data from wallet
			Local_wallet_instance.Rescan_From_Height(0) // we are setting it to zero, it will be automatically convert to start height if it is set
		}
	}
	js.Global().Set("DERO_JS_Rescan_Blockchain", js.NewCallback(js_Rescan_Blockchain))

	// this function does NOT report back anything
	js_SetOnline := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.SetOnlineMode()
		}
	}
	js.Global().Set("DERO_JS_SetOnline", js.NewCallback(js_SetOnline))

	// this function does NOT report back anything
	js_SetOffline := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.SetOfflineMode()
		}
	}
	js.Global().Set("DERO_JS_SetOffline", js.NewCallback(js_SetOffline))

	// this function does NOT report back anything
	js_ChangePassword := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.Set_Encrypted_Wallet_Password(params[0].String())
		}
	}
	js.Global().Set("DERO_JS_ChangePassword", js.NewCallback(js_ChangePassword))

	// this function does NOT report back anything
	js_SetInitialHeight := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.SetInitialHeight(int64(params[0].Int()))
		}
	}
	js.Global().Set("DERO_JS_SetInitialHeight", js.NewCallback(js_SetInitialHeight))

	// this function does NOT report back anything
	js_SetMixin := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.SetMixin((params[0].Int()))
		}
	}
	js.Global().Set("DERO_JS_SetMixin", js.NewCallback(js_SetMixin))

	// this function does NOT report back anything
	js_SetFeeMultiplier := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.SetFeeMultiplier(float32(params[0].Float()))
		}
	}
	js.Global().Set("DERO_JS_SetFeeMultiplier", js.NewCallback(js_SetFeeMultiplier))
        
        
        // this function does NOT report back anything
	js_SetSyncTime := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.SetDelaySync(int64(params[0].Int()))
		}
	}
	js.Global().Set("DERO_JS_SetSyncTime", js.NewCallback(js_SetSyncTime))

	// this function does NOT report back anything
	js_SetDaemonAddress := func(params []js.Value) {
		if Local_wallet_instance != nil {
			Local_wallet_instance.SetDaemonAddress(params[0].String())
		}
	}
	js.Global().Set("DERO_JS_SetDaemonAddress", js.NewCallback(js_SetDaemonAddress))
	// some apis to detect  parse validate address
	// this will setup some fields
	js_VerifyAddress := func(params []js.Value) {

		var address_main, address_paymentid, address_error string
		var address_valid, address_integrated bool

		address_error = "error"
		addr, err := globals.ParseValidateAddress(params[0].String())
		if err == nil {
			address_valid = true
			if addr.IsIntegratedAddress() {
				address_integrated = true
				address_paymentid = fmt.Sprintf("%x", addr.PaymentID)
			} else {
				address_integrated = false
			}
			address_error = "success"
		} else {
			address_error = err.Error()
			address_valid = false
			address_integrated = false
		}

		js.Global().Set("address_error", address_error)
		js.Global().Set("address_main", address_main)
		js.Global().Set("address_paymentid", address_paymentid)
		js.Global().Set("address_valid", address_valid)
		js.Global().Set("address_integrated", address_integrated)

	}

	js.Global().Set("DERO_JS_VerifyAddress", js.NewCallback(js_VerifyAddress))

	js_VerifyAmount := func(params []js.Value) {
		var amount_valid bool
		lamountstr := strings.TrimSpace(params[0].String())
		_, err := globals.ParseAmount(lamountstr)

		if err != nil {
			js.Global().Set("amount_valid", amount_valid)
			js.Global().Set("amount_error", err.Error())
			return
		}
		amount_valid = true
		js.Global().Set("amount_valid", amount_valid)
		js.Global().Set("amount_error", "success")
	}
	js.Global().Set("DERO_JS_VerifyAmount", js.NewCallback(js_VerifyAmount))

	js_VerifyPassword := func(params []js.Value) {
		password_error := "error"
		if Local_wallet_instance != nil {
			valid := Local_wallet_instance.Check_Password(params[0].String())
			if valid {
				password_error = "success"
			}
		}
		js.Global().Set("password_error", password_error)
	}
	js.Global().Set("DERO_JS_VerifyPassword", js.NewCallback(js_VerifyPassword))

	js_GetEncryptedCopy := func(params []js.Value) {
		wallet_encrypted_error := "error"
		var err error
		var encrypted_bytes []byte
		if Local_wallet_instance != nil {
			encrypted_bytes, err = Local_wallet_instance.GetEncryptedCopy()
			if err == nil {
				wallet_encrypted_error = "success"
			} else {
				wallet_encrypted_error = err.Error()
			}
		}

		typeu8array := js.TypedArrayOf(encrypted_bytes)
		js.Global().Set("wallet_encrypted_dump", typeu8array)
		typeu8array.Release()
		js.Global().Set("wallet_encrypted_error", wallet_encrypted_error)
	}

	js.Global().Set("DERO_JS_GetEncryptedCopy", js.NewCallback(js_GetEncryptedCopy))

}

// if this remain empty, default 127.0.0.1:20206 is used
var daemon_address = "" // this is setup below at runtime

// this wasm module exports necessary wallet apis to javascript
func main() {

	fmt.Printf("running go")
	globals.Arguments = map[string]interface{}{}

	globals.Arguments["--testnet"] = false

	globals.Config = config.Mainnet
	//globals.Initialize()

	debug.SetGCPercent(40) // start GC at 40%

	href := js.Global().Get("location").Get("href")
	u, err := url.Parse(href.String())
	if err == nil {
		r := strings.NewReplacer("0", "",
			"1", "",
			"2", "",
			"3", "",
			"4", "",
			"5", "",
			"6", "",
			"7", "",
			"8", "",
			"9", "",
			".", "",
			":", "",
		)
		rlog.Infof("u %+v", u)
		rlog.Infof("scheme %+v", u.Scheme)
		rlog.Infof("Host %+v", u.Host)
		if u.Scheme == "http" || u.Scheme == "" { // we do not support DNS names for http, for security reasons
			if len(r.Replace(u.Host)) == 0 { // number is an ipadress
				if strings.Contains(u.Host, ":") {
					daemon_address = u.Host // set the daemon address
				} else {
					daemon_address = u.Host + ":80" // set the daemon address
				}
			}
		} else if u.Scheme == "https" {
			if strings.Contains(u.Host, ":") {
				daemon_address = u.Scheme + "://" + u.Host // set the daemon address
			} else {
				daemon_address = u.Scheme + "://" + u.Host + ":443" // set the daemon address
			}
		}

		if len(daemon_address) == 0 {
			if globals.IsMainnet() {
				daemon_address = "127.0.0.1:20206"
			} else {
				daemon_address = "127.0.0.1:30306"
			}
		}

	}

	register_wallet_callbacks()
	go update_balance()

	select {} // if this return, program will exit
}

func update_balance() {
}

var i32_address, i32_address_paymentid string
var i8_address, i8_address_paymentid string

// generate integrated address at user demand
func generate_integrated_address() {
	if Local_wallet_instance != nil {

		i8 := Local_wallet_instance.GetRandomIAddress8()
		i32 := Local_wallet_instance.GetRandomIAddress32()

		js.Global().Set("random_i32_address", i32.String())
		js.Global().Set("random_i32_address_paymentid", fmt.Sprintf("%x", i32.PaymentID))

		js.Global().Set("random_i8_address", i8.String())
		js.Global().Set("random_i8_address_paymentid", fmt.Sprintf("%x", i8.PaymentID))

	}

}