package annotation

import (
	"errors"
	"fmt"
	"strings"
	"text/scanner"

	"github.com/mitchellh/mapstructure"
)

func extractAnnotation(line string) (Annotation, error) {
	withoutComment := strings.TrimLeft(strings.TrimSpace(line), "/")

	anno := Annotation{
		Name:  "",
		Attrs: make(map[string]string),
	}

	var s scanner.Scanner
	s.Init(strings.NewReader(withoutComment))

	var tok rune
	curAnnoType := Initiale
	var attrName string

	for tok != scanner.EOF && curAnnoType < Done {
		tok = s.Scan()
		switch tok {
		case '@':
			curAnnoType = AnnoName
		case '(':
			curAnnoType = AttrName
		case '=':
			curAnnoType = AttrValue
		case ',':
			curAnnoType = AttrName
		case ')':
			curAnnoType = Done
		case scanner.Ident:
			switch curAnnoType {
			case AnnoName:
				anno.Name = strings.ToLower(s.TokenText())
			case AttrName:
				attrName = s.TokenText()
			}
		default:
			switch curAnnoType {
			case AttrValue:
				anno.Attrs[strings.ToLower(attrName)] = strings.Trim(s.TokenText(), "\"")
			}
		}
	}
	if curAnnoType != Done {
		return anno, fmt.Errorf("invalid comments %v for annotation:%s", curAnnoType, line)
	}
	return anno, nil
}

// ParseRouter get router from comments
func ParseRouter(docLine string) (Router, error) {
	route := Router{}
	anno, err := extractAnnotation(docLine)
	if err != nil {
		return route, err
	}
	if anno.Name != ANNO_ROUTER {
		return route, errors.New("invalid router annotation")
	}
	err = mapstructure.Decode(anno.Attrs, &route)
	if err != nil {
		return route, err
	}
	if route.Path == "" || route.Method == "" {
		return route, errors.New("invalid router annotation, path or method must be defined")
	}
	return route, nil
}

// ParseRouterGroup get routergroup from comments
func ParseRouterGroup(docLine string) (RouterGroup, error) {
	rg := RouterGroup{}
	anno, err := extractAnnotation(docLine)
	if err != nil {
		return rg, err
	}
	if anno.Name != ANNO_ROUTERGROUP {
		return rg, errors.New("invalid router group annotation")
	}
	err = mapstructure.Decode(anno.Attrs, &rg)
	if err != nil {
		return rg, err
	}
	if rg.Prefix == "" {
		return rg, errors.New("invalid router group annotation, prefix must be defined")
	}
	return rg, nil
}

// ParseMiddleware get Middleware from comments
func ParseMiddleware(docLine string) (Middleware, error) {
	mw := Middleware{}
	anno, err := extractAnnotation(docLine)
	if err != nil {
		return mw, err
	}
	if anno.Name != ANNO_MIDDLEWARE {
		return mw, errors.New("invalid middleware annotation")
	}
	err = mapstructure.Decode(anno.Attrs, &mw)
	if err != nil {
		return mw, err
	}
	if mw.Name == "" {
		return mw, errors.New("invalid middleware annotation, name must be defined")
	}
	return mw, nil
}
