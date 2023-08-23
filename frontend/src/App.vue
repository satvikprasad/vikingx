<script setup lang="ts">
</script>

<template>
    <header>
        <div class="wrapper">
            <h1 class="white">VikingX.</h1>
        </div>
    </header>

    <main>
        <div v-if="render" class="wrapper">
            <p>{{ balance }}</p>
            <p>{{ time }}</p>
            <p>Bid: {{buy}}, Ask: {{sell}}</p>
        </div>
    </main>
</template>

<style scoped>
header {
    line-height: 1.5;
}

.logo {
    display: block;
    margin: 0 auto 2rem;
}

  @media (min-width: 1024px) {
      header {
          display: flex;
          place-items: center;
          padding-right: calc(var(--section-gap) / 2);
      }

      .logo {
          margin: 0 2rem 0 0;
      }

      header .wrapper {
          display: flex;
          place-items: flex-start;
          flex-wrap: wrap;
          flex-direction: column;
      }
  }
</style>

<script>
    export default {
        methods: {
            async fetchBalance() {
                try {
                    const response = await fetch('/api/balance')
                    const data = await response.json()
                    this.balance = data[0]
                    this.time = (new Date()).toString()
                } catch (error) {
                    console.error(error)
                }
                this.$forceUpdate()
            },
            async fetchTicker() {
                try {
                    const response = await fetch('/api/ticker')
                    const data = await response.json()
                    this.buy = data[0]
                    this.sell = data[1]
                    this.$forceUpdate()
                } catch (error) {
                    console.error(error)
                }
            }
        },
        mounted: function() {
            this.timer = setInterval(() => {
                this.fetchTicker() 
            }, 1000);
        },
        data() {
            return {
                balance: "",
                time: "",
                buy: 0,
                sell: 0,
                render: true,
                timer: 0,
            };
        },
        async created() {
            this.fetchTicker()
        },
    }


</script>
