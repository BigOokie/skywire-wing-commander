package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

// Bot provides management of the interface to the Telegram Bot
type Bot struct {
	config                 *Config
	telegram               *tgbotapi.BotAPI
	skyMgrMonitor          *SkyManagerMonitor
	commandHandlers        map[string]CommandHandler
	adminCommandHandlers   map[string]CommandHandler
	privateMessageHandlers []MessageHandler
	groupMessageHandlers   []MessageHandler
}

// BotContext provides context for Bot Messages
type BotContext struct {
	message *tgbotapi.Message
	User    *User
}

// CommandHandler provides an interface specification for command handlers
type CommandHandler func(*Bot, *BotContext, string, string) error

// MessageHandler provides an interface specification for message handlers
type MessageHandler func(*Bot, *BotContext, string) (bool, error)

// User is a structure to model the Telegram Bot user that is being interacted with
type User struct {
	ID        int    `json:"id"`
	UserName  string `db:"username" json:"username,omitempty"`
	FirstName string `db:"first_name" json:"first_name,omitempty"`
	LastName  string `db:"last_name" json:"last_name,omitempty"`
	Banned    bool   `json:"banned"`
	Admin     bool   `json:"admin"`

	exists bool
}

// NameAndTags is a helper function to append namess and tags to users within the group
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

/*
// Exists determines if the current user exists or not
func (u *User) Exists() bool {
	return u.exists
}
*/

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

func (bot *Bot) handleCommand(ctx *BotContext, command, args string) error {
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

	return fmt.Errorf("Command not found: %s", command)
}

func (bot *Bot) handlePrivateMessage(ctx *BotContext) error {
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
			errmsg := fmt.Sprintf("Sorry,'/%s' is an unknown command.\n\n%s", cmd, msgHelpShort)

			//log.Debugf("Command: '/%s %s' failed: %v", cmd, args, err)
			log.Debugf(errmsg)
			//return bot.Reply(ctx, "markdown", fmt.Sprintf("Command failed: %v", err))
			return bot.Reply(ctx, "markdown", errmsg)
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

func (bot *Bot) isReplyToMe(ctx *BotContext) bool {
	if re := ctx.message.ReplyToMessage; re != nil {
		if u := re.From; u != nil {
			if u.ID == bot.telegram.Self.ID {
				return true
			}
		}
	}
	return false
}

func (bot *Bot) handleGroupMessage(ctx *BotContext) error {
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

// Send will send a new message from the Bot using the provided BotContext
// The mode, format and text parameters are used to constuct the message and
// determine its format and delivery
func (bot *Bot) Send(ctx *BotContext, mode, format, text string) error {
	var msg tgbotapi.MessageConfig
	switch mode {
	case "whisper":
		msg = tgbotapi.NewMessage(int64(ctx.message.From.ID), text)
	case "reply":
		msg = tgbotapi.NewMessage(ctx.message.Chat.ID, text)
		msg.ReplyToMessageID = ctx.message.MessageID
	case "yell":
		msg = tgbotapi.NewMessage(bot.config.Telegram.ChatID, text)
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

/*
func (bot *Bot) Ask(ctx *BotContext, text string) error {
	msg := tgbotapi.NewMessage(ctx.message.Chat.ID, text)
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		Selective:  true,
	}
	msg.ReplyToMessageID = ctx.message.MessageID
	_, err := bot.telegram.Send(msg)
	return err
}
*/

// Reply will respond to a message recieved by the Bot in the BotContext (ctx).
// Specify the reply format and message text as parameters.
func (bot *Bot) Reply(ctx *BotContext, format, text string) error {
	return bot.Send(ctx, "reply", format, text)
}

func (bot *Bot) handleMessage(ctx *BotContext) error {
	// Check to ensure the User sending the message is registered in the Bots config
	// as the Admin user. Ignore any message or command from anyone else
	// Fixed #10
	if fmt.Sprintf("@%s", ctx.message.Chat.UserName) != bot.config.Telegram.Admin {
		log.Debugf("Bot.handleMessage: Ignoring message from non-owner user chat %d (%s)", ctx.message.Chat.ID, "@"+ctx.message.Chat.UserName)
		return nil
	}

	// If this is NOT a prive chat then DONT respond
	if !ctx.message.Chat.IsPrivate() {
		log.Debugf("Bot.handleMessage: Unknown chat %d (%s)", ctx.message.Chat.ID, ctx.message.Chat.UserName)
		return nil
	}

	log.Debug("Bot.handleMessage: handlePrivateMessage")
	return bot.handlePrivateMessage(ctx)

	/*
		if (ctx.message.Chat.IsGroup() || ctx.message.Chat.IsSuperGroup()) && ctx.message.Chat.ID == bot.config.Telegram.ChatID {
			log.Debug("Bot.handleMessage - handleGroupMessage")
			//return bot.handleGroupMessage(ctx)
			log.Printf("unknown chat %d (%s)", ctx.message.Chat.ID, ctx.message.Chat.UserName)
			return nil
		} else if ctx.message.Chat.IsPrivate() {
			log.Debug("Bot.handleMessage - handlePrivateMessage")
			return bot.handlePrivateMessage(ctx)
		} else {
			log.Debugf("unknown chat %d (%s)", ctx.message.Chat.ID, ctx.message.Chat.UserName)
			return nil
		}
	*/
}

// NewBot will create a new instance of a Bot struct based on the passed Config structure
// which supplies runtime configuration for the bot.
func NewBot(config *Config) (*Bot, error) {
	var bot = Bot{
		config:               config,
		commandHandlers:      make(map[string]CommandHandler),
		adminCommandHandlers: make(map[string]CommandHandler),
	}
	var err error

	bot.skyMgrMonitor = NewMonitor(config.SkyManager.Address)

	if bot.telegram, err = tgbotapi.NewBotAPI(config.Telegram.APIKey); err != nil {
		return nil, fmt.Errorf("Failed to initialize Telegram API: %v", err)
	}

	bot.telegram.Debug = config.Telegram.Debug

	chat, err := bot.telegram.GetChat(tgbotapi.ChatConfig{config.Telegram.ChatID, ""})
	if err != nil {
		return nil, fmt.Errorf("Failed to get chat info from Telegram: %v", err)
	}

	if !chat.IsPrivate() {
		return nil, fmt.Errorf("Only private chats are supported")
	}

	log.Printf("Bot User: %d %s", bot.telegram.Self.ID, bot.telegram.Self.UserName)
	log.Printf("Bot Chat: %s %d %s", chat.Type, chat.ID, chat.Title)

	bot.setCommandHandlers()

	return &bot, nil
}

func (bot *Bot) handleUpdate(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	ctx := BotContext{message: update.Message}

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

// Start will start the Bot running - the main duty being to monitor for and handle messages
func (bot *Bot) Start() error {
	log.Infoln("BOT: Starting.")
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	// Start the Bot Running (in the background)
	log.Infoln("Skywire Wing Commander Telegram Bot - Ready for duty.")
	defer log.Infoln("Skywire Wing Commander Telegram Bot - Signing off.")

	updates, err := bot.telegram.GetUpdatesChan(update)
	if err != nil {
		return fmt.Errorf("Failed to create Telegram updates channel: %v", err)
	}

	for update := range updates {
		if err := bot.handleUpdate(&update); err != nil {
			log.Errorf("Error: %v", err)
		}
	}
	log.Infoln("BOT: Stopped")
	return nil
}
