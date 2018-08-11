package Integration_Tests

type league struct {
	Id int
	Name string
	PublicView bool
	PublicJoin bool
}

type teams struct {

}

type user struct {
	Id int
	Email string
	Password string
	Leagues []*league

}