package wingcommander

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

type Bot struct {
	config                 *Config
	telegram               *tgbotapi.BotAPI
	commandHandlers        map[string]CommandHandler
	adminCommandHandlers   map[string]CommandHandler
	privateMessageHandlers []MessageHandler
	groupMessageHandlers   []MessageHandler
}

type Context struct {
	message *tgbotapi.Message
	User    *User
}

type CommandHandler func(*Bot, *Context, string, string) error
type MessageHandler func(*Bot, *Context, string) (bool, error)

type User struct {
	ID        int    `json:"id"`
	UserName  string `db:"username" json:"username,omitempty"`
	FirstName string `db:"first_name" json:"first_name,omitempty"`
	LastName  string `db:"last_name" json:"last_name,omitempty"`
	Banned    bool   `json:"banned"`
	Admin     bool   `json:"admin"`

	exists bool
}

func (u *User) NameAndTags() string {
	var tags []string
	if u.Banned {
		tags = append(tags, "banned")
	}
	if u.Admin {
		tags = append(tags, "admin")
	}

	// If username is hidden use userid
	identifier := u.UserName
	if identifier == "" {
		identifier = strconv.Itoa(u.ID)
	}

	if len(tags) > 0 {
		return fmt.Sprintf("%s (%s)", identifier, strings.Join(tags, ", "))
	}

	return identifier
}

func (u *User) Exists() bool {
	return u.exists
}

/*
func (bot *Bot) enableUser(u *User) ([]string, error) {
	var actions []string
	if !u.Exists() {
		actions = append(actions, "created")
	}
	if u.Banned {
		u.Banned = false
		actions = append(actions, "unbanned")
	}
	//if !u.Enlisted {
	//	u.Enlisted = true
	//	actions = append(actions, "enlisted")
	//}
	//if len(actions) > 0 {
	//	if err := bot.db.PutUser(u); err != nil {
	//		return nil, fmt.Errorf("failed to change user status: %v", err)
	//	}
	//}
	return actions, nil
}
*/

/*
func (bot *Bot) handleForwardedMessageFrom(ctx *Context, id int) error {
	args := tgbotapi.ChatConfigWithUser{bot.config.ChatID, "", id}
	member, err := bot.telegram.GetChatMember(args)
	if err != nil {
		return fmt.Errorf("failed to get chat member from telegram: %v", err)
	}

	if !member.IsMember() && !member.IsCreator() && !member.IsAdministrator() {
		return bot.Reply(ctx, "that user is not a member of the chat")
	}

	user := member.User
	log.Printf("forwarded from user: %#v", user)
	dbuser := bot.db.GetUser(user.ID)
	if dbuser == nil {
		dbuser = &User{
			ID:        user.ID,
			UserName:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}

	return bot.enableUserVerbosely(ctx, dbuser)
}
*/

func (bot *Bot) handleCommand(ctx *Context, command, args string) error {
	if !ctx.User.Banned {
		handler, found := bot.commandHandlers[command]
		if found {
			return handler(bot, ctx, command, args)
		}
	}

	if ctx.User.Admin {
		handler, found := bot.adminCommandHandlers[command]
		if found {
			return handler(bot, ctx, command, args)
		}
	}

	return fmt.Errorf("command not found: %s", command)
}

func (bot *Bot) handlePrivateMessage(ctx *Context) error {
	/*
		if ctx.User.Admin {
			// let admin force add users by forwarding their messages
			if u := ctx.message.ForwardFrom; u != nil {
				if err := bot.handleForwardedMessageFrom(ctx, u.ID); err != nil {
					return fmt.Errorf("failed to add user %s: %v", u.String(), err)
				}
				return nil
			}
		}
	*/
	if ctx.message.IsCommand() {
		cmd, args := ctx.message.Command(), ctx.message.CommandArguments()
		err := bot.handleCommand(ctx, cmd, args)
		if err != nil {
			log.Printf("command '/%s %s' failed: %v", cmd, args, err)
			return bot.Reply(ctx, fmt.Sprintf("command failed: %v", err))
		}
		return nil
	}

	for i := len(bot.privateMessageHandlers) - 1; i >= 0; i-- {
		handler := bot.privateMessageHandlers[i]
		next, err := handler(bot, ctx, ctx.message.Text)
		if err != nil {
			return fmt.Errorf("private message handler failed: %v", err)
		}
		if !next {
			break
		}
	}

	return nil
}

/*
func (bot *Bot) handleUserJoin(ctx *Context, user *tgbotapi.User) error {
	if user.ID == bot.telegram.Self.ID {
		log.Printf("i have joined the group")
		return nil
	}
	dbuser := bot.db.GetUser(user.ID)
	if dbuser == nil {
		dbuser = &User{
			ID:        user.ID,
			UserName:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}
	dbuser.Enlisted = true
	if err := bot.db.PutUser(dbuser); err != nil {
		log.Printf("failed to save the user")
		return err
	}

	log.Printf("user joined: %s", dbuser.NameAndTags())
	return nil
}
*/

/*
func (bot *Bot) handleUserLeft(ctx *Context, user *tgbotapi.User) error {
	if user.ID == bot.telegram.Self.ID {
		log.Printf("i have left the group")
		return nil
	}
	dbuser := bot.db.GetUser(user.ID)
	if dbuser != nil {
		dbuser.Enlisted = false
		if err := bot.db.PutUser(dbuser); err != nil {
			log.Printf("failed to save the user")
			return err
		}

		log.Printf("user left: %s", dbuser.NameAndTags())
	}
	return nil
}
*/

func (bot *Bot) removeMyName(text string) (string, bool) {
	var removed bool
	var words []string
	for _, word := range strings.Fields(text) {
		if word == "@"+bot.telegram.Self.UserName {
			removed = true
			continue
		}
		words = append(words, word)
	}
	return strings.Join(words, " "), removed
}

func (bot *Bot) isReplyToMe(ctx *Context) bool {
	if re := ctx.message.ReplyToMessage; re != nil {
		if u := re.From; u != nil {
			if u.ID == bot.telegram.Self.ID {
				return true
			}
		}
	}
	return false
}

func (bot *Bot) handleGroupMessage(ctx *Context) error {
	var gerr error
	/*
		if u := ctx.message.NewChatMembers; u != nil {
			for _, user := range *u {
				if err := bot.handleUserJoin(ctx, &user); err != nil {
					gerr = err
				}
			}
		}
	*/
	/*
		if u := ctx.message.LeftChatMember; u != nil {
			if err := bot.handleUserLeft(ctx, u); err != nil {
				gerr = err
			}
		}
	*/
	if ctx.User != nil {
		msgWithoutName, mentioned := bot.removeMyName(ctx.message.Text)

		if mentioned || bot.isReplyToMe(ctx) {
			for i := len(bot.groupMessageHandlers) - 1; i >= 0; i-- {
				handler := bot.groupMessageHandlers[i]
				next, err := handler(bot, ctx, msgWithoutName)
				if err != nil {
					return fmt.Errorf("group message handler failed: %v", err)
				}
				if !next {
					break
				}
			}
		}
	}
	return gerr
}

func (bot *Bot) Send(ctx *Context, mode, format, text string) error {
	var msg tgbotapi.MessageConfig
	switch mode {
	case "whisper":
		msg = tgbotapi.NewMessage(int64(ctx.message.From.ID), text)
	case "reply":
		msg = tgbotapi.NewMessage(ctx.message.Chat.ID, text)
		msg.ReplyToMessageID = ctx.message.MessageID
	case "yell":
		msg = tgbotapi.NewMessage(bot.config.ChatID, text)
	default:
		return fmt.Errorf("unsupported message mode: %s", mode)
	}
	switch format {
	case "markdown":
		msg.ParseMode = "Markdown"
	case "html":
		msg.ParseMode = "HTML"
	case "text":
		msg.ParseMode = ""
	default:
		return fmt.Errorf("unsupported message format: %s", format)
	}
	_, err := bot.telegram.Send(msg)
	return err
}

/*
func (bot *Bot) ReplyAboutEvent(ctx *Context, text string, event *Event) error {
	return bot.Send(ctx, "reply", "markdown", fmt.Sprintf(
		"%s\n%s", text, formatEventAsMarkdown(event, false),
	))
}
*/

func (bot *Bot) Ask(ctx *Context, text string) error {
	msg := tgbotapi.NewMessage(ctx.message.Chat.ID, text)
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  true,
	}
	msg.ReplyToMessageID = ctx.message.MessageID
	_, err := bot.telegram.Send(msg)
	return err
}

func (bot *Bot) Reply(ctx *Context, text string) error {
	return bot.Send(ctx, "reply", "text", text)
}

func (bot *Bot) handleMessage(ctx *Context) error {
	if (ctx.message.Chat.IsGroup() || ctx.message.Chat.IsSuperGroup()) && ctx.message.Chat.ID == bot.config.ChatID {
		return bot.handleGroupMessage(ctx)
	} else if ctx.message.Chat.IsPrivate() {
		return bot.handlePrivateMessage(ctx)
	} else {
		log.Printf("unknown chat %d (%s)", ctx.message.Chat.ID, ctx.message.Chat.UserName)
		return nil
	}
}

func NewBot(config Config) (*Bot, error) {
	var bot = Bot{
		config:               &config,
		commandHandlers:      make(map[string]CommandHandler),
		adminCommandHandlers: make(map[string]CommandHandler),
	}
	var err error

	if bot.telegram, err = tgbotapi.NewBotAPI(config.Bot.Token); err != nil {
		return nil, fmt.Errorf("failed to initialize telegram api: %v", err)
	}

	bot.telegram.Debug = config.Bot.Debug

	chat, err := bot.telegram.GetChat(tgbotapi.ChatConfig{config.ChatID, ""})
	if err != nil {
		return nil, fmt.Errorf("failed to get chat info from telegram: %v", err)
	}
	if !chat.IsGroup() && !chat.IsSuperGroup() {
		return nil, fmt.Errorf("only group and supergroups are supported")
	}

	log.Printf("user: %d %s", bot.telegram.Self.ID, bot.telegram.Self.UserName)
	log.Printf("chat: %s %d %s", chat.Type, chat.ID, chat.Title)

	bot.setCommandHandlers()

	return &bot, nil
}

func (bot *Bot) handleUpdate(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	ctx := Context{message: update.Message}

	if u := ctx.message.From; u != nil {
		ctx.User = &User{
			ID:        u.ID,
			UserName:  u.UserName,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		}
	}

	return bot.handleMessage(&ctx)
}

func (bot *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.telegram.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("failed to create telegram updates channel: %v", err)
	}

	//go bot.maintain()

	for update := range updates {
		if err := bot.handleUpdate(&update); err != nil {
			log.Printf("error: %v", err)
		}
	}
	log.Printf("stopped")
	return nil
}
