<template>
<div id="new">
    <h1 class="title-page">CREATE NEW WALLET</h1>
    <v-card class="menu">
        <v-alert v-show="alertShow" :type="alertType">
            {{ alertMessage }}
        </v-alert>
        <span>Please fill in all the blanks</span>
        <v-divider class="new-div"></v-divider>
        <v-text-field autocomplete="new-password" v-model="walletName" label="Wallet name for quick reference" :color="$selectColor" filled :disabled="btnDisabled"></v-text-field>
        <v-text-field autocomplete="new-password" type="password" v-model="password" label="Password" :color="$selectColor" filled :disabled="btnDisabled"></v-text-field>
        <v-text-field autocomplete="new-password" type="password" label="Confirm password" :color="$selectColor" filled :disabled="btnDisabled"></v-text-field>
        <v-btn outlined class="create-wallet" @click.prevent="newWallet()" :disabled="btnDisabled">Create Wallet</v-btn>
    </v-card>
</div>
</template>

<script>
import * as wallet from '../wallet/wallet'

export default {
    data() {
        return {
            walletName: "",
            password: "",
            alertType: "error",
            alertShow: false,
            alertMessage: "An error as occured!",
            btnDisabled: false
        }
    },
    methods: {
        async newWallet() {
            this.btnDisabled = true
            let result = wallet.createWallet(this.password) //TODO async
            if (result === "success") {
                this.alertType = "success"
                this.alertMessage = "Wallet successfully created!"
                this.alertShow = true

                setTimeout(() => {
                    this.$router.push('/receive')
                } , 1500) //1.5s
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

    padding-top: 2%;
    padding-bottom: 2%;
    padding-left: 10%;
    padding-right: 10%;

    display: flex;
    flex-direction: column;
}

.new-label {
    text-align: center;
}

.create-wallet {
    width: 35%;
    margin: auto;
    margin-top: 6%;
}

.new-div {
    margin: 5%;
}

@media screen and (max-width: 960px) {
    .menu {
        margin-top: 3%;
        margin-left: 2%;
        margin-right: 2%;
        margin-bottom: 3%;

        padding-top: 5%;
        padding-bottom: 5%;
    }

    .create-wallet {
        width: auto;
    }

    .title-page {
        margin-top: 5%;
    }
}
</style>