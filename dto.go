package main

type BlockRes struct {
	BlockNum   uint64 `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	BlockTime  uint64 `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}

type BlockByIdRes struct {
	BlockNum     uint64   `json:"block_num"`
	BlockHash    string   `json:"block_hash"`
	BlockTime    uint64   `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"transactions"`
}

type TxRes struct {
	TxHash string   `json:"tx_hash"`
	From   string   `json:"from"`
	To     string   `json:"to"`
	Nonce  uint64   `json:"nonce"`
	Data   string   `json:"data"`
	Value  string   `json:"value"`
	Logs   []LogRes `json:"logs"`
}

type LogRes struct {
	Index uint64 `json:"index"`
	Data  string `json:"data"`
}
