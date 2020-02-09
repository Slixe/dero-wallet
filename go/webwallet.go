// +build js,wasm

package main

import (
	"encoding/hex"
	"encoding/json"
	"syscall/js"
)

import "github.com/deroproject/derosuite/walletapi"
import "github.com/deroproject/derosuite/globals"
import "github.com/deroproject/derosuite/config"
import "github.com/deroproject/derosuite/address"
//import "github.com/deroproject/derosuite/transaction"
//import "github.com/deroproject/derosuite/crypto"

type walletInfo struct {
	Version string

	UnlockedBalance uint64
	LockedBalance uint64
	TotalBalance uint64

	DaemonHeight uint64
	DaemonTopoHeight uint64
	WalletDaemonAddress string

	WalletHeight uint64
	WalletTopoHeight uint64
	WalletInitialHeight int64
	WalletAddress string

	WalletAvailable bool
	WalletComplete bool
	WalletOnline bool
	WalletMixin int
	WalletFeesMultiplier float64
	WalletSyncTime int64
	WalletMinimumTopoHeight int64
}

type integratedAddress struct {
	Address string
	PaymentId string
}

type transferResult struct {
	TxId string
	TxHex string
	Fee string
	Amount string
	InputsSum string
	Change string
}

type addressValid struct {
	Valid bool
	Integrated bool
	Address string
	PaymentId string
	Err string
}

var daemon_address = "https://wallet.dero.io:443"
var Local_wallet_instance *walletapi.Wallet
var done = make(chan struct{})

func main() {
	globals.Arguments = map[string]interface{}{}
	globals.Arguments["--testnet"] = false
	globals.Config = config.Mainnet

	js.Global().Set("GO_Pinger", js.FuncOf(goPinger))

	js.Global().Set("DERO_CreateNewWallet", js.FuncOf(createNewWallet))
	js.Global().Set("DERO_CreateEncryptedWalletFromRecoveryWords", js.FuncOf(createEncryptedWalletFromRecoveryWords))
	js.Global().Set("DERO_OpenEncryptedWallet", js.FuncOf(openEncryptedWallet))
	js.Global().Set("DERO_CreateEncryptedWalletViewOnly", js.FuncOf(createEncryptedWalletViewOnly))
	js.Global().Set("DERO_GenerateIntegratedAddress", js.FuncOf(generateIntegratedAddress))
	js.Global().Set("DERO_GetInfos", js.FuncOf(getInfos))
	js.Global().Set("DERO_DumpEncryptedWallet", js.FuncOf(getEncryptedWalletDump))
	js.Global().Set("DERO_OnlineMode", js.FuncOf(setOnlineMode))
	js.Global().Set("DERO_CloseWallet", js.FuncOf(closeWallet))
	js.Global().Set("DERO_GetSeedInLanguage", js.FuncOf(getSeedInLanguage))
	js.Global().Set("DERO_Transfer", js.FuncOf(transfer))
	js.Global().Set("DERO_TransferEverything", js.FuncOf(transferEverything))
	js.Global().Set("DERO_VerifyPassword", js.FuncOf(verifyPassword))
	js.Global().Set("DERO_ChangePassword", js.FuncOf(changePassword))
	js.Global().Set("DERO_ValidateAddress", js.FuncOf(validateAddress))


	js.Global().Set("DERO_Callback_GetTxHistory", js.FuncOf(getTxHistory))
	js.Global().Set("DERO_Callback_CreateNewWallet", js.FuncOf(createNewWalletCallback))
	js.Global().Set("DERO_Callback_Transfer", js.FuncOf(transfer))

	<-done
}

func goPinger(this js.Value, inputs []js.Value) interface{} {
	return "OK"; //nothing
}

func createNewWallet(this js.Value, params []js.Value) interface{} {
	error_message := "error"
	filename := params[0].String()
	password := params[1].String()

	w, err := walletapi.Create_Encrypted_Wallet_Random(filename, password)

	if err == nil {
		error_message = "success"
		Local_wallet_instance = w
		Local_wallet_instance.SetDaemonAddress(daemon_address)
	} else {
		error_message = err.Error()
	}

	return error_message
}

func createNewWalletCallback(this js.Value, params []js.Value) interface{} {
	callback := params[len(params)-1:][0]

	go func() {
		error_message := "error"
		filename := params[0].String()
		password := params[1].String()
	
		w, err := walletapi.Create_Encrypted_Wallet_Random(filename, password)
	
		if err == nil {
			error_message = "success"
			Local_wallet_instance = w
			Local_wallet_instance.SetDaemonAddress(daemon_address)
			Local_wallet_instance.SetInitialHeight(0)
			Local_wallet_instance.Rescan_From_Height(0)
		} else {
			error_message = err.Error()
		}
	
		callback.Invoke(error_message)
	}()

	return nil
}

func createEncryptedWalletFromRecoveryWords(this js.Value, params []js.Value) interface{} {
	error_message := "error"

	w, err := walletapi.Create_Encrypted_Wallet_From_Recovery_Words(params[0].String(), params[1].String(), params[2].String())

	if err == nil {
		error_message = "success"
		Local_wallet_instance = w
		Local_wallet_instance.SetDaemonAddress(daemon_address)
	} else {
		error_message = err.Error()
	}

	return error_message
}

/*
 Params: filename, password, encrypted wallet in array
*/
func openEncryptedWallet(this js.Value, params []js.Value) interface{} { //here, js should send the encrypted wallet as param
	error_message := "error"

	src := []byte(params[2].String())
	db_array := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(db_array, src)
	db_array = db_array[:n]
	if err != nil {
		return err
		//rlog.Warnf("error decoding hex string \n", err)
	}

	//rlog.Infof("i passed DB of size %d\n", len(db_array))
	w, err := walletapi.Open_Encrypted_Wallet(params[0].String(), params[1].String(), db_array)
	if err == nil {
		error_message = "success"
		Local_wallet_instance = w
		Local_wallet_instance.SetDaemonAddress(daemon_address)
		//rlog.Infof("Successfully opened wallet\n")
	} else {
		error_message = err.Error()

		//rlog.Warnf("Error opened wallet %s\n", err)
	}
	return error_message
}

func createEncryptedWalletViewOnly(this js.Value, params []js.Value) interface{} {
	filename := params[0].String()
	password := params[1].String()
	viewkey := params[2].String()
	error_message := "error"

	wallet, err := walletapi.Create_Encrypted_Wallet_ViewOnly(filename, password, viewkey)

	if err != nil {
		error_message = err.Error()
	} else {
		error_message = "success"
		Local_wallet_instance = wallet
		Local_wallet_instance.SetDaemonAddress(daemon_address)
	}

	return error_message
}

func generateIntegratedAddress(this js.Value, params []js.Value) interface{} {
	var result integratedAddress

	if Local_wallet_instance != nil {
		i8 := Local_wallet_instance.GetRandomIAddress8()
		result.Address = i8.String()
		dst := make([]byte, hex.EncodedLen(len(i8.PaymentID)))
		hex.Encode(dst, i8.PaymentID)
		result.PaymentId = string(dst)
		//result.i32 = Local_wallet_instance.GetRandomIAddress32().String()
	}

	res, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return js.ValueOf(string(res))
}

func getInfos(this js.Value, params []js.Value) interface{} {

	var result walletInfo

	result.Version = config.Version.String()

	if Local_wallet_instance != nil {
		result.UnlockedBalance, result.LockedBalance = Local_wallet_instance.Get_Balance()
		result.TotalBalance = result.UnlockedBalance + result.LockedBalance

		result.WalletHeight = Local_wallet_instance.Get_Height()
		result.DaemonHeight = Local_wallet_instance.Get_Daemon_Height()
        result.WalletTopoHeight = uint64(Local_wallet_instance.Get_TopoHeight())
		result.DaemonTopoHeight = uint64(Local_wallet_instance.Get_Daemon_TopoHeight())
		
		result.WalletAddress = Local_wallet_instance.GetAddress().String()
		result.WalletAvailable = true
		result.WalletComplete = !Local_wallet_instance.Is_View_Only()

		result.WalletInitialHeight = Local_wallet_instance.GetInitialHeight()
		result.WalletOnline = Local_wallet_instance.GetMode()
		result.WalletMixin = Local_wallet_instance.GetMixin()
		result.WalletFeesMultiplier = float64(Local_wallet_instance.GetFeeMultiplier())
		result.WalletDaemonAddress = Local_wallet_instance.Daemon_Endpoint
		
		result.WalletSyncTime = Local_wallet_instance.SetDelaySync(0)
        result.WalletMinimumTopoHeight = Local_wallet_instance.GetMinimumTopoHeight()
	} else {
		result.UnlockedBalance = uint64(0)
		result.LockedBalance = uint64(0)
		result.TotalBalance = uint64(0)

		result.WalletHeight = uint64(0)
		result.DaemonHeight = uint64(0)
		result.WalletTopoHeight = uint64(0)
		result.DaemonTopoHeight = uint64(0)

		result.WalletInitialHeight = int64(0)

		result.WalletAddress = ""

		result.WalletAvailable = false
		result.WalletComplete = true
		result.WalletOnline = false
		result.WalletMixin = 5
		result.WalletFeesMultiplier = float64(1.5)
		result.WalletDaemonAddress = ""
		result.WalletSyncTime = int64(0)
		result.WalletMinimumTopoHeight = int64(-1)
	}

	res, err := json.Marshal(result)

	if err != nil {
		return err
	}

	return js.ValueOf(string(res))
}

func getEncryptedWalletDump(this js.Value, params []js.Value) interface{} {
	var encrypted_bytes []byte
	var err error
	if Local_wallet_instance != nil {
		encrypted_bytes, err = Local_wallet_instance.GetEncryptedCopy()

		if err != nil {
			return err.Error()
		} 
	}

	dst := make([]byte, hex.EncodedLen(len(encrypted_bytes)))
	hex.Encode(dst, encrypted_bytes)

	return string(dst)
}

func setOnlineMode(this js.Value, params []js.Value) interface{} {
	var currentState = false

	if Local_wallet_instance != nil {
		if (params[0].Bool()) {
			Local_wallet_instance.SetOnlineMode()
		} else {
			Local_wallet_instance.SetOfflineMode()
		}
		currentState = Local_wallet_instance.GetMode()
	}

	return currentState
}

func closeWallet(this js.Value, params []js.Value) interface{} {
	if Local_wallet_instance != nil {
		Local_wallet_instance.Close_Encrypted_Wallet()
		Local_wallet_instance = nil
	}

	return Local_wallet_instance == nil
}

func getSeedInLanguage(this js.Value, params []js.Value) interface{} {
	seed := "Some error occurred"
	if Local_wallet_instance != nil && len(params) == 1 {
		seed = Local_wallet_instance.GetSeedinLanguage(params[0].String())
	}

	return seed
}

func getTxHistory(this js.Value, params []js.Value) interface{} {
	go func() {
		callback := params[len(params)-1:][0] //js function passed as param for callback

		error_message := "Wallet is Closed"
		var buffer []byte
		var err error

		if Local_wallet_instance != nil {
			//min_height, _ := strconv.ParseUint(params[6].String(), 0, 64)
			//max_height, _ := strconv.ParseUint(params[7].String(), 0, 64)
			entries := Local_wallet_instance.Show_Transfers(params[0].Bool(), params[1].Bool(), params[2].Bool(), params[3].Bool(), params[4].Bool(), params[5].Bool(), 0, 0)

			if len(entries) != 0 {
				buffer, err = json.Marshal(entries)
				if err != nil {
					error_message = err.Error()
				} else {
					error_message = "success"
				}
			} else {
				error_message = "no entries"
			}
		}
		callback.Invoke(error_message, string(buffer)) //first param is error message, second is result as json 
	}()

	return nil
}

func transfer(this js.Value, params []js.Value) interface{} {
	go func() {
		callback := params[len(params)-1:][0]
		var result transferResult
		transfer_error := "error"

		defer func() {
			buffer, err := json.Marshal(result)
			if err != nil {
				transfer_error = err.Error()
			}

			callback.Invoke(transfer_error, string(buffer))
		}()

		var address_list []address.Address
		var amount_list []uint64

		if params[0].Length() != params[1].Length() {
			return
		}

		for i := 0; i < params[0].Length(); i++ { // convert string address to our native form
			a, err := globals.ParseValidateAddress(params[0].Index(i).String())
			if err != nil {
				transfer_error = err.Error()
				return
			}
			address_list = append(address_list, *a)
		}

		for i := 0; i < params[1].Length(); i++ { // convert string address to our native form
			amount, err := globals.ParseAmount(params[1].Index(i).String())
			if err != nil {
				transfer_error = err.Error()
				return
			}

			amount_list = append(amount_list, amount)
		}

		payment_id := params[2].String()

		if len(payment_id) > 0 && !(len(payment_id) == 64 || len(payment_id) == 16) {
			transfer_error = "Invalid payment ID"
			return
		}
		if _, err := hex.DecodeString(payment_id); err != nil {
			transfer_error = "Invalid payment ID"
			return
		}

		unlock_time := uint64(0)
		fees_per_kb := uint64(0)
		mixin := uint64(0)

		tx, inputs, input_sum, change, err := Local_wallet_instance.Transfer(address_list, amount_list, unlock_time, payment_id, fees_per_kb, mixin)
		_ = inputs
		if err != nil {
			transfer_error = err.Error()
			return

		}

		amount := uint64(0)
		for i := range amount_list {
			amount += amount_list[i]
		}

		result.Fee = globals.FormatMoney12(tx.RctSignature.Get_TX_Fee())
		result.Amount = globals.FormatMoney12(amount)
		result.Change = globals.FormatMoney12(change)
		result.InputsSum = globals.FormatMoney12(input_sum)
		result.TxId = tx.GetHash().String()
		result.TxHex = hex.EncodeToString(tx.Serialize())
		transfer_error = "success"
	}()

	return nil
}

func transferEverything(this js.Value, params []js.Value) interface{} {
	go func() {
		callback := params[len(params)-1:][0]
		var result transferResult
		transfer_error := "error"

		defer func() {
			buffer, err := json.Marshal(result)
			if err != nil {
				transfer_error = err.Error()
			}

			callback.Invoke(transfer_error, string(buffer))
		}()

		var address_list []address.Address
		var amount_list []uint64

		if params[0].Length() != 1 {
			return
		}

		for i := 0; i < params[0].Length(); i++ { // convert string address to our native form
			a, err := globals.ParseValidateAddress(params[0].Index(i).String())
			if err != nil {
				transfer_error = err.Error()
				return
			}
			address_list = append(address_list, *a)
		}

		payment_id := params[1].String()

		if len(payment_id) > 0 && !(len(payment_id) == 64 || len(payment_id) == 16) {
			transfer_error = "Invalid payment ID"
			return
		}
		if _, err := hex.DecodeString(payment_id); err != nil {
			transfer_error = "Invalid payment ID"
			return
		}

		//unlock_time := uint64(0)
		fees_per_kb := uint64(0)
		mixin := uint64(0)

		tx, inputs, input_sum, err := Local_wallet_instance.Transfer_Everything(address_list[0], payment_id, 0, fees_per_kb, mixin)
		_ = inputs
		if err != nil {
			transfer_error = err.Error()
			return

		}

		amount := uint64(0)
		for i := range amount_list {
			amount += amount_list[i]
		}
		amount = uint64(input_sum - tx.RctSignature.Get_TX_Fee())
		change := uint64(0)

		result.Fee = globals.FormatMoney12(tx.RctSignature.Get_TX_Fee())
		result.Amount = globals.FormatMoney12(amount)
		result.Change = globals.FormatMoney12(change)
		result.InputsSum = globals.FormatMoney12(input_sum)
		result.TxId = tx.GetHash().String()
		result.TxHex = hex.EncodeToString(tx.Serialize())
		transfer_error = "success"
	}()

	return nil
}

func verifyPassword(this js.Value, params []js.Value) interface{} {
	result := false
	if Local_wallet_instance != nil {
		result = Local_wallet_instance.Check_Password(params[0].String())
	}

	return result
}

func changePassword(this js.Value, params []js.Value) interface{} {
	result := false

	if Local_wallet_instance != nil {
		Local_wallet_instance.Set_Encrypted_Wallet_Password(params[0].String())
		result = true
	}

	return result
}

func validateAddress(this js.Value, params []js.Value) interface {} {
	var result addressValid

	addr, err := globals.ParseValidateAddress(params[0].String())

	if err == nil {
		result.Valid = true
		if addr.IsIntegratedAddress() {
			result.Integrated = true

			dst := make([]byte, hex.EncodedLen(len(addr.PaymentID)))
			hex.Encode(dst, addr.PaymentID)
			result.PaymentId = string(dst)

		} else {
			result.Integrated = false
		}
		result.Err = "success"
	} else {
		result.Valid = false
		result.Integrated = false
		result.Err = err.Error()
	}

	res, err := json.Marshal(result)
	if err != nil {
		result.Err = err.Error()
	}

	return string(res)
}