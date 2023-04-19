package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "net/http"
	// "time"
	"strconv"
)

type Block struct {
	gorm.Model
	BlockNum     uint          `gorm:"column:block_num"`
	BlockHash    string        `gorm:"column:block_hash"`
	BlockTime    uint          `gorm:"column:block_time"` // 创建时间，时间戳
	ParentHash   string        `gorm:"column:parent_hash"`
	Transactions []Transaction `gorm:"foreignKey:BlockID"`
}
type Transaction struct {
	BlockID         uint // 預設foriegn key格式
	TransactionHash string
}

type BlockRes struct {
	BlockNum   uint
	BlockHash  string
	BlockTime  uint
	ParentHash string
}

var block = Block{}
var blockRes = BlockRes{}
var getBlocksRes = []BlockRes{}

var txHashStrList []string

func (block Block) TableName() string {
	// 绑定MYSQL表名为block
	return "block"
}

func (transaction Transaction) TableName() string {
	// 绑定MYSQL表名为transaction
	return "transaction"
}

// 配置MySQL连接参数
var username = "admin"                                            //账号
var password = "Aaa6542005"                                       //密码
var host = "aws-mysql-1.cjzwlfgsosmn.us-east-1.rds.amazonaws.com" //数据库地址，可以是Ip或者域名
var port = 3306                                                   //数据库端口
var Dbname = "eth"                                                //数据库名

// 通过前面的数据库参数，拼接MYSQL DSN， 其实就是数据库连接串（数据源名称）
// 类似{username}使用花括号包着的名字都是需要替换的参数
var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)

// 连接MYSQL
var db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

func main() {
	router := gin.Default()
	router.GET("/blocks", getBlocks)
	router.GET("/blocks/:id", getBlockByID)

	router.Run("localhost:8080")
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
		fmt.Println("找不到记录")
		return
	}

	for _, a := range block.Transactions {
		txHashStrList = append(txHashStrList, a.TransactionHash)
	}

	c.JSON(200, gin.H{
		"block_num":    block.BlockNum,
		"block_hash":   block.BlockHash,
		"block_time":   block.BlockTime,
		"parent_hash":  block.ParentHash,
		"transactions": &txHashStrList,
	})
}
