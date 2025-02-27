package common

import (
	"net/url"
	"strings"

	"github.com/webx-top/echo"
)

var DefaultNextURLVarName = `next`

func GetNextURL(ctx echo.Context, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	next := ctx.Form(varName)
	if next == ctx.Request().URL().Path() {
		next = ``
	}
	return next
}

func ReturnToCurrentURL(ctx echo.Context, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	next := ctx.Form(varName)
	if len(next) == 0 {
		next = ctx.Request().URI()
	}
	return next
}

func WithURLParams(urlStr string, key string, value string, args ...string) string {
	if strings.Contains(urlStr, `?`) {
		urlStr += `&`
	} else {
		urlStr += `?`
	}
	urlStr += key + `=` + url.QueryEscape(value)
	var k string
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			k = args[i]
			continue
		}
		urlStr += `&` + k + `=` + url.QueryEscape(args[i])
		k = ``
	}
	if len(k) > 0 {
		urlStr += `&` + k + `=`
	}
	return urlStr
}

func WithNextURL(ctx echo.Context, urlStr string, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	withVarName := varName
	if len(varNames) > 1 && len(varNames[1]) > 0 {
		withVarName = varNames[1]
	}

	next := GetNextURL(ctx, varName)
	if len(next) == 0 || next == urlStr {
		return urlStr
	}
	if next[0] == '/' {
		if len(urlStr) > 8 {
			var urlCopy string
			switch strings.ToLower(urlStr[0:7]) {
			case `https:/`:
				urlCopy = urlStr[8:]
			case `http://`:
				urlCopy = urlStr[7:]
			}
			if len(urlCopy) > 0 {
				p := strings.Index(urlCopy, `/`)
				if p > 0 && urlCopy[p:] == next {
					return urlStr
				}
			}
		}
	}

	return WithURLParams(urlStr, withVarName, next)
}

func GetOtherURL(ctx echo.Context, next string) string {
	if len(next) == 0 {
		return next
	}
	urlInfo, _ := url.Parse(next)
	if urlInfo == nil || urlInfo.Path == ctx.Request().URL().Path() {
		next = ``
	}
	return next
}

func FullURL(domianURL string, myURL string) string {
	if IsFullURL(myURL) {
		return myURL
	}
	if !strings.HasPrefix(myURL, `/`) && !strings.HasSuffix(domianURL, `/`) {
		myURL = `/` + myURL
	}
	myURL = domianURL + myURL
	return myURL
}

func IsFullURL(purl string) bool {
	if len(purl) == 0 {
		return false
	}
	if purl[0] == '/' {
		return false
	}
	// find "://"
	firstPos := strings.Index(purl, `/`)
	if firstPos < 0 || firstPos == len(purl)-1 {
		return false
	}
	if firstPos > 1 && purl[firstPos-1] == ':' && purl[firstPos+1] == '/' {
		return true
	}
	return false
}
