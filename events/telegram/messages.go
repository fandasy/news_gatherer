package telegram

import (
	"strconv"
	"strings"
)

const msgHelp = `Я могу сохранять и оберегать ваши новостные страницы. Я также могу предложить их вам для чтения. 

Чтобы сохранить страницу, просто пришлите мне ссылку на нее.
Поддерживаемые страницы: VK, RSS, XML

Чтобы просмотреть весь список сохраненных страниц, отправьте мне команду /list.

Чтобы удалить страницу из списка, отправьте мне команду /rm pageURL.

Чтобы отобразить новости со всех страниц, отправьте мне команду /allnews.

Чтобы отобразить новости с определенной страницы или страницы с определенной платформы, отправьте мне команду /news filter
Примеры: 
 - /news https://lenta.ru/rss
 - /news https://vk.com/club228003824
 - /news VK
 - /news RSS`

const msgHello = "Привет! \n\n" + msgHelp

const (
	msgUnknownCommand      = "Неизвестная команда"
	msgNoSavedPages        = "У вас нет сохраненных страниц"
	msgNoSavedPagesRm      = "У вас нет сохраненной страницы: "
	msgSaved               = "Сохранено!"
	msgRemove              = "Удалено!"
	msgAlreadyExists       = "У вас уже есть эта страница в списке"
	msgNotContainNewsFeed  = "Сайт не содержит ленту новостей или не содержит RSS"
	msgTypeOrPageNotExist  = "Этот тип или страница не существуют"
	msgNotValidateGroup    = "Группа не прошла валидацию, возможно, она не существует или закрыта"
	msgImpossibleRetelling = "Невозможно пересказать"
	msgRetellingStarted    = "Пересказ начат"
)

func generateListMsg(pages []string, count int) string {
	pageList := strings.Join(pages, `
`)
	msgList := `There are ` + strconv.Itoa(count) + ` pages in your list:
` + pageList

	return msgList
}
