package markdown

type MdManager interface {
	StoreMarkdown(leagueId int, markdown, oldFile string) (string, error)
	GetMarkdown(fileName string) (string, error)
}
