package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	// "time"
	"strconv"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var myAlbum album
var myAlbums []album

type Block struct {
	ID         int64  `gorm:"column:id"` // 主键
	BlockNum   int64  `gorm:"column:block_num"`
	BlockHash  string `gorm:"column:block_hash"`
	BlockTime  int64  `gorm:"column:block_time"` // 创建时间，时间戳
	ParentHash string
}

var block = Block{}
var blocks = []Block{}

// 设置表名，可以通过给struct类型定义 TableName函数，返回当前struct绑定的mysql表名是什么
func (block Block) TableName() string {
	// 绑定MYSQL表名为users
	return "block"
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

// if err != nil {
// 	panic("连接数据库失败, error=" + err.Error())
// }

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.GET("/blocks", getBlocks)

	// //定义一个用户，并初始化数据
	// block := Block{
	// 	BlockNum:   666,
	// 	BlockHash:  "0x123456789",
	// 	BlockTime:  time.Now().Unix(),
	// 	ParentHash: "0x999999999",
	// }

	// //插入一条用户数据
	// //自动生成SQL语句：INSERT INTO `block` (`block_num`,`block_hash`,`block_time`, `parent_hash`) VALUES (666,'0x123456789',date(),"0x999999999")
	// if err := db.Create(&block).Error; err != nil {
	// 	fmt.Println("插入失败", err)
	// 	return
	// }

	router.Run("localhost:8080")
}

// [GET] /blocks?limit=n
func getBlocks(c *gin.Context) {
	limit := c.Query("limit")
	i, _ := strconv.Atoi(limit)
	db.Limit(i).Find(&blocks)
	c.IndentedJSON(http.StatusOK, blocks)
}

func getBlockByID(c *gin.Context) {
	//查询并返回第一条数据
	//自动生成sql： SELECT * FROM `block`  WHERE (block_num = 666) LIMIT 1
	result := db.Where("block_num = ?", 666).First(&block)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("找不到记录")
		return
	}
	//打印查询到的数据
	// fmt.Println(block.ID, block.BlockNum, block.BlockHash)

	c.IndentedJSON(http.StatusOK, block)
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
