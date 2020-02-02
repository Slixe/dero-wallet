<template>
<div id="main">
    <h1 class="title-page">RECOVER WALLET</h1>
    <v-card class="menu">
        <h2>Available Wallets</h2>
        <v-select :items="availableWallets" :color="$selectColor" :item-color="$selectColor" class="wallet-select"></v-select>
        <div class="select-buttons">
            <v-btn color="success" class="select-button" outlined>Open Selected Wallet</v-btn>
            <v-btn color="error" class="select-button" outlined>Delete Selected Wallet</v-btn>
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
export default {
    data() {
        return {
            availableWallets: ["Main Wallet", "Test Wallet", "Slixe"],
            recoverButtons: ["Recover using Seed", "Recover using Recovery Key", "Recover using File"]
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
    width: 35%;
}

.recover-view-wallet {
    margin-top: 6%;
    width: 35%;
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