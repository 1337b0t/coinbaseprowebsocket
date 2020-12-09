# coinbaseprowebsocket
Quick and dirty websocket example based off github.com/preichenberger/go-coinbasepro

Highlights:

- TCP Server and Listener for calculations and processing
- Use this example to start your crytp trading app and save yourself a headache from unexpected disconnects

## Usage
```sh
$ make build
$ chmod +x run.sh
$ bash run.sh
```

## Output

```sh
[{market:BTC-USD,price:18449.94,size:0.01224668,side: sell, bestBid:18449.94, bestAsk:18449.95}]
[{market:ETH-USD,price:574.4,size:5.81068411,side: sell, bestBid:574.40, bestAsk:574.42}]
[{market:ETH-USD,price:574.41,size:0.42669127,side: buy, bestBid:574.40, bestAsk:574.41}]
[{market:BTC-USD,price:18449.95,size:0.00738907,side: buy, bestBid:18449.94, bestAsk:18449.95}]
[{market:BTC-USD,price:18449.95,size:0.00777901,side: buy, bestBid:18449.94, bestAsk:18449.95}]
```
