package main

import (
  "regexp"
  "strings"
  "log"
  "net/http"
  "reflect"
  "html/template"
  "os"
  "io/ioutil"
  "path/filepath"
)

type InitiumError struct {
  message string
  code int
}

func CreateError(message string, code int) *InitiumError {
  return &InitiumError{message: message, code: code}
}

type RequestParameters struct {
}

type InitiumRequest struct {
  Request *http.Request
  Writer http.ResponseWriter
  vars map[string]string
}

type RequestFunction func(*InitiumRequest) (*InitiumError)

type ControllerRoute struct {
  uri string
  name string
  auth bool
  method string
  call RequestFunction
}

type InitiumController interface {
  // InitializeController(app ApplicationInterface)
  RoutingRegister() []ControllerRoute
}

type ApplicationInterface interface {
  RenderTemplate(*InitiumRequest, string, interface{}) error
}

type RoutingCollection struct {
  auth bool
  name string
  expr *regexp.Regexp
  method string
  params []string
  handler RequestFunction
}

type InitiumApp struct {
  routes []*RoutingCollection
  templates *template.Template

  Debug bool
}

func (app* InitiumApp) TemplateWalk(path string, file os.FileInfo, err error) error {
  if !file.IsDir() {
    var str, end = strings.Index(path, "/"), strings.Index(path, ".")
    if str == -1 || end == -1 {
      log.Println("[Warning] Can not extract namespace information for file:", path)
      return nil
    }
    var name = strings.Replace(path[str+1:end], "/", ".", -1)

    buf, err := ioutil.ReadFile(path)
    if err != nil {
      log.Println("Error while reading template:", err)
      return err
    }

    var localTemplate *template.Template
    if app.templates == nil {
      localTemplate = template.New(name)
      app.templates = localTemplate
    } else {
      localTemplate = app.templates.New(name)
    }

    _, err = localTemplate.Parse(string(buf))
    if err != nil {
      log.Println("Error while loading template", name, err)
      return err
    }
    log.Println("Loaded template", name)
  }
  return nil
}

func (app* InitiumApp) LoadTemplates(root string) {
  err := filepath.Walk(root, app.TemplateWalk)
  if err != nil {
    log.Println("Error while template loading:", err)
  }
}

func (app *InitiumApp) RenderTemplate(request *InitiumRequest, name string, data interface{}) error {
  log.Println("Requesting template:", name)
  app.templates.ExecuteTemplate(request.Writer, name, data)
  return nil
}

func (app* InitiumApp) RegisterController(controller InitiumController) {
  for _, v := range controller.RoutingRegister() {
    urlparts := strings.Split(v.uri, "/")

    var params []string
    for idx, part := range urlparts {
      if strings.HasPrefix(part, "{") {
        params = append(params, part[1:len(part) - 1])
        urlparts[idx] = "([^/]*?)"
      }
    }
    expr, err := regexp.Compile("^" + strings.Join(urlparts, "/") + "$")
    if err != nil {
      log.Println("[Warn] Can not compile regular expression for route:", v.uri)
      continue;
    }

    app.routes = append(app.routes, &RoutingCollection{
      expr: expr,
      auth: v.auth,
      name: v.name,
      params: params,
      method: v.method,
      handler: v.call,
    })

    log.Print("Registered route [", v.name, "] ", v.uri, " => ", reflect.TypeOf(controller))
  }
}

func (app* InitiumApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println("Handling", r.Method, "path", r.URL.Path)

  if strings.Contains(r.URL.Path, ".") {
    log.Println("File handle for path:", r.URL.Path)
    http.ServeFile(w, r, "public" + r.URL.Path)
    return
  }

  /* {{{ Application handles routing tables */
  for _, route := range app.routes {
    if route.expr.MatchString(r.URL.Path) && ((route.method != "" && route.method == r.Method) || (route.method == "" && r.Method == "GET")) {
      var params = make(map[string]string)
      var uriScheme = route.expr.FindStringSubmatch(r.URL.Path)
      if uriScheme[0] != r.URL.Path {
        continue;
      }

      if len(route.params) > 0 {
        for value := range uriScheme[1:] {
          params[route.params[value - 1]] = uriScheme[value]
        }
      }

      for param, val := range r.URL.Query() {
        params[param] = val[0];
      }

      var requestType = &InitiumRequest{Request: r, Writer: w, vars: params}
      var err *InitiumError = nil

      err = route.handler(requestType)
      if err != nil {
        if app.Debug {
          app.RenderTemplate(requestType, "debug.error", err)
        } else {
          app.RenderTemplate(requestType, "error", err)
        }
      }
      break
    }
  } /* Application routing end }}} */

}

