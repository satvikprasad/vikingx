const chartElement = document.getElementById("chart");

const chart = LightweightCharts.createChart(chartElement, {
    height: 700,
    width: 1000,
});

const candleSeries = chart.addCandlestickSeries({
    upColor: "#ffffff",
    downColor: "#000000",
    borderVisible: true,
    wickUpColor: "#000000",
    wickDownColor: "#000000",
    borderColor: "#000000",
});

const form = document.getElementById("symbol-form");
const input = document.getElementById("symbol");

form.addEventListener(
    "submit",
    (e) => {
        e.preventDefault();

        fetchCandles(input.value);
    },
    false
);

chart.timeScale().fitContent();

async function fetchCandles(symbol) {
    let res = await fetch("/api/candles/" + symbol);

    let data = await res.json();
    if (data.error) {
        console.log(data.error);
    }

    candleData = [];
    for (let i = 0; i < data.length; i++) {
        let date = new Date(data[i]["Timestamp"]);

        candleData.push({
            time: date / 1000,
            open: data[i]["Open"],
            high: data[i]["High"],
            low: data[i]["Low"],
            close: data[i]["Close"],
        });
    }

    candleSeries.setData(candleData);
}

const currentLocale = window.navigator.languages[0];
const formatter = Intl.NumberFormat(currentLocale, {
    style: "currency",
    currency: "USD",
}).format;

chart.applyOptions({
    rightPriceScale: {
        ticksVisible: true,
    },
    localization: {
        priceFormatter: formatter
    }
});

fetchCandles("BTCUSDT");
