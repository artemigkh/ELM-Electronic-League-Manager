package markdown

type MdManager interface {
	StoreMarkdown(markdown, oldFile string) (string, error)
	GetMarkdown(fileName string) (string, error)
}
