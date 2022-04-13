package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)

const (
	NftkeyURL = "https://nftkey.app/graphql"
)

type NftkeyCollection struct {
	Data struct {
		Erc721CollectionByAlias struct {
			FloorPrice   int `json:"floor"`
			typename     string `json:"__typename"`
		} `json:"erc721CollectionByAlias"`
	} `json:"data"`
}

type NftkeyResult struct {
	Floorprice float64
}

func GetNftkeyData(collection string) (NftkeyResult, error) {
	var result NftkeyResult
	var payload = fmt.Sprintf(`{
    "operationName": "GetERC721Collection",
    "variables": {
        "alias": "%s"
    },
    "query": "query GetERC721Collection($alias: String!) {\n  erc721CollectionByAlias(alias: $alias) {\n    ...ERC721CollectionInfo\n    __typename\n  }\n}\n\nfragment ERC721CollectionInfo on ERC721CollectionInfo {\n  floor\n __typename\n}\n"
}`, collection)
	var jsonStr = []byte(payload)
	req, err := http.NewRequest("POST", NftkeyURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return result, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return result, err
	}
	results, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	myString := string(results)
	var r NftkeyCollection
	err = json.Unmarshal([]byte(myString), &r)
	if err != nil {
		return result, err
	}
	result.Floorprice = float64(r.Data.Erc721CollectionByAlias.FloorPrice)
	return result, nil
}
