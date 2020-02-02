<template>
<div id="receive">
    <v-card class="first-menu">
        <div class="first-side-menu">
            <h1>SYNC INFO</h1>
            <v-divider class="div-space div-bott"></v-divider>
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
        </div>
        <v-card color="grey darken-2" class="middle-menu">
            <h4 class="wallet-name">MAIN WALLET</h4>
            <h2>1024 DERO</h2>
            <span>{{ 1024 * 0.4 }} €</span>
            <v-divider class="div-space"></v-divider>
            <span>dERokevAZEZVJ2N7o39VH81BXBqX9ojtncnPTDMyiVbmYiTXQY93AUCLcor9xsWCKWhYy25ja89ikZWXWab9kXRB7LYfUmbQyS</span>
        </v-card>
        <div class="first-side-menu">
            <h1>BALANCE</h1>
            <v-divider class="div-space div-bott"></v-divider>
            <div>
                <h4 class="espaced">Total: 1024</h4>
                <h4 class="espaced">Locked: 0</h4>
                <h4 class="espaced">Unlocked: 0</h4>
            </div>
        </div>
    </v-card>
    <v-card class="chart" v-if="chartReady">
        <apexchart type="line" :options="priceChart.options" :series="priceChart.datas"></apexchart>
    </v-card>
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
.first-menu  {
    margin-left: 2%;
    margin-right: 2%;
    display: flex;
    flex-wrap: nowrap;
    flex-direction: row;
}

.first-side-menu {
    width: 25%;
    height: 100%;
    margin-left: auto;
    margin-right: auto;
    padding: 2%;
}

.middle-menu {
    padding: 2%;
    flex: 1;
    margin: auto;
}

.chart {
    width: 50%;
    margin: 2%;
    margin-top: 3%;
    margin-bottom: 3%;

    padding: 2%;
}

.wallet-name {
    margin-bottom: 2%;
}

.div-space {
    margin: 3%;
}

.div-bott {
    margin-bottom: 10%;
}

.espaced {
    margin-top: 2%;
    margin-bottom: 2%;
}

</style>