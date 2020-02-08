<template>
<div id="main">
    <h1 class="title-page">RECOVER WALLET</h1>
    <v-card class="menu">
        <h2>Available Wallets</h2>
        <v-select v-model="walletSelected" :items="availableWallets" item-text="name" :color="$selectColor" :item-color="$selectColor" class="wallet-select"></v-select>
        <div class="select-buttons">
            <v-btn color="success" class="select-button" @click="openWallet()" outlined>Open Selected Wallet</v-btn>
            <v-btn color="error" class="select-button" @click="deleteWallet()" outlined>Delete Selected Wallet</v-btn>
        </div>
        <v-divider class="div-space"></v-divider>
        <div class="recover-buttons">
            <v-btn outlined v-for="(button, key) in recoverButtons" :key="key" class="recover-button">{{ button }}</v-btn>
        </div>
        <v-btn outlined class="recover-view-wallet">Recover View Only Wallet</v-btn>
    </v-card>
</div>
</template>

<script>
import * as wallet from '../wallet/wallet'

export default {
    data() {
        return {
            walletSelected: null,
            availableWallets: [],
            recoverButtons: ["Recover using Seed", "Recover using Recovery Key", "Recover using File"]
        }
    },
    mounted() {
        this.availableWallets = wallet.getEncryptedWallets()
    },
    methods: {
        openWallet() {
            if (this.walletSelected != null)
            {
                let walletDump = wallet.getEncryptedWallet(this.walletSelected)
                if (wallet.openEncryptedWallet(prompt("Password for this wallet:"), walletDump)) {
                    this.$router.push('/receive')
                }
            }
        },
        deleteWallet() {
           if (this.walletSelected != null)
           {
               wallet.removeEncryptedWallet(this.walletSelected)
               this.availableWallets = wallet.getEncryptedWallets()
           } 
        }
    }
}
</script>

<style scoped>
.menu {
    margin-top: 3%;
    margin-left: 20%;
    margin-right: 20%;
    margin-bottom: 3%;

    padding: 2%;
}

.wallet-select {
    width: 50%;
    margin: auto;
}

.select-buttons {
    display: flex;
    justify-content: space-evenly;
}

.recover-buttons {
    margin-top: 8%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
}

.recover-button {
    margin-bottom: 2%;
    width: auto;
}

.recover-view-wallet {
    margin-top: 6%;
    width: auto;
}

.div-space {
    margin: 5%;
    margin-left: 15%;
    margin-right: 15%;
}

@media screen and (max-width: 960px)
{
    .menu {
        margin-top: 3%;
        margin-left: 2%;
        margin-right: 2%;
        margin-bottom: 3%;

        padding: 2%;
    }
    .select-buttons {
        display: flex;
        flex-direction: column;
        justify-content: space-around;
    }
    .recover-button {
        margin-bottom: 2%;
        width: auto;
    }
    .recover-view-wallet {
        margin-top: 6%;
        margin-bottom: 5%;
        width: auto;
    }
    .select-button {
        margin: auto;
        margin-bottom: 5%;
    }
    .title-page {
        margin-top: 5%;
    }
}
</style>