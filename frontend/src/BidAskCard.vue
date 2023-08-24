<script setup lang="ts">
    interface BidAsk {
        bid: number,
        ask: number,
    }

    interface Props {
        bidask: BidAsk,
        ticker: string,
    }
    
    const props = defineProps<Props>()
</script>

<template>
    <div class="container mx-auto shadow-lg shadow-pink-400/50 bg-pink-50 hover:shadow-2xl m-2 p-2">
        <h2 class="text-pink-400 text-lg font-bold">{{ticker}}</h2>
        <p class="">Bid: {{bidask.bid}}</p>
        <p class="">Ask: {{bidask.ask}}</p>
    </div>
</template>
