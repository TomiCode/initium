package controllers

import "log"
import "path/filepath"
import "os"
import "strings"
import "html/template"

const BaseCategory = 0

type MenuOption struct {
  Title string
  route int
  access int
}

type MenuCategory struct {
  Title string
  options []MenuOption
}

type Categories []MenuCategory
type MenuContent map[int]Categories

var rootTemplate *template.Template

func registerOption(controller int, title string, category int, route int) {
}

func registerCategory(controller int, title string) int {
  return 0
}

func routeGenerate() string {
  return "meh"
}

func loadTemplates(path string, file os.FileInfo, err error) error {
  if file.IsDir() {
    return nil
  }
  if filepath.Ext(path) != ".tmpl" {
    return nil
  }

  var begin, end = strings.Index(path, "/"), strings.Index(path, ".")
  var templateNamespace string = strings.Replace(path[begin + 1:end], "/", ".", -1)

  var createTemplate *template.Template
  if rootTemplate != nil {
    createTemplate = rootTemplate.New(templateNamespace)
  } else {
    createTemplate = template.New(templateNamespace)
    rootTemplate = createTemplate.Funcs(template.FuncMap{"route": routeGenerate})
  }

  _, lerr := createTemplate.ParseFiles(path)
  if lerr != nil {
    log.Println(lerr)
    return lerr
  }

  log.Println("Parsed", templateNamespace)
  return nil
}

func init() {
  log.Println("InitiumViews package init.")
  if err := filepath.Walk("views/", loadTemplates); err != nil {

  }
}