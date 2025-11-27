package main

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/katzterd/ui/dialog"
)

var (
	dialogNodes = []dialog.Node{
		{ID: "start", Text: "Start Node", Keyboard: [][]dialog.Button{{{Text: "Go to node 2", Goto: "2"}, {Text: "Go to node 3", Goto: "3"}}, {{Text: "Go Telegram UI", URL: "https://github.com/katzterd/ui"}}}},
		{ID: "2", Text: "node 2 without keyboard"},
		{ID: "3", Text: "node 3", Keyboard: [][]dialog.Button{
			{{Text: "Go to start", Goto: "start"}, {Text: "Go to node 4", Goto: "4"}}}},
		{ID: "4", Text: "node 4", Keyboard: [][]dialog.Button{
			{{Name: "btn1", Text: "Run func & go to start", Goto: "start", Handler: func(ctx context.Context, b *bot.Bot, update *models.Update) {
				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            "Going back to start...",
				})
			}}},
			{{Name: "btn2", Text: "Back to 3", Goto: "3", Handler: func(ctx context.Context, b *bot.Bot, update *models.Update) {
				b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
					Text:            "Going back to 3...",
				})
			}}},
		}},
	}
)

func handlerDialog(ctx context.Context, b *bot.Bot, update *models.Update) {
	p := dialog.New(b, dialogNodes, dialog.WithPrefix("dialog"))

	p.Show(ctx, b, update.Message.Chat.ID, "start")
}

func handlerDialogInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	p := dialog.New(b, dialogNodes, dialog.Inline())

	p.Show(ctx, b, update.Message.Chat.ID, "start")
}
