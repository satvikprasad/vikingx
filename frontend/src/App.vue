<script lang="ts">
    import BidAskCard from './BidAskCard.vue'

    import { inject, ref } from 'vue'
    import type { Ref } from 'vue'

    interface BidAsk {
        bid: number,
        ask: number,
    }

    interface Instrument {
        instType: string
        instId: string
        instFamily: string
        uly: string
        category: string
        baseCcy: string
        quoteCcy: string
        settleCcy: string
        ctVal: string
        ctMult: string
        ctValCcy: string
        optType: string
        stk: string
        listTime: string
        expTime: string
        lever: string
        tickSz: string
        lotSz: string
        minSz: string
        ctType: string
        alias: string
        state: string
        maxLmtSz: string
        maxMktSz: string
        maxTwapSz: string
        maxIcebergSz: string
        maxTriggerSz: string
        maxStopSz: string
    }

    export default {
        setup() {
            const balance = ref('')
            const bidasks: Ref<Map<string, BidAsk>> = ref(new Map<string, BidAsk>())
            const bidAskTimer = 0
            const instruments: Ref<Instrument[]> = ref([])
            return {
                balance,
                bidasks,
                bidAskTimer,
                instruments
            }
        },
        async created() {
            await this.fetchBalance()
            await this.fetchInstruments()
                .then(async() => {
                    await this.fetchBidAsks()
                })
        },
        methods: {
            async fetchBalance() {
                await fetch('/api/balance')
                    .then(async (res) => {
                        const data = await res.json()
                        this.balance = data[0]
                    }).catch((error) => {return error})
            },
            async fetchBidAsks() {
                for (let i = 0; i < this.instruments.length; i++) {
        const instrument = this.instruments[i];
        await fetch('/api/bidask/' + instrument.instId)
        .then(async (res) => {
                            const data = await res.json()
                            console.log(data)
                            this.bidasks.set(instrument.instId, {
                                bid: data[0],
                                ask: data[1]
                            })
                        })
                }
            },
            async fetchInstruments() {
                await fetch('/api/instruments/SWAP')
                    .then(async (res) => {
                        const data: Instrument[] = await res.json()
                        this.instruments = data.filter((inst) => {
                            return inst.settleCcy == "USDT" || inst.baseCcy == "USDT" 
                        }).slice(0, 5)
                    })
            }
        },
        components: {
            BidAskCard
        },
        mounted: function() {
            this.bidAskTimer = setInterval(() => {
                this.fetchBidAsks()
            }, 2000)
        }
    }
</script>

<template>
    <header>
        <div class="container mx-auto">
            <h1 class="font-bold  text-xl text-center antialiased">{{balance}} USDT</h1>
        </div>
    </header>

    <main class="bg-gradient-to-b from-pink-200 to-pink-50">
        <div class="flex flex-row p-2 space-x-2">
            <BidAskCard v-for="[ticker, bidask] in bidasks" :ticker="ticker" :bidask="bidask" />
        </div>
    </main>
</template>

<style scoped>
</style>

