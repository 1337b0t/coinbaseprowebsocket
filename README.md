# coinbaseprowebsocket

Quick and dirty websocket example based off <https://github.com/preichenberger/go-coinbasepro>

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
[{market:BTC-USD,price:18282.34,size:0.1,side:buy,bestBid:18279.05,bestAsk:18283.99,tm:1607627963}]
[{market:BTC-USD,price:18283.99,size:0.17330283,side:buy,bestBid:18279.05,bestAsk:18283.99,tm:1607627963}]

```

## To DO

- Technical Indicators
- K-Lines
