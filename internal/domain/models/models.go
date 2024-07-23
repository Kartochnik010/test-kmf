package models

import (
	"encoding/xml"
)

type RateItem struct {
	FullName    string `xml:"fullname"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Quant       string `xml:"quant"`
	Index       string `xml:"index"`
	Change      string `xml:"change"`
}

type RatesXML struct {
	XMLName     xml.Name   `xml:"rates"`
	Generator   string     `xml:"generator"`
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	Copyright   string     `xml:"copyright"`
	Date        string     `xml:"date"`
	Items       []RateItem `xml:"item"`
}

type Rate struct {
	Title string `json:"title"`
	Code  string `json:"code"`
	Value string `json:"value"`
	Date  string `json:"date"`
}
