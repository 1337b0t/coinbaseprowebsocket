# coinbaseprowebsocket

Quick and dirty websocket example based off github.com/preichenberger/go-coinbasepro

Highlights:

- TCP Server and Listener for calculations and processing
- Use this example to start your crytpo trading app, and save yourself a headache from unexpected disconnects

## Usage

```sh
$make build
$chmod +x run.sh
$bash run.sh
```

## Output

```sh
[{market:BTC-USD,price:18282.34,size:0.1,side: buy, bestBid:18279.05, bestAsk:18283.99, tm:1607627963}]
[{market:BTC-USD,price:18283.99,size:0.17330283,side: buy, bestBid:18279.05, bestAsk:18283.99, tm:1607627963}]
[{market:BTC-USD,price:18282.35,size:0.00380949,side: buy, bestBid:18279.76, bestAsk:18282.35, tm:1607627963}]
[{market:ETH-USD,price:559.99,size:0.85,side: sell, bestBid:559.99, bestAsk:560.00, tm:1607627963}]
[{market:BTC-USD,price:18281.08,size:0.05361368,side: buy, bestBid:18280.49, bestAsk:18281.08, tm:1607627963}]
[{market:ETH-USD,price:559.99,size:1.6,side: sell, bestBid:559.80, bestAsk:560.00, tm:1607627963}]
```

## To DO

- Technical Indicators
- K-Lines
