package handler

import (
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"../interfaces"
	"strings"
	"errors"
	"io/ioutil"
)

// UranaiHandler responses constellation divination.
type UranaiHandler struct {}

// DoHandle handles a message.
func (h UranaiHandler) DoHandle(m interfaces.Message, obj interfaces.Postable) (err error) {
	response, err := h.process(m)
	if err != nil {
		return
	}

	m.Text = response
	obj.PostMessage(m)
	return
}

// Constellations map
var Constellations = map[string]string {
	"ohituji": "aries",
	"oushi": "taurus",
	"hutago": "gemini",
	"kani": "cancer",
	"shishi": "leo",
	"otome": "virgo",
	"tenbin": "libra",
	"sasori": "scorpio",
	"ite": "sagittarius",
	"yagi": "capricorn",
	"mizugame": "aquarius",
	"uo": "pisces",
}

func (h UranaiHandler) process(m interfaces.Message) (response string, err error) {
	prefix := "uranai: "
	if m.Type == "message" && strings.HasPrefix(m.Text, prefix) {
		constellation := strings.Split(m.Text, prefix)[1]
		url := fmt.Sprintf("http://fortune.yahoo.co.jp/12astro/%s", Constellations[constellation])
		uTitle, uBody, uBody2 := h.getPage(url)
		response = fmt.Sprintf("%s\n%s\n\n開運おまじない\n%s", uTitle, uBody, uBody2)
		err = nil
	} else {
		err = errors.New("Cannot parse.")
	}
	return
}

func (h UranaiHandler) getPage(url string) (uTitle string, uBody string, uBody2 string) {
	doc, _ := goquery.NewDocument(url)
	u1DlChildren := doc.Find(".mg10a .yftn12a-md48 dl").Children()
	u2DlChildren := doc.Find(".mg10b .yftn12a-md48 dl").Children()


	uTitle = h.retrieveText(u1DlChildren, "dt")
	uBody  = h.retrieveText(u1DlChildren, "dd")
	uBody2 = h.retrieveText(u2DlChildren, "dd")
	return
}

func (h UranaiHandler) retrieveText(elements *goquery.Selection, filter string) string {
	ret, _ := eucjpToUtf8(elements.Filter(filter).Text())
	return ret
}

func eucjpToUtf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.EUCJP.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}
