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
  "runtime"
  "time"
)

type InitiumError struct {
  message string
  code int
}

func CreateError(message string, code int) *InitiumError {
  return &InitiumError{message: message, code: code}
}

func (err* InitiumError) Error() string {
  return err.message
}

type InitiumRequest struct {
  Request *http.Request
  Writer http.ResponseWriter
  vars map[string]string
  Session ApplicationSession
}

type RequestFunction func(*InitiumRequest) error

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
  sessions *SessionStorage
  templates *template.Template

  Debug bool
  Stats *runtime.MemStats
}

type TemplateParameter struct {
  SessionId string
  Debug bool
  Controller interface{}
}

func CreateInitium(debug bool) (*InitiumApp) {
  var app = &InitiumApp{Debug: debug}
  app.Initialize()

  return app
}

func (app* InitiumApp) Initialize() {
  app.sessions = CreateSessionStorage("_initium_session", 16)
  app.Stats = &runtime.MemStats{}
  go app.UpdateMemoryStats()
}

func (app* InitiumApp) UpdateMemoryStats() {
  runtime.ReadMemStats(app.Stats)
  log.Print("Update memory: Alloc: ", (app.Stats.Alloc / 1024), " KB, System: ", (app.Stats.Sys / 1024), " KB")

  time.AfterFunc(time.Duration(time.Second * 16), app.UpdateMemoryStats)
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
  var templateParam = &TemplateParameter{SessionId: request.Session.GetSessionId(), Debug: app.Debug, Controller: data}
  var err = app.templates.ExecuteTemplate(request.Writer, name, templateParam)
  if err != nil {
    log.Println("Error occurred while", name, "render:", err);
    return CreateError("Template render error", 104)
  }
  return nil
}

/* {{{ RegisterController */
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
} // }}}

/* {{{ ServeHTTP */
func (app* InitiumApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if strings.Contains(r.URL.Path, ".") {
    log.Print("File request ", r.Method, ": ", r.URL.Path)
    http.ServeFile(w, r, "public" + r.URL.Path)
    return
  }
  log.Print("Router request ", r.Method, ": ", r.URL.Path)

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
      var err error

      err = app.sessions.StartSession(requestType)
      if err != nil {
        log.Println("Session broken, reason:", err.Error())
        break;
      }

      err = route.handler(requestType)
      if err != nil {
        app.RenderTemplate(requestType, "error", err)
      }
      break
    }
  }

} // }}}
