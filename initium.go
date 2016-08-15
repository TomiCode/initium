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

type RequestParameters struct {
}

type InitiumRequest struct {
  request *http.Request
  writer http.ResponseWriter

  vars map[string]string
}

type RequestFunction func(*InitiumRequest)(interface{})

type ControllerRoute struct {
  uri string
  method string
  template string
  call RequestFunction
}

type InitiumController interface {
  routingRegister() []ControllerRoute
}

type RoutingCollection struct {
  method string
  params []string
  template string
  expr *regexp.Regexp
  fn RequestFunction
}

type InitiumApp struct {
  routes []*RoutingCollection
  tmpls  *template.Template
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

    var tmpl *template.Template
    if app.tmpls == nil {
      app.tmpls = template.New(name)
      tmpl = app.tmpls
    } else {
      tmpl = app.tmpls.New(name)
    }

    _, err = tmpl.Parse(string(buf))
    if err != nil {
      log.Println("Error while parsing template", name, err)
      return err
    }
    log.Println("Parsed", name, "from", path)
  }
  return nil
}

func (app* InitiumApp) LoadTemplates(root string) {
  err := filepath.Walk(root, app.TemplateWalk)
  if err != nil {
    log.Println("Error while template loading:", err)
  }
}

func (app* InitiumApp) RegisterController(controller InitiumController) {
  for _, v := range controller.routingRegister() {
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
      method: v.method,
      template: v.template,
      fn: v.call,
      expr: expr,
      params: params,
    })

    log.Print("Registered route ", v.uri, " => ", reflect.TypeOf(controller))
  }
}

func (app* InitiumApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println("Handling", r.Method, "path", r.URL.Path)

  if strings.Contains(r.URL.Path, ".") {
    log.Println("File handle for path:", r.URL.Path)
    http.ServeFile(w, r, "public" + r.URL.Path)
    return
  }

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

      fnRequest := &InitiumRequest{request: r, writer: w, vars: params}
      if result := route.fn(fnRequest); result != nil && route.template != "" {
        log.Println("Execute template", route.template);
        app.tmpls.ExecuteTemplate(w, route.template, result)
      }
      return
    }
  }
}

