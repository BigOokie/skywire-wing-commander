// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
package telegrambot

import (
	"gopkg.in/telegram-bot-api.v4"
)

// CreateMarkup will return a tgbotapi.InlineKeyboardMarkup. Supports a single button row only
func CreateMarkup(btns ...string) tgbotapi.InlineKeyboardMarkup {
	row := tgbotapi.NewInlineKeyboardRow()
	for _, btn := range btns {
		inlineBtn := tgbotapi.NewInlineKeyboardButtonData(btn, btn)
		row = append(row, inlineBtn)
	}
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

// CreateMultiLineMarkup will build and return a multiline tgbotapi.InlineKeyboardMarkup
// The "|" character is used in the btns input to define the end of a row
// for example "1", "2", "|", "3", "4" will produce two rows with "1", "2" and "3", "4"
func CreateMultiLineMarkup(btns ...string) tgbotapi.InlineKeyboardMarkup {
	row := tgbotapi.NewInlineKeyboardRow()
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, btn := range btns {
		if btn == "|" {
			rows = append(rows, row)
			row = tgbotapi.NewInlineKeyboardRow()
		} else {
			inlineBtn := tgbotapi.NewInlineKeyboardButtonData(btn, btn)
			row = append(row, inlineBtn)
		}
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
