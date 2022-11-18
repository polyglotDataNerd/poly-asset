package erc1155

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	a "github.com/polyglotDataNerd/poly-Go-utils/aws"
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func Get1155(assetID string, tokenID string) ([]byte, error) {
	var uriTransformed string
	awsCli := a.Settings{AWSConfig: &aws.Config{Region: aws.String("us-west-2")}}
	client, _ := ethclient.Dial(a.SSMParams("/poly/eth/alchemy/wss/prod", 0, awsCli.SessionGenerator()))
	tokenString, _ := strconv.Atoi(tokenID)
	token, _ := NewToken(common.HexToAddress(assetID), client)

	uri, terr := token.Uri(&bind.CallOpts{}, big.NewInt(int64(tokenString)))
	if strings.Contains(uri, "ipfs") {
		uriTransformed = strings.ReplaceAll(strings.ReplaceAll(uri, "ipfs://", "https://ipfs.io/ipfs/"), "{id}", tokenID)
	} else {
		uriTransformed = strings.ReplaceAll(uri, "{id}", tokenID)
	}
	resp, err := http.Get(uriTransformed)
	if err != nil || resp == nil {
		log.Error.Println("may not be a valid asset, try again:", fmt.Sprintf("%s/%s", assetID, tokenID))
	}
	defer resp.Body.Close()
	payload, _ := io.ReadAll(resp.Body)
	return []byte(strings.ReplaceAll(string(payload), "ipfs://", "https://ipfs.io/ipfs/")), terr
}
