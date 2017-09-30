package app

import "os"
import "log"
import "bytes"
import "strings"
import "io/ioutil"
import "path/filepath"
import "html/template"

// Main template frame.
type TemplateFrame struct {
  template string
  content interface{}
}

// Template function maps.
var templateFuncs = template.FuncMap{
  "yeld": tmplFunc_yeld,
}

// App root template pointer.
var appTemplate *template.Template

// Yeld function.
func tmplFunc_yeld(frame TemplateFrame) template.HTML {
  log.Println("AppFrame yeld", frame)

  var content bytes.Buffer
  var err = appTemplate.ExecuteTemplate(&content, frame.template, frame.content)
  if err != nil {
    log.Println("Yeld execution error:", err)
    return ""
  }
  return template.HTML(content.String())
}

// Register the template directory.
func (app *Initium) SetTemplateDir(dir string) *Initium {
  if appTemplate != nil {
    log.Println("Warning: Existing template store exists in memory!")
  }
  if err := filepath.Walk(dir, registerTemplate); err != nil {
    log.Fatal(err)
  }

  log.Println("Templates load finished.")
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

// Execute application frame.
func (handler *Handler) View(template string, content interface{}) error {
  log.Println("Creating response view from", template)

  var frame = TemplateFrame{template: template, content: content}
  var err = appTemplate.ExecuteTemplate(handler, "app", frame)
  if err != nil {
    log.Println("Error while appframe execute:", err)
  }
  return nil
}
