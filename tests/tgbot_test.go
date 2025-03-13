package tests

import (
	"go-snippets-test/snippets"
	"testing"
)

func TestTgBot(t *testing.T) {
	snippets.RunTelegramBot()
}
