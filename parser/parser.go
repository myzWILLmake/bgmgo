package parser

import "github.com/myzWILLmake/bgmgo/parser/websites"

var ParserCtor map[string]func() websites.WebParser

func init() {
	ParserCtor = make(map[string]func() websites.WebParser)
	ParserCtor["dmhy"] = websites.DmhyCtor
}
