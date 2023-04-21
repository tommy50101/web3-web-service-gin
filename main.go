package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

type Block struct {
	gorm.Model
	ID           uint
	BlockNum     uint
	BlockHash    string
	BlockTime    uint
	ParentHash   string
	Transactions []Transaction
}
type Transaction struct {
	gorm.Model
	ID      uint
	TxHash  string
	From    string
	To      string
	Nonce   uint
	Data    string
	Value   string
	BlockID uint
	Logs    []Log
}
type BlockRes struct {
	BlockNum   uint   `json:"block_num"`
	BlockHash  string `json:"block_hash"`
	BlockTime  uint   `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}
type BlockByIdRes struct {
	BlockNum     uint     `json:"block_num"`
	BlockHash    string   `json:"block_hash"`
	BlockTime    uint     `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"transactions"`
}
type TxRes struct {
	TxHash string `json:"tx_hash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Nonce  uint   `json:"nonce"`
	Data   string `json:"data"`
	Value  string `json:"value"`
	Logs   []Log  `json:"logs"`
}
type Log struct {
	Index         uint   `json:"index"`
	Data          string `json:"data"`
	TransactionID uint   `json:"-"`
}

var (
	block       = Block{}
	blocks      = []Block{}
	transaction = Transaction{}
)

var (
	dsn string
	db  *gorm.DB
)

func (block Block) TableName() string {
	// 绑定MYSQL表名為block
	return "block"
}

func (transaction Transaction) TableName() string {
	// 绑定MYSQL表名為transaction
	return "transaction"
}

func (log Log) TableName() string {
	// 绑定MYSQL表名為log
	return "log"
}

func main() {
	// 初始化DB
	initDb()

	router := gin.Default()
	router.GET("/blocks", getBlocks)
	router.GET("/blocks/:id", getBlockByID)
	router.GET("/transaction/:txHash", getTxByTxHash)

	router.Run("localhost:8080")
}

func initDb() {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PWD")
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	Dbname := os.Getenv("DB_NAME")
	// 連線Db
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// [GET] /blocks?limit=n
func getBlocks(c *gin.Context) {
	limit := c.Query("limit")
	i, _ := strconv.Atoi(limit)
	db.Model(&Block{}).Limit(i).Find(&blocks)

	blocksRes := []BlockRes{}
	for _, value := range blocks {
		blockRes := BlockRes{
			BlockNum:   value.BlockNum,
			BlockHash:  value.BlockHash,
			BlockTime:  value.BlockTime,
			ParentHash: value.ParentHash,
		}
		blocksRes = append(blocksRes, blockRes)
	}

	c.JSON(200, gin.H{
		"blocks": blocksRes,
	})
}

// [GET] /blocks/:id
func getBlockByID(c *gin.Context) {
	id := c.Param("id")
	result := db.Model(&Block{}).Preload("Transactions").Where("id = ?", id).First(&block)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("找不到紀錄")
		return
	}

	var txHashStrList []string
	for _, a := range block.Transactions {
		txHashStrList = append(txHashStrList, a.TxHash)
	}

	blockByIdRes := BlockByIdRes{
		BlockNum:     block.BlockNum,
		BlockHash:    block.BlockHash,
		BlockTime:    block.BlockTime,
		ParentHash:   block.ParentHash,
		Transactions: txHashStrList,
	}
	c.IndentedJSON(http.StatusOK, blockByIdRes)
}

// [GET] /transaction/:txHash
func getTxByTxHash(c *gin.Context) {
	txHash := c.Param("txHash")
	result := db.Model(&Transaction{}).Preload("Logs").Where("tx_hash = ?", txHash).First(&transaction)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("找不到紀錄")
		return
	}

	var txRes = TxRes{
		TxHash: transaction.TxHash,
		From:   transaction.From,
		To:     transaction.To,
		Nonce:  transaction.Nonce,
		Data:   transaction.Data,
		Value:  transaction.Value,
		Logs:   transaction.Logs,
	}
	c.IndentedJSON(http.StatusOK, txRes)
}
