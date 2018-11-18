package icons

type IconManager interface {
	StoreNewIcon(tempLoc string) (string, string, error)
}
