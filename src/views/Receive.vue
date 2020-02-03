<template>
<div id="receive">
    <div class="container-menu">
        <v-card class="default-menu">
            <h1>SYNC INFO</h1>
            <v-divider class="espaced"></v-divider>
            <div class="div-space">
                <label>5736723 / 5736723</label>
                <v-progress-linear class="espaced" color="green" value="100" height="17">
                    <strong>100%</strong>
                </v-progress-linear>
            </div>
            <div class="div-space">
                <h4>Daemon:</h4>
                <span>https://wallet.dero.io:443</span>
            </div>
        </v-card>
        <v-card color="grey darken-2" class="default-menu">
            <h4 class="wallet-name">MAIN WALLET</h4>
            <h2>1024 DERO</h2>
            <span>{{ 1024 * 0.4 }} €</span>
            <v-divider class="div-space"></v-divider>
            <span class="div-space">dERokevAZEZVJ2N7o39VH81BXBqX9ojtncnPTDMyiVbmYiTXQY93AUCLcor9xsWCKWhYy25ja89ikZWXWab9kXRB7LYfUmbQyS</span>
        </v-card>
        <v-card class="default-menu">
            <h1>BALANCE</h1>
            <v-divider class="div-space div-bott"></v-divider>
            <div>
                <h4 class="espaced">Total: 1024</h4>
                <h4 class="espaced">Locked: 0</h4>
                <h4 class="espaced">Unlocked: 0</h4>
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

export default {
    components: {
        apexchart: VueApexCharts,
    },
    data() {
        return {
            priceChart: {},
            chartReady: false
        }
    },
    mounted() {
        chart.priceChart().then(data => {
            this.priceChart = data
            this.chartReady = true
        })
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