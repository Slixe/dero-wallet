<template>
<div id="new">
    <h1 class="title-page">CREATE NEW WALLET</h1>
    <v-card class="menu">
        <span>Please fill in all the blanks</span>
        <v-divider class="new-div"></v-divider>
        <v-text-field label="Wallet name for quick reference" :color="$selectColor" filled></v-text-field>
        <v-text-field label="Password" :color="$selectColor" filled></v-text-field>
        <v-text-field label="Confirm password" :color="$selectColor" filled></v-text-field>
        <v-btn outlined class="create-wallet">Create Wallet</v-btn>
    </v-card>
</div>
</template>

<script>
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