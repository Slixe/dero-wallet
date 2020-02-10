/* eslint-disable */
importScripts('./wasm_exec.js')
import registerPromiseWorker from 'promise-worker/register'

const go = new Go();
let wasmReady = false

let secretid

registerPromiseWorker(async (message) => {
    const id = message.id
    const name = message.name
    const params = message.params

    console.log(message)

    if (secretid != null && secretid != id)
    {
        return "Invalid secretId"
    }

    let result

    switch(name) {
        case "start": {
            await useWASM()
            secretid = id
            result = true
            break
        }
        case "ping": {
            result = GO_Pinger()
            break
        }
        case "createWallet": {
            result = DERO_CreateNewWallet("", params[0])
            break
        }
        case "openEncryptedWallet": {
            result = DERO_OpenEncryptedWallet("", params[0], params[1])
            break
        }
        case "recoverWalletSeed": {
            result = DERO_CreateEncryptedWalletFromRecoveryWords("", params[0], params[1])
            break
        }
        case "recoverViewWallet": {
            result = DERO_CreateEncryptedWalletViewOnly("", params[0], params[1])
            break
        }
        case "onlineMode": {
            DERO_OnlineMode(params[0])
            break
        }
        case "getInfos": {
            result = JSON.parse(DERO_GetInfos())
            break
        }
        case "generateIntegratedAddress": {
            result = JSON.parse(DERO_GenerateIntegratedAddress())
            break
        }
        case "getSeedInLanguage": {
            result = DERO_GetSeedInLanguage(params[0])
            break
        }
        case "getTxHistory": {
            result = new Promise((resolve, reject) => {
                DERO_Callback_GetTxHistory(true, true, true, true, true, false, 0, 0, (err, result) => {
                    if (err != "success") {
                        reject(err)
                    }
                    else {
                        resolve(JSON.parse(result))
                    }
                })
            })
            break
        }
        case "transfer": {
            result = new Promise((resolve, reject) => {
                DERO_Callback_Transfer(params[0], params[1], params[2], (err, result) => {
                    if (err) {
                        reject(err)
                    }
                    else {
                        resolve(JSON.parse(result))
                    }
                })
            })
            break
        }
        case "dumpEncryptedWallet": {
            result = DERO_DumpEncryptedWallet()
            break
        }
        case "closeWallet": {
            result = DERO_CloseWallet()
            break
        }
    }

    return result
})

async function useWASM()
{
    if (wasmReady) {
        console.log("WASM already ready!!")
    }

    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    let result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
    go.run(result.instance);
    console.log("running go");
    wasmReady = true
}