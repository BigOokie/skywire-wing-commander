// Copyright © 2018 BigOokie
//
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
package telegrambot

import (
	"testing"
)

func Test_CreateMarkup(t *testing.T) {
	kbmarkup := CreateMarkup("1", "2", "3", "4")
	if len(kbmarkup.InlineKeyboard) != 1 {
		t.Error("Expected 1 row in the Keyboard")
	}

	row := kbmarkup.InlineKeyboard[0]
	if len(row) != 4 {
		t.Error("Expected 4 buttons in the first keyboard row")
	}

	if btn := row[0]; btn.Text != "1" {
		t.Errorf("Unexpected Text for Button 1: %s", btn.Text)
	}
	if btn := row[1]; btn.Text != "2" {
		t.Errorf("Unexpected Text for Button 2: %s", btn.Text)
	}
	if btn := row[2]; btn.Text != "3" {
		t.Errorf("Unexpected Text for Button 3: %s", btn.Text)
	}
	if btn := row[3]; btn.Text != "4" {
		t.Errorf("Unexpected Text for Button 4: %s", btn.Text)
	}
}

func Test_CreateMultiLineMarkup_SingleLine(t *testing.T) {
	kbmarkup := CreateMultiLineMarkup("1", "2", "3", "4")

	if len(kbmarkup.InlineKeyboard) != 1 {
		t.Error("Expected 1 row in the Keyboard")
	}

	row := kbmarkup.InlineKeyboard[0]
	if len(row) != 4 {
		t.Error("Expected 4 buttons in the first keyboard row")
	}

	if btn := row[0]; btn.Text != "1" {
		t.Errorf("Unexpected Text for Button 1: %s", btn.Text)
	}
	if btn := row[1]; btn.Text != "2" {
		t.Errorf("Unexpected Text for Button 2: %s", btn.Text)
	}
	if btn := row[2]; btn.Text != "3" {
		t.Errorf("Unexpected Text for Button 3: %s", btn.Text)
	}
	if btn := row[3]; btn.Text != "4" {
		t.Errorf("Unexpected Text for Button 4: %s", btn.Text)
	}
}

func Test_CreateMultiLineMarkup_TwoLines(t *testing.T) {
	kbmarkup := CreateMultiLineMarkup("1", "2", "|", "3", "4", "5")

	if len(kbmarkup.InlineKeyboard) != 2 {
		t.Error("Expected 2 row in the Keyboard")
	}

	row := kbmarkup.InlineKeyboard[0]
	if len(row) != 2 {
		t.Error("Expected 2 buttons in the first keyboard row")
	}

	if btn := row[0]; btn.Text != "1" {
		t.Errorf("Unexpected Text for Row 1 Button 1: %s", btn.Text)
	}
	if btn := row[1]; btn.Text != "2" {
		t.Errorf("Unexpected Text for Row 1 Button 2: %s", btn.Text)
	}

	row = kbmarkup.InlineKeyboard[1]
	if len(row) != 3 {
		t.Error("Expected 3 buttons in the second keyboard row")
	}

	if btn := row[0]; btn.Text != "3" {
		t.Errorf("Unexpected Text for Row 2 Button 1: %s", btn.Text)
	}
	if btn := row[1]; btn.Text != "4" {
		t.Errorf("Unexpected Text for Row 2 Button 2: %s", btn.Text)
	}
	if btn := row[2]; btn.Text != "5" {
		t.Errorf("Unexpected Text for Row 2 Button 3: %s", btn.Text)
	}
}

func Test_CreateMultiLineMarkup_ThreeLines(t *testing.T) {
	kbmarkup := CreateMultiLineMarkup("1", "|", "2", "|", "3")

	if len(kbmarkup.InlineKeyboard) != 3 {
		t.Error("Expected 3 row in the Keyboard")
	}

	row := kbmarkup.InlineKeyboard[0]
	if len(row) != 1 {
		t.Error("Expected 1 button in the first keyboard row")
	}
	if btn := row[0]; btn.Text != "1" {
		t.Errorf("Unexpected Text for Row 1 Button 1: %s", btn.Text)
	}

	row = kbmarkup.InlineKeyboard[1]
	if len(row) != 1 {
		t.Error("Expected 1 buttons in the second keyboard row")
	}
	if btn := row[0]; btn.Text != "2" {
		t.Errorf("Unexpected Text for Row 2 Button 1: %s", btn.Text)
	}

	row = kbmarkup.InlineKeyboard[2]
	if len(row) != 1 {
		t.Error("Expected 1 buttons in the third keyboard row")
	}
	if btn := row[0]; btn.Text != "3" {
		t.Errorf("Unexpected Text for Row 3 Button 1: %s", btn.Text)
	}

}
