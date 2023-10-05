package entity

import (
	"regexp"
	"strconv"
)

type Client struct {
	ID       uint   `json:"id"`
	Tel      uint   `json:"tel"`
	Provider uint   `json:"provider"`
	Tag      string `json:"tag"`
	Tz       int    `json:"tz"`
}

type HttpToClient struct {
	Tel uint   `json:"tel"`
	Tag string `json:"tag"`
	Tz  int    `json:"tz"`
}

func NewClient(htc HttpToClient) Client {
	return Client{
		Tel:      htc.Tel,
		Provider: getProvider(htc.Tel),
		Tag:      htc.Tag,
		Tz:       htc.Tz,
	}
}

func getProvider(tel uint) uint {
	prov, err := strconv.Atoi(strconv.Itoa(int(tel))[1:4])
	if err != nil {
		return 000
	}
	return uint(prov)
}

func IsTelValid(tel uint) bool {
	result, err := regexp.Match(`7\d{10}`, []byte(strconv.Itoa(int(tel))))
	if err != nil {
		return false
	}
	return result
}
