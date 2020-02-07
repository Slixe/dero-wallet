// +build js,wasm

package main

import (
	"encoding/hex"
	"encoding/json"
	//"fmt"
	"syscall/js"
	//"time"
	//"strings"
)

import "github.com/deroproject/derosuite/walletapi"
import "github.com/deroproject/derosuite/globals"
import "github.com/deroproject/derosuite/config"
//import "github.com/deroproject/derosuite/address"
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
	/*js.Global().Set("DERO_", js.FuncOf())
	js.Global().Set("DERO_", js.FuncOf())
	js.Global().Set("DERO_", js.FuncOf())
	js.Global().Set("DERO_", js.FuncOf())*/

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

	return string(encrypted_bytes)
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