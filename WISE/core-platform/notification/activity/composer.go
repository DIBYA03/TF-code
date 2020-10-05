package activity

import "log"

//TextComposer simple interface to compose a text
//Ideally we could use text/template package to generate text

type MessageComposer interface {
	String(Language) (string, error)
	Title(string) string
	Body(string) string
}

type TextComposer interface {
	Compose(TemplateName, interface{}) MessageComposer
}

//NewTextComposer returns a new text composer service
func NewTextComposer() TextComposer {
	return &composer{}
}

func (c *composer) String(lang Language) (string, error) {
	return c.template.NewWithLang(c.templName, lang, c.v)
}

func (c *composer) Title(lang string) string {
	return ""
}

func (c *composer) Body(lang string) string {
	text, err := c.template.NewWithLang(c.templName, Language(lang), c.v)
	if err != nil {
		log.Printf("Error composing body for push notification err:%v", err)
	}
	return text
}

type composer struct {
	template  Template
	templName TemplateName
	v         interface{}
}

func (c *composer) Compose(template TemplateName, v interface{}) MessageComposer {
	c.templName = template
	c.v = v
	return c
}
