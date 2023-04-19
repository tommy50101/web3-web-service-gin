package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "net/http"
	// "time"
	_ "github.com/joho/godotenv/autoload"
	"os"
	// "reflect"
	"strconv"
)

type Block struct {
	gorm.Model
	BlockNum     uint          `gorm:"column:block_num"`
	BlockHash    string        `gorm:"column:block_hash"`
	BlockTime    uint          `gorm:"column:block_time"`
	ParentHash   string        `gorm:"column:parent_hash"`
	Transactions []Transaction `gorm:"foreignKey:BlockID"`
}
type Transaction struct {
	gorm.Model
	TxHash  string
	From    string
	To      string
	Nonce   uint
	Data    string
	Value   string
	BlockID uint
	Logs    []Log
}
type Log struct {
	Index         uint
	Data          string
	TransactionID uint
}
type BlockRes struct {
	BlockNum   uint
	BlockHash  string
	BlockTime  uint
	ParentHash string
}

var (
	block        = Block{}
	getBlocksRes = []BlockRes{}
	transaction  = Transaction{}
	logs         = []Log{}

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
	// 連線
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// [GET] /blocks?limit=n
func getBlocks(c *gin.Context) {
	limit := c.Query("limit")
	i, _ := strconv.Atoi(limit)
	db.Model(&Block{}).Limit(i).Find(&getBlocksRes)

	c.JSON(200, gin.H{
		"blocks": &getBlocksRes,
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

	c.JSON(200, gin.H{
		"block_num":    block.BlockNum,
		"block_hash":   block.BlockHash,
		"block_time":   block.BlockTime,
		"parent_hash":  block.ParentHash,
		"transactions": &txHashStrList,
	})
}

// [GET] /transaction/:txHash
func getTxByTxHash(c *gin.Context) {
	txHash := c.Param("txHash")
	result := db.Model(&Transaction{}).Preload("Logs").Where("tx_hash = ?", txHash).First(&transaction)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("找不到紀錄")
		return
	}

	c.JSON(200, gin.H{
		"tx_hash": transaction.TxHash,
		"from":    transaction.From,
		"to":      transaction.To,
		"nonce":   transaction.Nonce,
		"data":    transaction.Data,
		"value":   transaction.Value,
		"logs":    transaction.Logs,
	})
}
