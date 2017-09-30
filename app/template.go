package app

import "os"
import "log"
import "strings"
import "io/ioutil"
import "path/filepath"
import "html/template"

// Template function maps.
var templateFuncs = template.FuncMap{}

// App root template pointer.
var appTemplate *template.Template

// Register the template directory.
func (app *Initium) SetTemplateDir(dir string) *Initium {
  if appTemplate != nil {
    log.Println("Warning: Existing template store exists in memory!")
  }

  if err := filepath.Walk(dir, registerTemplate); err != nil {
    log.Fatal(err)
  }
  return app
}

// Callback for every single template file.
func registerTemplate(path string, file os.FileInfo, err error) error {
  if file.IsDir() || !strings.Contains(path, ".tmpl") {
    return nil
  }
  var end, start = strings.Index(path, "."), strings.Index(path, "/") + 1
  var alias = strings.Replace(path[start:end], "/", ".", -1)

  content, err := ioutil.ReadFile(path)
  if err != nil {
    log.Println("Error while reading template content:", err)
    return nil
  }

  log.Println("Registering template alias:", alias)
  if appTemplate == nil {
    appTemplate, err = template.New(alias).Funcs(templateFuncs).Parse(string(content))
  } else {
    _, err = appTemplate.New(alias).Parse(string(content))
  }

  if err != nil {
    log.Println("Parsing error:", err)
  }
  return nil
}

// Execute template from handler(?)
func (handler *Handler) View(template string, content interface{}) error {
  log.Println("Creating response view from", template)
  var err = appTemplate.ExecuteTemplate(handler, template, content)
  if err != nil {
    log.Println("Error while template execute:", err)
  }
  return nil
}

// Json response handler(?)
func (handler *Handler) Json(content interface{}) error {
  log.Println("Creating json response for this request.")
  return nil
}
