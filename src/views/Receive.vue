<template>
<div id="receive">
    <div class="container-menu">
        <v-card class="default-menu">
            <h1>SYNC INFO</h1>
            <v-divider class="espaced"></v-divider>
            <div class="div-space">
                <label>{{ walletTopoHeight }} / {{ daemonTopoHeight }}</label>
                <v-progress-linear class="espaced" :color="(daemonTopoHeight / walletTopoHeight) == 1 ? 'green' : 'red'" :value="walletTopoHeight > 0 ? (daemonTopoHeight / walletTopoHeight).toFixed(0) * 100 : 100" height="17">
                    <strong>{{ walletTopoHeight > 0 ? (daemonTopoHeight / walletTopoHeight).toFixed(2) * 100 : 100 }} %</strong>
                </v-progress-linear>
            </div>
            <div class="div-space">
                <h4>Daemon:</h4>
                <span>{{ daemonAddress }}</span>
            </div>
        </v-card>
        <v-card color="grey darken-2" class="default-menu">
            <h4 class="wallet-name">MAIN WALLET</h4>
            <h2>{{ totalBalance }} DERO</h2>
            <span>{{ totalBalance * 0.4 }} €</span>
            <v-divider class="div-space"></v-divider>
            <span class="div-space">{{ address }}</span>
        </v-card>
        <v-card class="default-menu">
            <h1>BALANCE</h1>
            <v-divider class="div-space div-bott"></v-divider>
            <div>
                <h4 class="espaced">Total: {{ totalBalance }}</h4>
                <h4 class="espaced">Locked: {{ lockedBalance }}</h4>
                <h4 class="espaced">Unlocked: {{ unlockedBalance }}</h4>
            </div>
        </v-card>
    </div>
    <div class="container-menu">
        <v-card class="default-menu second-menu chart" v-if="chartReady">
            <h1>Price Chart</h1>
            <apexchart type="line" :options="priceChart.options" :series="priceChart.datas"></apexchart>
        </v-card>
        <v-card color="grey darken-2" class="default-menu second-menu" elevation="10">
            <h1>Last Transactions</h1>
            <v-divider class="div-space div-bott"></v-divider>
        </v-card>
    </div>
</div>
</template>

<script>
import * as chart from '../charts'
import VueApexCharts from 'vue-apexcharts'
import * as wallet from '../wallet/wallet'

export default {
    components: {
        apexchart: VueApexCharts,
    },
    data() {
        return {
            priceChart: {},
            chartReady: false,
            totalBalance: 0,
            lockedBalance: 0,
            unlockedBalance: 0,
            address: "Loading...",
            daemonTopoHeight: 0,
            walletTopoHeight: 0,
            daemonAddress: "No daemon"
        }
    },
    async mounted() {
        await wallet.waitWASM()
        
        setInterval(() => {
            this.updateInfos()
        }, 1000) //every 100ms

        chart.priceChart().then(data => {
            this.priceChart = data
            this.chartReady = true
        })
    },
    methods: {
        updateInfos() {
            /* eslint-disable */
            let infos = wallet.getInfos()
            this.totalBalance = infos.TotalBalance
            this.lockedBalance = infos.LockedBalance
            this.unlockedBalance = infos.UnlockedBalance
            this.address = infos.WalletAddress
            this.daemonTopoHeight = infos.DaemonTopoHeight
            this.walletTopoHeight = infos.WalletTopoHeight
            this.daemonAddress = infos.WalletDaemonAddress
        }
    }
}
</script>

<style scoped>
.container-menu  {
    margin-top: 5%;
    margin-left: 2%;
    margin-right: 2%;
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    word-break: break-all;
    align-items: stretch;
}

.default-menu {
    width: 100%;
    margin-bottom: 3%;
    padding: 3%;
}

.second-menu {
    width: auto;
}

.chart {
    flex: 1;
    padding: 0%;
    padding-top: 1%;
}

.div-space {
    margin: 3%;
}

.espaced {
    margin-top: 2%;
    margin-bottom: 2%;
}

@media screen and (max-width: 960px) {
    .container-menu  {
        flex-direction: column;
    }
}
</style>