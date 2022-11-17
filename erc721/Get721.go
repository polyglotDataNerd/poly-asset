package erc721

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	a "github.com/polyglotDataNerd/poly-Go-utils/aws"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func Get721(assetID string, tokenID string) ([]byte, error) {
	var uriTransformed string
	awsCli := a.Settings{AWSConfig: &aws.Config{Region: aws.String("us-east-1")}}
	client, _ := ethclient.Dial(a.SSMParams("/curio/eth/quicknode/wss/prod", 0, awsCli.SessionGenerator()))
	tokenString, _ := strconv.Atoi(tokenID)
	token, _ := NewToken(common.HexToAddress(assetID), client)
	uri, terr := token.TokenURI(&bind.CallOpts{}, big.NewInt(int64(tokenString)))
	if strings.Contains(uri, "ipfs") {
		uriTransformed = strings.ReplaceAll(uri, "ipfs://", "https://ipfs.io/ipfs/")
	} else {
		uriTransformed = uri
	}
	resp, err := http.Get(uriTransformed)
	if err != nil || resp == nil {
		log.Error.Println("may not be a valid asset, try again:", fmt.Sprintf("%s/%s", assetID, tokenID))
	}
	defer resp.Body.Close()
	payload, _ := ioutil.ReadAll(resp.Body)
	return []byte(strings.ReplaceAll(string(payload), "ipfs://", "https://ipfs.io/ipfs/")), terr
}
