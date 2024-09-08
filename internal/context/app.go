package context

type ApplicationContext struct {
	GistID string
}

func (ctx *ApplicationContext) NewApplicationContext(gistID string) *ApplicationContext {
	return &ApplicationContext{GistID: gistID}
}
