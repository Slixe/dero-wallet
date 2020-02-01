<template>
<div id="main">
    <h1>DERO WALLET</h1>
    <v-card dark class="menu">
        <h2>Available Wallets</h2>
        <v-select :items="availableWallets" class="wallet-select"></v-select>
        <div class="select-buttons">
            <v-btn color="success" outlined>Open Selected Wallet</v-btn>
            <v-btn color="error" outlined>Delete Selected Wallet</v-btn>
        </div>
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
    text-align: center;
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
</style>