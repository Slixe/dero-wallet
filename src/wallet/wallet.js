import { EventBus } from '../event-bus';

/* eslint-disable */
const go = new Go();
let wasmReady = false

let hasWallet = false
let walletName = "No Wallet"

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

    EventBus.$on('closeWallet', closeWallet => {
        if (closeWallet) {
            DERO_OnlineMode(false)
            let result = DERO_CloseWallet()
            hasWallet = false
            walletName = "No Wallet"
            console.log("Encrypted Wallet closed: " + result)
        }
      })
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

    return result
}

export function getWalletsNames()
{
    let names = []
    for (const value of getEncryptedWallets())
    {
        names.push(value.name)
    }

    return names
}

export function getEncryptedWallets()
{
    return JSON.parse(localStorage.getItem("wallets")) || []
}

export function getEncryptedWallet(walletName)
{
    const wallets = getEncryptedWallets()
    let wallet = null
    let i = 0
    while (wallet == null && i < wallets.length)
    {
        let value = wallets[i]

        if (value.name === walletName)
            wallet = value.wallet
        i++
    }
    return wallet
}

export function addEncryptedWallet(walletName, walletDump)
{
    let walletJson = {name: walletName, wallet: walletDump }
    let wallets = getEncryptedWallets()
    wallets.push(walletJson)
    localStorage.setItem("wallets", JSON.stringify(wallets))
}

export function removeEncryptedWallet(walletName)
{
    let wallets = getEncryptedWallets()
    let found = false
    let i = 0
    while (!found && i < wallets.length)
    {
        let wallet = wallets[i]

        if (wallet.name === walletName)
        {
            wallets.splice(i, 1)
            found = true
        }

        i++
    }
    
    localStorage.setItem("wallets", JSON.stringify(wallets))
}

export function createWallet(walletName, password)
{
    let result = DERO_CreateNewWallet("", password)
    if (result === "success")
    {
        setWalletName(walletName)
        DERO_OnlineMode(true)
        let walletDump = dumpEncryptedWallet()
       
        addEncryptedWallet(walletName, walletDump)
    }

    return result
}

export function openEncryptedWallet(walletName, password, wallet_data)
{
    let result = DERO_OpenEncryptedWallet("", password, wallet_data)
    result = result === "success"

    if (result)
    {
        setWalletName(walletName)
        DERO_OnlineMode(true)
    }


    return result
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

export function setWalletName(name)
{
    walletName = name
    hasWallet = name != null

    EventBus.$emit('isWalletOpen', hasWallet)
}

export function getWalletName()
{
    return walletName
}

export function hasWalletOpen()
{
    return hasWallet
}