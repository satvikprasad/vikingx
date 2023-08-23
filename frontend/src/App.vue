<template>
    <header>
    </header>

    <main>
        {{balance}} USDT in account.
        <div v-for="[ticker, bidask] in bidasks">
            <h3>{{ticker}}</h3>
            Bid: {{bidask.bid}}
            Ask: {{bidask.ask}}
        </div>
    </main>
</template>

<style scoped>
</style>

<script lang="ts">
    import { ref } from 'vue'
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
                        console.log(data)
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
                        this.instruments = data.slice(0, 10)
                    })
            }
        },
        mounted: function() {
            this.bidAskTimer = setInterval(() => {
                this.fetchBidAsks()
            }, 2000)
        }
    }
</script>
