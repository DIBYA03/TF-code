package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wiseco/go-lib/log"
)

const templateDIR = "./templates/www"

//Response is all the info we need to properly render json ResponseWriter, Data, Logger, Status
type Response struct {
	Writer http.ResponseWriter
	Data   interface{}
	Logger log.Logger
	Status int
}

//ResponseTPL is all the info we need to properly render json ResponseWriter, Data, Logger, Tempate name, Status
type ResponseTPL struct {
	Writer http.ResponseWriter
	Data   interface{}
	Logger log.Logger
	Name   string
	Status int
}

var templates *template.Template

//TestTemplates is just used for testing
var TestTemplates *template.Template

//RenderJSON takes a Response struct and attempts to write json out to the repsonse writer
func RenderJSON(r Response) {
	j, err := json.Marshal(r.Data)

	if err != nil {
		r.Logger.ErrorD("Error marshalling repsonse data", log.Fields{"err": err.Error()})

		r.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Status > 0 {
		r.Writer.WriteHeader(r.Status)
	}

	r.Writer.Header().Set("Content-Type", "application/json")

	r.Writer.Write(j)
}

//RenderTemplate takes a ResponseTPL struct and attempts to write a template out to the repsonse writer
func RenderTemplate(r ResponseTPL) {
	t, err := getTemplates()

	if err != nil {
		r.Logger.ErrorD("Unable to parse templates", log.Fields{"err": err.Error()})

		r.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(r.Writer, r.Name, r.Data)

	if err != nil {
		r.Logger.ErrorD("Unable to execute template", log.Fields{"err": err.Error()})

		r.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Status > 0 {
		r.Writer.WriteHeader(r.Status)
	}

	return
}

func getTemplates() (*template.Template, error) {
	//ONLY FOR TESTING
	if TestTemplates != nil {
		return TestTemplates, nil
	}

	if templates != nil {
		return templates, nil
	}

	t := template.New("")

	if _, err := os.Stat(templateDIR); os.IsNotExist(err) {
		return t, fmt.Errorf("template dir does not exist: %s", templateDIR)
	}

	err := filepath.Walk(templateDIR, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, err = t.ParseFiles(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return t, err
	}

	templates = t

	return t, nil
}
