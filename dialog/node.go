package dialog

import (
	"github.com/go-telegram/bot/models"
)

type Button struct {
	Name    string
	Text    string
	Goto    string
	URL     string
	Handler OnSelect
}

type Node struct {
	ID       string
	Text     string
	Keyboard [][]Button
}

func (n Node) buildKB(prefix string) models.ReplyMarkup {
	if len(n.Keyboard) == 0 {
		return nil
	}

	var kb [][]models.InlineKeyboardButton

	for _, row := range n.Keyboard {
		var kbRow []models.InlineKeyboardButton
		for _, btn := range row {
			b := models.InlineKeyboardButton{
				Text: btn.Text,
			}
			if btn.URL != "" {
				b.URL = btn.URL
			} else if btn.Handler != nil {
				b.CallbackData = prefix + btn.Name + btn.Goto
			} else {
				b.CallbackData = prefix + btn.Goto
			}
			kbRow = append(kbRow, b)
		}
		kb = append(kb, kbRow)
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: kb,
	}
}
