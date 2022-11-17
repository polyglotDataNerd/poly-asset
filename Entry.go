package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	a "github.com/polyglotDataNerd/poly-Go-utils/aws"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	r "github.com/polyglotDataNerd/poly-asset/erc1155"
	e "github.com/polyglotDataNerd/poly-asset/erc721"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	gin.SetMode(gin.ReleaseMode)

}

func main() {
	client := gin.Default()
	awsCli := a.Settings{AWSConfig: &aws.Config{Region: aws.String("us-east-1")}}
	conn, _ := ethclient.Dial(a.SSMParams("/poly/eth/alchemy/wss/prod", 0, awsCli.SessionGenerator()))
	client.GET("/:id/:token", func(c *gin.Context) {
		asset := strings.Split(c.Request.RequestURI, "/")
		address := asset[1]
		token := asset[2]
		tokenString, _ := strconv.Atoi(token)
		t721, _ := e.NewToken(common.HexToAddress(address), conn)
		ownerCheck, _ := t721.OwnerOf(&bind.CallOpts{}, big.NewInt(int64(tokenString)))

		if ownerCheck == common.HexToAddress("0x0000000000000000000000000000000000000000") {
			metadata, err1155 := r.Get1155(address, token)
			if err1155 != nil {
				log.Error.Printf(fmt.Sprintf("asset not found: %s/%s", address, token))
				c.Error(err1155)
			} else {
				c.Data(http.StatusOK, "application/json", metadata)
			}
		} else {
			metadata, err721 := e.Get721(address, token)
			if err721 != nil {
				log.Error.Printf(fmt.Sprintf("asset not found: %s/%s", address, token))
				c.Error(err721)
			} else {
				c.Data(http.StatusOK, "application/json", metadata)
			}
		}
	})
	// health checks
	client.GET("/", func(c *gin.Context) {
		return
	})
	if serr := client.Run(); serr != nil {
		log.Error.Println("could not run server: ", serr)
	}

}
