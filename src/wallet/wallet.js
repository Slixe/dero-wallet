/* eslint-disable */
const go = new Go();
let wasmReady = false

function arrayBufferToHex(view) {
    let result = "";

    for (var i = 0; i < view.length; i++)
    {
        let value = view[i].toString(16)
        result += (value.length === 1 ? '0' + value : value)
    }

    return result
}

export async function useWASM()
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

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
 }

export async function waitWASM()
{
    while (!wasmReady)
    {
        console.log("Waiting WASM...")
        await sleep(150)
    }
}

export function pinger()
{
    setInterval(go_pinger(), 50); //we ping it every 50ms
}

export function dumpEncryptedWallet()
{
    let result = DERO_DumpEncryptedWallet()

    if (result.length > 500) {
        result = new Uint8Array(result)
    }


    return result
}

export function createWallet(walletName, password)
{
    let result = DERO_CreateNewWallet("", password)
    if (result === "success")
    {
        DERO_OnlineMode(true)
        let walletDump = dumpEncryptedWallet()
        localStorage.setItem("wallet_" + walletName, walletDump) //todo save encrypted wallet
    }

    return result
}

export function openEncryptedWallet(password, wallet_data)
{
    DERO_OpenEncryptedWallet("", password, wallet_data)
}

export function recoverWalletSeed(walletName, password, seed)
{
    DERO_CreateEncryptedWalletFromRecoveryWords("", password, seed)
}

export function recoverViewWallet(walletName, password, viewKey)
{
    DERO_CreateEncryptedWalletViewOnly("", password, viewKey)
}

export function getInfos()
{
    let result = DERO_GetInfos()
    return JSON.parse(result)
}