package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsn string
	db  *gorm.DB
)

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
	// 連線Db
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", DB_USERNAME, DB_PWD, DB_HOST, DB_PORT, DB_NAME)
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// [GET] /blocks?limit=n
func getBlocks(c *gin.Context) {
	blocks := []Block{}
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
	block := Block{}
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
	transaction := Transaction{}
	txHash := c.Param("txHash")
	result := db.Model(&Transaction{}).Preload("Logs").Where("tx_hash = ?", txHash).First(&transaction)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("找不到紀錄")
		return
	}

	var logsRes []LogRes
	for _, a := range transaction.Logs {
		logRes := LogRes{
			Index: a.Index,
			Data:  a.Data,
		}
		logsRes = append(logsRes, logRes)
	}

	var txRes = TxRes{
		TxHash: transaction.TxHash,
		From:   transaction.From,
		To:     transaction.To,
		Nonce:  transaction.Nonce,
		Data:   transaction.Data,
		Value:  transaction.Value,
		Logs:   logsRes,
	}
	c.IndentedJSON(http.StatusOK, txRes)
}
