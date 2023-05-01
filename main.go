package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbName string
	dsn    string
	gormDb *gorm.DB
)

func main() {
	checkArgs()
	initDb()

	router := gin.Default()
	router.GET("/blocks", getBlocks)
	router.GET("/blocks/:id", getBlockByID)
	router.GET("/transaction/:txHash", getTxByTxHash)
	router.Run("localhost:8080")
}

// 判斷開啟哪個鏈的服務
func checkArgs() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("請輸入要查詢哪個鏈上的資訊: 1.BSC testnet   2.Ethereum testnet(goerli)   3.Ethereum mainnet ")
	fmt.Print("-> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if text[:len(text)-1] != "1" && text[:len(text)-1] != "2" && text[:len(text)-1] != "3" {
		log.Fatal("不合法的輸入")
	}

	// 依照輸入選擇不同庫
	if text[:len(text)-1] == "1" {
		dbName = "bsc_testnet"
	} else if text[:len(text)-1] == "2" {
		dbName = "eth_testnet_goerli"
	} else {
		dbName = "eth_mainnet"
	}
}

func initDb() {
	// 連線Db
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", DB_USERNAME, DB_PWD, DB_HOST, DB_PORT, dbName)
	gormDb, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// [GET] /blocks?limit=n
func getBlocks(c *gin.Context) {
	blocks := []Block{}
	limit := c.Query("limit")
	i, _ := strconv.Atoi(limit)
	gormDb.Model(&Block{}).Limit(i).Order("block_num desc").Find(&blocks)

	blocksRes := []BlockRes{}
	for _, value := range blocks {
		blockRes := BlockRes{
			BlockNum:   value.BlockNum,
			BlockHash:  value.BlockHash,
			BlockTime:  value.BlockTime,
			ParentHash: value.ParentHash,
			IsStable:   value.IsStable,
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
	result := gormDb.Model(&Block{}).Preload("Transactions").Where("block_num = ?", id).First(&block)
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
		IsStable:     block.IsStable,
		Transactions: txHashStrList,
	}
	c.JSON(http.StatusOK, blockByIdRes)
}

// [GET] /transaction/:txHash
func getTxByTxHash(c *gin.Context) {
	transaction := Transaction{}
	txHash := c.Param("txHash")
	result := gormDb.Model(&Transaction{}).Preload("Logs").Where("tx_hash = ?", txHash).First(&transaction)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("找不到紀錄")
		return
	}

	var logsRes []LogRes
	for _, a := range transaction.Logs {
		logRes := LogRes{
			Index: a.Index,
			Data:  "0x" + hex.EncodeToString(a.Data),
		}
		logsRes = append(logsRes, logRes)
	}

	var txRes = TxRes{
		TxHash: transaction.TxHash,
		From:   transaction.From,
		To:     transaction.To,
		Nonce:  transaction.Nonce,
		Data:   "0x" + hex.EncodeToString(transaction.Data),
		Value:  transaction.Value,
		Logs:   logsRes,
	}
	c.JSON(http.StatusOK, txRes)
}
