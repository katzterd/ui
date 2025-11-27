package dialog

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type OnSelect func(ctx context.Context, bot *bot.Bot, update *models.Update)
type OnErrorHandler func(err error)

type Dialog struct {
	prefix  string
	onError OnErrorHandler
	nodes   []Node
	inline  bool

	callbackHandlerID string
}

func New(b *bot.Bot, nodes []Node, opts ...Option) *Dialog {
	p := &Dialog{
		prefix:  bot.RandomString(16),
		onError: defaultOnError,
		nodes:   nodes,
	}

	for _, opt := range opts {
		opt(p)
	}

	for _, node := range nodes {
		for _, row := range node.Keyboard {
			for _, btn := range row {
				if btn.Handler != nil {
					b.RegisterHandler(bot.HandlerTypeCallbackQueryData, p.prefix+btn.Name, bot.MatchTypePrefix, p.btnHandler(btn.Handler, btn.Goto))
				}
			}
		}
	}

	p.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, p.prefix, bot.MatchTypePrefix, p.callback)

	return p
}

func (d *Dialog) btnHandler(handler OnSelect, id string) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handler(ctx, b, update)
		if id != "" {
			d.navToNode(ctx, b, update, id)
		}
	}
}

// Prefix returns the prefix of the widget
func (d *Dialog) Prefix() string {
	return d.prefix
}

func defaultOnError(err error) {
	log.Printf("[TG-UI-DIALOG] [ERROR] %s", err)
}

func (d *Dialog) showNode(ctx context.Context, b *bot.Bot, chatID any, node Node) (*models.Message, error) {
	params := &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        node.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: node.buildKB(d.prefix),
	}

	return b.SendMessage(ctx, params)
}

func (d *Dialog) Show(ctx context.Context, b *bot.Bot, chatID any, Goto string) (*models.Message, error) {
	node, ok := d.findNode(Goto)
	if !ok {
		return nil, fmt.Errorf("failed to find node with id %s", Goto)
	}

	return d.showNode(ctx, b, chatID, node)
}

func (d *Dialog) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	if err != nil {
		d.onError(err)
	}
	if !ok {
		d.onError(fmt.Errorf("failed to answer callback query"))
	}

	Goto := strings.TrimPrefix(update.CallbackQuery.Data, d.prefix)
	d.navToNode(ctx, b, update, Goto)
}

func (d *Dialog) navToNode(ctx context.Context, b *bot.Bot, update *models.Update, id string) {
	node, ok := d.findNode(id)
	if !ok {
		d.onError(fmt.Errorf("failed to find node with id %s", id))
		return
	}

	if d.inline {
		_, errEdit := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        node.Text,
			ParseMode:   models.ParseModeMarkdown,
			ReplyMarkup: node.buildKB(d.prefix),
		})
		if errEdit != nil {
			d.onError(errEdit)
		}
		return
	}

	_, errSend := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        node.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: node.buildKB(d.prefix),
	})
	if errSend != nil {
		d.onError(errSend)
	}
}

func (d *Dialog) findNode(id string) (Node, bool) {
	for _, node := range d.nodes {
		if node.ID == id {
			return node, true
		}
	}

	return Node{}, false
}
