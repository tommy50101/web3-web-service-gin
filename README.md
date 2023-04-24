# Web-Service-Gin

## Instructions

取得已獲取之區塊、交易資訊

## Setup

```zsh
$ go get .
```

## Start

```zsh
$ go run .
```

## Usage

```
# 回傳最新的 n 個 blocks
[GET]http://localhost:8080/blocks?limit=n

// example
http://localhost:8080/blocks?limit=n

# 回傳單一 block by block id
[GET]http://localhost:8080/blocks/:id

// example
http://localhost:8080/blocks/36

# 回傳 transaction data with event logs
[GET]http://localhost:8080/transaction/:tx_hash

// example
http://localhost:8080/transaction/0xef4e0f4e7fd15bc84fb22471292e438b1ed390252f2809f94054c0a408330639
```
