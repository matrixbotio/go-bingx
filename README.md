# BingX Go Client

This repository contains a Go client library for interacting with the BingX API. It provides a convenient way to access various trading functionalities, including spot trading and swap trading, using Go programming language.

## Features

- **Spot Trading**: 
  - Retrieve account balance
  - Create, cancel, and retrieve orders
  - Manage open orders and view order history

- **Swap Trading**:
  - Create orders with detailed parameters

## Installation

To use this library in your project, you can simply import it into your Go application. Ensure you have Go installed and set up on your machine.

```bash
go get github.com/Sagleft/go-bingx
```

## Usage

### Initialization

First, create a new client instance using your API key and secret key:

```go
import "github.com/Sagleft/go-bingx"

client := bingxgo.NewClient("your_api_key", "your_secret_key")
```

### Spot Trading

#### Get Account Balance

```go
spotClient := bingxgo.NewSpotClient(client)
balances, err := spotClient.GetBalance()
if err != nil {
    log.Fatal(err)
}
fmt.Println(balances)
```

#### Create Order

```go
order := bingxgo.SpotOrderRequest{
    Symbol:   "BTCUSDT",
    Side:     "BUY",
    Type:     "LIMIT",
    Quantity: 1.0,
    Price:    50000.0,
}

orderResponse, err := spotClient.CreateOrder(order)
if err != nil {
    log.Fatal(err)
}
fmt.Println(orderResponse)
```

#### Cancel Order

```go
err := spotClient.CancelOrder("BTCUSDT", "order_id")
if err != nil {
    log.Fatal(err)
}
```

### Swap Trading

#### Create Order

```go
tradeClient := bingxgo.TradeClient{client: client}
swapOrder := bingxgo.OrderRequest{
    Symbol:       "BTCUSDT",
    Side:         "BUY",
    PositionSide: "LONG",
    Type:         "LIMIT",
    Quantity:     1.0,
    Price:        50000.0,
}

swapOrderResponse, err := tradeClient.CreateOrder(swapOrder)
if err != nil {
    log.Fatal(err)
}
fmt.Println(swapOrderResponse)
```

## Contributing

We welcome contributions to this project. Please fork the repository and submit a pull request with your changes. Ensure your code follows the existing style and includes tests where applicable.

## Donate

In the comment to the transaction, indicate what kind of library to BingX, then I will understand that it needs to be further developed.

* USDT-Ton or TON: `UQD4otGPbsBePwdSL3PPnM6Qw_29AOpJnAcqOMOetZ-8OHHW`
* USDT-TRC20: `TQwTR3oJcFtXXynBoWsh6wNSM7id8K1k2U`

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

## Contact

For any questions or issues, please open an issue on the GitHub repository or contact the maintainers directly.
 