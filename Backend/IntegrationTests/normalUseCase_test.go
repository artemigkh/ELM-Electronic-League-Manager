package IntegrationTests

import (
	"testing"
)

func Test_NormalUseCase(t *testing.T) {
	doneSetup = make(chan bool)
	createRouterAndHttpclient()
	<-doneSetup
	_ = createUser(t)
}
