package app

import "os"
import "log"
import "strings"
import "path/filepath"
import "html/template"

// App root template pointer.
var appTemplate *template.Template

// Register the template directory.
func (app *Initium) SetTemplateDir(dir string) {
  if err := filepath.Walk(dir, registerTemplate); err != nil {
    log.Fatal(err)
  }
}

// Callback for every single template file.
func registerTemplate(path string, file os.FileInfo, err error) error {
  if file.IsDir() || !strings.Contains(path, ".tmpl") {
    return nil
  }
  var end, start = strings.Index(path, "."), strings.Index(path, "/") + 1
  var alias = strings.Replace(path[start:end], "/", ".", -1)

  log.Println("Registering template alias:", alias)
  if appTemplate == nil {
    appTemplate, err = template.New(alias).Funcs(templateFuncs()).ParseFiles(path)
  } else {
    _, err = appTemplate.New(alias).ParseFiles(path)
  }

  if err != nil {
    log.Println("Parsing error:", err)
  }
  return nil
}

// Template function mapping for the whole application.
func templateFuncs() template.FuncMap {
  return template.FuncMap{}
}
