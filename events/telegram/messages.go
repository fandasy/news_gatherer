package telegram

import (
	"strconv"
	"strings"
)

const msgHelp = `I can save and preserve your news pages. I can also offer them for you to read. 

To save a page, just send me a link to it.
Supported pages: VK, RSS, XML

To view the entire list of saved pages, send me the /list command.

To remove a page from the list, send me the /rm pageURL command.

To display news from all pages, send me the /allnews command.

To display news from a specific page or a page from a specific platform, send me /news filter
Examples: 
 - /news https://lenta.ru/rss
 - /news https://vk.com/club228003824
 - /news VK
 - /news RSS`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand     = "Unknown command"
	msgNoSavedPages       = "You have no saved pages"
	msgNoSavedPagesRm     = "You don't have a saved page: "
	msgSaved              = "Saved!"
	msgRemove             = "Removed!"
	msgAlreadyExists      = "You have already have this page in your list"
	msgNotContainNewsFeed = "The site does not contain a news feed or does not contain RSS"
	msgTypeOrPageNotExist = "This type or page doesn't exist"
	msgNotValidateGroup   = "The group has not been validated, perhaps it does not exist or is closed"
)

func generateListMsg(pages []string, count int) string {
	pageList := strings.Join(pages, `
`)
	msgList := `There are ` + strconv.Itoa(count) + ` pages in your list:
` + pageList

	return msgList
}
