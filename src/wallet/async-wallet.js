import PromiseWorker from 'promise-worker'
import Worker from 'worker-loader!./worker'
import { EventBus } from '../event-bus';

/* eslint-disable */
const worker = new Worker()
const promiseWorker = new PromiseWorker(worker)

let hasWallet = false
let walletName = "No Wallet"
let wasmRunning = false
let secretid

//from original dero web wallet (wallet.js)
let getRandomBytes = ((typeof self !== 'undefined' && (self.crypto || self.msCrypto)) ? function() { // Browsers
    let crypto = (self.crypto || self.msCrypto), QUOTA = 65536
    return function(n) {
        let a = new Uint8Array(n)
        for (var i = 0; i < n; i += QUOTA)
        {
            crypto.getRandomValues(a.subarray(i, i + Math.min(n - i, QUOTA)))
        }
        return a;
    }
} : function() { // Node
    return require("crypto").randomBytes;
}
)();

function toHexString(byteArray) {
    return byteArray.reduce((output, elem) => (output + ('0' + elem.toString(16)).slice(-2)), '')
}
//end


function run(name, ...params) {
    return promiseWorker.postMessage({
        id: secretid,
        name,
        params
    })
}

export async function init()
{
    if (wasmRunning)
        return

    secretid = toHexString(getRandomBytes(20))
    wasmRunning = await run("start")

    EventBus.$on('closeWallet', async value => {
        if (value) {
            closeWallet()
        }
    })

    //setInterval(() => run("ping"), 50) //we ping it every 50ms to prevent 'Go program has already exited'
}

export async function createWallet(walletName, password)
{
    let result = await run("createWallet", password)
    if (result === "success")
    {
        setWalletName(walletName)
        await onlineMode(true)
        let walletDump = await dumpEncryptedWallet()
       
        addEncryptedWallet(walletName, walletDump)
    }

    return result
}

export async function openEncryptedWallet(walletName, password, wallet_data)
{
    let result = await run("openEncryptedWallet", password, wallet_data)
    result = result === "success"

    if (result)
    {
        setWalletName(walletName)
        await onlineMode(true)
    }

    return result
}

export async function recoverWalletSeed(walletName, password, seed)
{
    let result = await run("recoverWalletSeed", password, seed)
    console.log(result)
    if (result === "success") {
        setWalletName(walletName)
        await onlineMode(true)
        let walletDump = await dumpEncryptedWallet()

        addEncryptedWallet(walletName, walletDump)
    }

    return result
}

export async function recoverViewWallet(walletName, password, viewKey)
{
    await run("recoverViewWallet", password, viewKey) //TODO
}

export async function getInfos()
{
    return await run("getInfos")
}

export async function generateIntegratedAddress()
{
    return await run("generateIntegratedAddress")
}

export async function getSeedInLanguage(language)
{
    return await run("getSeedInLanguage", language)
}

export async function getTxHistory()
{
    return await run("getTxHistory")
}

export async function transfer(addresses, amounts, paymentId)
{
    return await run("transfer", addresses, amounts, paymentId)
}

export async function dumpEncryptedWallet()
{
    return await run("dumpEncryptedWallet")
}

export async function onlineMode(online)
{
    return await run("onlineMode", online)
}

export async function closeWallet()
{
    await onlineMode(false)
    let walletDump = await dumpEncryptedWallet()
    updateEncryptedWallet(walletName, walletDump)
    /*let result = */await run("closeWallet")
    hasWallet = false
    walletName = "No Wallet"
    //console.log("Encrypted Wallet closed: " + result)
}

/* END WORKER */

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

export function updateEncryptedWallet(walletName, walletDump)
{
    removeEncryptedWallet(walletName)
    addEncryptedWallet(walletName, walletDump)
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