const chartElement = document.getElementById("chart")

const chart = LightweightCharts.createChart(chartElement, {height: 700});
const candleSeries = chart.addCandlestickSeries({
    upColor: '#ffffff', downColor: '#000000', borderVisible: true,
    wickUpColor: '#000000', wickDownColor: '#000000', borderColor: '#000000'
})

chart.timeScale().fitContent();

async function fetchCandles() {
    let res = await fetch("/api/candles")
    let data = await res.json()

    candleData = []

    data.sort(function (a, b){
        return new Date(a["Timestamp"]) - new Date(b["Timestamp"])
    });
    
    for (let i = 0; i < data.length; i++) {
        let date = new Date(data[i]["Timestamp"]);

        if (date.getDay() == "0" || date.getDay() == "6") {
            continue
        }

        const day = ("0" + date.getDate()).slice(-2)

        const month = ("0" + (date.getMonth() + 1)).slice(-2)

        candleData.push({
            time:date.getFullYear()+"-"+month+"-"+day,            
            open:data[i]["Open"],
            high:data[i]["High"],
            low:data[i]["Low"],
            close:data[i]["Close"]
        })
    }

    console.log(candleData)

    candleSeries.setData(candleData)
}

fetchCandles()
