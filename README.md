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

執行後，可透過輸入以選擇要獲取哪個鏈的資料，參數對應如下:

1:BSC testnet  2:Ethereum testnet(Goerli)  3:Ethereum mainnet

服務成功啟動後，可呼叫以下3支API

```
1. 回傳最新的 n 個 blocks
[GET]http://localhost:8080/blocks?limit=n


2. 回傳單一 block by block id ( 這邊的 id 為 Block Number )
[GET]http://localhost:8080/blocks/:id


3. 回傳 transaction data with event logs
[GET]http://localhost:8080/transaction/:tx_hash

```

以下為實際使用範例:

http://localhost:8080/blocks?limit=3

```JSON
{
    "blocks": [
        {
            "block_num": 29416071,
            "block_hash": "0xa8e04350f61bbc49eba595357b7c03a34689fbbb819e9dcb49354b869d143f37",
            "block_time": 1682940170,
            "parent_hash": "0xc21cf3190370359aa49298f256d8030119dd86d14535746e10292f4c0c8a4fdb"
        },
        {
            "block_num": 29416070,
            "block_hash": "0xc21cf3190370359aa49298f256d8030119dd86d14535746e10292f4c0c8a4fdb",
            "block_time": 1682940167,
            "parent_hash": "0x532c6341f7cd9fb8b03e593f398f189c55663ed156033797fe1690d1631f251b"
        },
        {
            "block_num": 29416069,
            "block_hash": "0x532c6341f7cd9fb8b03e593f398f189c55663ed156033797fe1690d1631f251b",
            "block_time": 1682940164,
            "parent_hash": "0xce390e89149143e07d41ad37ff8785d20fdcbec7dc7e1a1f7af2ee957d29c428"
        }
    ]
}

```

http://localhost:8080/blocks/29416069

```JSON
{
    "block_num": 29416069,
    "block_hash": "0x532c6341f7cd9fb8b03e593f398f189c55663ed156033797fe1690d1631f251b",
    "block_time": 1682940164,
    "parent_hash": "0xce390e89149143e07d41ad37ff8785d20fdcbec7dc7e1a1f7af2ee957d29c428",
    "transactions": [
        "0x4f6a80e9a00fdd3de6a0fd582ad1157f1b6752eed06a7a13171e9ceef899d462",
        "0x38eceecd7a121aa9d7a71b3a7a1d719592e1490ca6c6cff78663a82c56a7d55f",
        "0xb2212e28841b647ed40c1729bc3a970808089aa7f5ffbb3567572680a7077da3"
    ]
}

```

http://localhost:8080/transaction/0xb2212e28841b647ed40c1729bc3a970808089aa7f5ffbb3567572680a7077da3

```JSON
{
    "tx_hash": "0xb2212e28841b647ed40c1729bc3a970808089aa7f5ffbb3567572680a7077da3",
    "from": "0x35552c16704d214347f29Fa77f77DA6d75d7C752",
    "to": "0x0000000000000000000000000000000000001000",
    "nonce": 2652631,
    "data": "0xf340fa0100000000000000000000000035552c16704d214347f29fa77f77da6d75d7c752",
    "value": "3128508000000000",
    "logs": [
        {
            "index": 5,
            "data": "0x000000000000000000000000000000000000000000000000000a00d347371c00"
        },
        {
            "index": 4,
            "data": "0x00000000000000000000000000000000000000000000000000011c8940cd3c00"
        }
    ]
}

```
