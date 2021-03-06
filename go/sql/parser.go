/*
   Copyright 2016 GitHub Inc.
	 See https://github.com/github/gh-ost/blob/master/LICENSE
*/

package sql

import (
	"regexp"
	"strconv"
)

var (
	renameColumnRegexp = regexp.MustCompile(`(?i)change\s+(column\s+|)([\S]+)\s+([\S]+)\s+`)
)

type Parser struct {
	columnRenameMap map[string]string
}

func NewParser() *Parser {
	return &Parser{
		columnRenameMap: make(map[string]string),
	}
}

func (this *Parser) ParseAlterStatement(alterStatement string) (err error) {
	allStringSubmatch := renameColumnRegexp.FindAllStringSubmatch(alterStatement, -1)
	for _, submatch := range allStringSubmatch {
		if unquoted, err := strconv.Unquote(submatch[2]); err == nil {
			submatch[2] = unquoted
		}
		if unquoted, err := strconv.Unquote(submatch[3]); err == nil {
			submatch[3] = unquoted
		}

		this.columnRenameMap[submatch[2]] = submatch[3]
	}
	return nil
}

func (this *Parser) GetNonTrivialRenames() map[string]string {
	result := make(map[string]string)
	for column, renamed := range this.columnRenameMap {
		if column != renamed {
			result[column] = renamed
		}
	}
	return result
}

func (this *Parser) HasNonTrivialRenames() bool {
	return len(this.GetNonTrivialRenames()) > 0
}
