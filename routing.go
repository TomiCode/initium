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
  req *http.Request
  resw http.ResponseWriter

  params map[string]string
  template string
}

type RequestFunction func(*InitiumRequest)(interface{})

type ControllerRoute struct {
  path string
  method string
  fn RequestFunction
}

type InitiumController interface {
  routingRegister() []ControllerRoute
}

type RoutingCollection struct {
  expr *regexp.Regexp
  params []string
  method string
  fn RequestFunction
}

type InitiumRouter struct {
  routes []*RoutingCollection
  tmpls  *template.Template
}

func (router* InitiumRouter) TemplateWalk(path string, file os.FileInfo, err error) error {
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
    if router.tmpls == nil {
      router.tmpls = template.New(name)
      tmpl = router.tmpls
    } else {
      tmpl = router.tmpls.New(name)
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

func (router* InitiumRouter) LoadTemplates(root string) {
  err := filepath.Walk(root, router.TemplateWalk)
  if err != nil {
    log.Println("Error while template loading:", err)
  }
}

func (router* InitiumRouter) RegisterController(controller InitiumController) {
  for _, v := range controller.routingRegister() {
    urlparts := strings.Split(v.path, "/")

    var params []string
    for idx, part := range urlparts {
      if strings.HasPrefix(part, "{") {
        params = append(params, part[1:len(part) - 1])
        urlparts[idx] = "([^/]*?)"
      }
    }
    expr, err := regexp.Compile("^" + strings.Join(urlparts, "/") + "$")
    if err != nil {
      log.Println("[Warn] Can not compile regular expression for route:", v.path)
      continue;
    }
    router.routes = append(router.routes, &RoutingCollection{
      fn: v.fn,
      method: v.method,
      expr: expr,
      params: params,
    })

    log.Println("Registered", v.method, "route", v.path, "for", reflect.TypeOf(controller))
  }
}

func (router* InitiumRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  for _, route := range router.routes {
    if route.expr.MatchString(r.URL.Path) && route.method == r.Method {
      var params = make(map[string]string)

      if len(route.params) > 0 {
        // Check for [0] == r.URL.Path
        for idx, val := range route.expr.FindStringSubmatch(r.URL.Path) {
          if idx == 0 || val == "" {
            continue;
          }
          params[route.params[idx - 1]] = val
        }
      }
      for param, val := range r.URL.Query() {
        params[param] = val[0];
      }

      fnRequest := &InitiumRequest{req: r, resw: w, params: params}
      if result := route.fn(fnRequest); result != nil {
        log.Println("Rendering", fnRequest.template);
        router.tmpls.ExecuteTemplate(w, fnRequest.template, result)
      }

      return
    }
  }
}

