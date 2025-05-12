// Пакет commands предоставляет реализации различных команд бота.
package commands

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"command-bot/pkg/command"
)

// QuoteCommand предоставляет случайные вдохновляющие цитаты
type QuoteCommand struct {
	quotes []Quote
	rng    *rand.Rand
}

// Quote представляет цитату с автором
type Quote struct {
	Text   string
	Author string
}

// NewQuoteCommand создает новую команду quote
func NewQuoteCommand() *QuoteCommand {
	source := rand.NewSource(time.Now().UnixNano())

	// Предопределенный список цитат
	quotes := []Quote{
		{Text: "The only way to do great work is to love what you do.", Author: "Steve Jobs"},
		{Text: "Life is what happens when you're busy making other plans.", Author: "John Lennon"},
		{Text: "The future belongs to those who believe in the beauty of their dreams.", Author: "Eleanor Roosevelt"},
		{Text: "In the middle of difficulty lies opportunity.", Author: "Albert Einstein"},
		{Text: "Success is not final, failure is not fatal: It is the courage to continue that counts.", Author: "Winston Churchill"},
		{Text: "The best way to predict the future is to create it.", Author: "Peter Drucker"},
		{Text: "Believe you can and you're halfway there.", Author: "Theodore Roosevelt"},
		{Text: "It does not matter how slowly you go as long as you do not stop.", Author: "Confucius"},
		{Text: "The only limit to our realization of tomorrow will be our doubts of today.", Author: "Franklin D. Roosevelt"},
		{Text: "The journey of a thousand miles begins with one step.", Author: "Lao Tzu"},
		{Text: "Don't watch the clock; do what it does. Keep going.", Author: "Sam Levenson"},
		{Text: "The only person you are destined to become is the person you decide to be.", Author: "Ralph Waldo Emerson"},
		{Text: "The best revenge is massive success.", Author: "Frank Sinatra"},
		{Text: "The purpose of our lives is to be happy.", Author: "Dalai Lama"},
		{Text: "You miss 100% of the shots you don't take.", Author: "Wayne Gretzky"},
	}

	return &QuoteCommand{
		quotes: quotes,
		rng:    rand.New(source),
	}
}

// Name возвращает основное имя команды
func (c *QuoteCommand) Name() string {
	return "quote"
}

// Aliases возвращает альтернативные имена для команды
func (c *QuoteCommand) Aliases() []string {
	return []string{"inspire", "wisdom"}
}

// Description возвращает краткое описание того, что делает команда
func (c *QuoteCommand) Description() string {
	return "Provides a random inspirational quote"
}

// Usage возвращает строку, показывающую, как использовать команду
func (c *QuoteCommand) Usage() string {
	return "quote"
}

// RequiredPermissions возвращает список разрешений, необходимых для выполнения этой команды
func (c *QuoteCommand) RequiredPermissions() []string {
	return []string{} // Специальные разрешения не требуются
}

// Execute выполняет команду с заданным контекстом и возвращает ответ
func (c *QuoteCommand) Execute(ctx context.Context, cmdCtx command.CommandContext) (string, error) {
	// Выбираем случайную цитату из списка
	quote := c.quotes[c.rng.Intn(len(c.quotes))]

	return fmt.Sprintf("\"%s\"\n— %s", quote.Text, quote.Author), nil
}
