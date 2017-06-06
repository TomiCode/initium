package main

import "regexp"
import "strings"
import "log"
import "net/http"

import "html/template"
import "os"
import "fmt"
import "math/rand"
import "io/ioutil"
import "path/filepath"
import "runtime"
import "time"

/* Database driver. */
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import _ "github.com/mattn/go-sqlite3"

/* Local initium packages. */
import _ "initium/controllers"
import _ "initium/models"
// import _ "initium/views"

const (
  Permission_None       = 0x00
  Permission_Auth_None  = 0x10
  Permission_Auth_User  = 0x11
  Permission_Auth_Mod   = 0x12
  Permission_Auth_Admin = 0x14
  Permission_Auth_Owner = 0x18
  Permission_NoAuth     = 0x20
)

const (
  RequestType_NORMAL = 0x00
  RequestType_FORMS  = 0x01
  RequestType_JSON   = 0x02

  /* Authorization needs to be different. */
  RequestType_API    = 0x04
)


type InitiumError struct {
  message string
  code int
}

func CreateError(message string, code int) *InitiumError {
  return &InitiumError{message: message, code: code}
}

func (err *InitiumError) Error() string {
  return err.message
}

type PermissionNode struct {
  Controller string
  Value uint8
}

type ControllerPermissions []*PermissionNode

func (permissions ControllerPermissions) Value(controller string) uint8 {
  for _, permission := range permissions {
    if controller == permission.Controller {
      return permission.Value
    }
  }
  return 0x00
}

type InitiumRequest struct {
  Permissions ControllerPermissions
  Session ApplicationSession
  Request *http.Request
  Writer http.ResponseWriter
  Route *RoutingCollection
  User *InitiumUser
}

func (request *InitiumRequest) Redirect(url string) error {
  http.Redirect(request.Writer, request.Request, url, http.StatusFound)
  return nil
}

/*
  Handle user permission as functions? 
*/
func (request *InitiumRequest) HasAccess(route *RoutingCollection) bool {
  if route == nil {
    log.Println("Routing collection is nil. Can not access permissions.")
    return false
  }

  if route.permission == Permission_NoAuth && request.User != nil {
    return false
  }

  if (route.permission & Permission_Auth_None) == Permission_Auth_None {
    if request.User == nil {
      return false
    }

    var currentPermission uint8 = request.Permissions.Value(route.controller)
    if (route.permission & currentPermission) != currentPermission {
      return false
    }
  }
  return true
}

type InitiumParameter struct {
  Name string
  Value string
}

type RequestParameters struct {
  params []*InitiumParameter
}

func (reqParams *RequestParameters) GetValue(key string) string {
  if reqParams == nil {
    return ""
  }

  for _, param := range reqParams.params {
    if key == param.Name {
      return param.Value
    }
  }
  return ""
}

type RequestFunction func(*InitiumRequest, *RequestParameters) error

type ControllerRoute struct {
  uri string
  alias string
  access uint8
  method string
  call RequestFunction
  mode uint8
}

type InitiumController interface {
  RegisterModule() *InitiumModule
  RegisterRouting() []*ControllerRoute
  RegisterOptions() []*InitiumModuleCategory
}

type ApplicationInterface interface {
  AuthenticateLogin (string, string, ApplicationSession) error
  RenderTemplate(*InitiumRequest, string, interface{}) error
  RenderData(*InitiumRequest, interface{}) error
  GetDatabase() (*sql.DB)

  Route(string, ...interface{}) string
  // RouteRedirect(*InitiumRequest, string,  ...interface{}) error
}

type RoutingCollection struct {
  mode uint8
  expr *regexp.Regexp
  method string
  params []string
  handler RequestFunction
  controller string

  // Routing alias to url path.
  abstract string
  permission uint8
}

type ModuleOption struct {
  Name string
  Route string
}

type ModuleElement struct {
  Title string
  Hash uint
}

type ModuleCategoryCollection struct {
  Name string
  Collection []*ModuleOption
}

type ModuleCollection struct {
  Header *ModuleOption
  Options []*ModuleCategoryCollection
  Controller string
}

type InternalModuleCategories struct {
  Title string
}

type InternalModule struct {
  *ModuleElement

}

type InitiumModule struct {
  Title string
  RouteName string
  ControllerAlias string
}

type InitiumModuleCategory struct {
  Title string
  Options []*ModuleOption
}

type InitiumDebugInformation struct {
  SessionId string
  AuthToken string
}

type InitiumHeader struct {
  Current *ModuleCollection
  Elements []*ModuleOption

  Alerts []*InitiumAlert
  Debug *InitiumDebugInformation
  User *InitiumUser
}

type TemplateParameter struct {
  Header *InitiumHeader
  Self interface{}
}

type InitiumApp struct {
  routes map[uint]*RoutingCollection
  modules []*ModuleCollection
  sessions *SessionStorage
  database *sql.DB
  templates *template.Template

  Debug bool
  Stats *runtime.MemStats
}

func CreateInitium(debug bool) (*InitiumApp) {
  var app = &InitiumApp{Debug: debug}
  return app.Initialize()
}

func (app *InitiumApp) Initialize() (*InitiumApp) {
  app.routes = make(map[uint]*RoutingCollection, 0)
  app.Stats = &runtime.MemStats{}
  app.CreateSessionStorage()

  go app.UpdateMemoryStats()
  return app
}

func (app *InitiumApp) GenerateUUID(size int) string {
  rand.Seed(time.Now().UnixNano())
  var result_id = make([]byte, size)
  rand.Read(result_id)

  return fmt.Sprintf("%02x", result_id)
}

func (app *InitiumApp) GenerateHash(alias string) (result uint) {
  result = 0x1EEF
  for _, r := range alias {
    result = result ^ 0x1EEF
    result = result * uint(r)
  }
  return
}

func (app *InitiumApp) OpenDatabase(connection string) {
  var err error
  log.Println("Opening database connection.")
  // app.database, err = sql.Open("mysql", connection)
  app.database, err = sql.Open("sqlite3", "./database.db")
  if err != nil {
    log.Println("Error while opening database connection:", err)
  }
}

func (app *InitiumApp) GetDatabase() (*sql.DB) {
  log.Println("Accessing database")
  if app.database == nil {
    log.Println("Error: Accessing uninitialized database connection.")
    return nil;
  }

  var err = app.database.Ping()
  if err != nil {
    log.Println("Error while database connection test:", err)
    return nil
  }
  return app.database
}

func (app *InitiumApp) CloseDatabase() {
  if app.database != nil {
    log.Println("Closing database connection.")
    var err = app.database.Close()
    if err != nil {
      log.Println("Error while closing database connection:", err)
    }
  }
}

func (app *InitiumApp) UpdateMemoryStats() {
  runtime.ReadMemStats(app.Stats)
  log.Println("Memory allocated:", (app.Stats.Alloc / 1024), "KiB, system:", (app.Stats.Sys / 1024), "KiB")

  time.AfterFunc(time.Duration(time.Second * 16), app.UpdateMemoryStats)
}

func (app *InitiumApp) Route(name string, params ...interface{}) string {
  var hash = app.GenerateHash(name)
  route, valid := app.routes[hash]
  if !valid {
    log.Println("No route for alias", name, "hash", hash)
    return ""
  }

  var appRoute string = fmt.Sprintf(route.abstract, params...)
  var err = strings.Index(appRoute, "%!")
  if err != -1 {
    log.Println("Route generation error:", appRoute)
    return appRoute[0:err]
  }
  return appRoute
}

func (app *InitiumApp) TemplateWalk(path string, file os.FileInfo, err error) error {
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
      app.templates = localTemplate.Funcs(template.FuncMap{"route": app.Route})
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

func (app *InitiumApp) LoadTemplates(root string) {
  err := filepath.Walk(root, app.TemplateWalk)
  if err != nil {
    log.Println("Error while template loading:", err)
  }
  log.Println("Loaded templates from path:", root)
}

func (app *InitiumApp) RenderTemplate(request *InitiumRequest, template string, data interface{}) error {
  log.Println("Requesting template:", template)
  var appHeader *InitiumHeader = &InitiumHeader{User: request.User}

  for _, module := range app.modules {
    if module.Controller == request.Route.controller {
      appHeader.Current = &ModuleCollection{Header: module.Header}

      for _, category := range module.Options {
        var categoryCollection = &ModuleCategoryCollection{Name: category.Name}
        for _, option := range category.Collection {
          // option.Route
          if option.Route != "" && request.HasAccess(app.routes[app.GenerateHash(option.Route)]) {
            categoryCollection.Collection = append(categoryCollection.Collection, option)
          }
        }
        if len(categoryCollection.Collection) > 0 {
          appHeader.Current.Options = append(appHeader.Current.Options, categoryCollection)
        }
      }
    } else {
      if module.Header.Route != "" && request.HasAccess(app.routes[app.GenerateHash(module.Header.Route)]) {
        appHeader.Elements = append(appHeader.Elements, module.Header)
      }
    }
  }

  appHeader.Alerts = request.PullAlerts()
  if app.Debug {
    appHeader.Debug = &InitiumDebugInformation{SessionId: request.Session.GetSessionId()}
    if request.User != nil {
      appHeader.Debug.AuthToken = request.User.AuthToken
    }
  }

  var templateParam = &TemplateParameter{Header: appHeader, Self: data}
  var err = app.templates.ExecuteTemplate(request.Writer, template, templateParam)
  if err != nil {
    log.Println("Error occurred while", template, "render:", err);
    return CreateError("Template render error", 104)
  }
  return nil
}

func (app *InitiumApp) RegisterController(controller InitiumController) {
  var module = controller.RegisterModule()
  var controllerAlias string
  var aliasHash uint

  if module != nil && module.ControllerAlias != "" {
    controllerAlias = module.ControllerAlias
  } else {
    controllerAlias = app.GenerateUUID(5)
  }
  log.Println("Started controller registration:", controllerAlias)

  for _, v := range controller.RegisterRouting() {
    var urlparts []string = strings.Split(v.uri, "/")
    var params []string

    for idx, part := range urlparts {
      if strings.HasPrefix(part, ":") {
        log.Println("Parsed route parameter:", part)
        params = append(params, part[1:len(part)])
        urlparts[idx] = "%v"
      }
    }

    var abstractUri = strings.Join(urlparts, "/")
    log.Println("Abstract url:", abstractUri)
    expr, err := regexp.Compile("^" + strings.Replace(abstractUri, "%v", "([^/]*?)", -1) + "$")
    if err != nil {
      log.Println("[Warn] Regular expression error:", v.uri)
      continue
    }

    log.Printf("Url parts: %v, routing params: %v\n", urlparts, params)

    var routingTable = &RoutingCollection{
      mode: v.mode,
      expr: expr,
      params: params,
      method: v.method,
      handler: v.call,
      abstract: abstractUri,
      permission: v.access,
      controller: controllerAlias,
    }

    if v.alias != "" {
      aliasHash = app.GenerateHash(v.alias)
    } else {
      var unamed string = app.GenerateUUID(4)
      aliasHash = app.GenerateHash(unamed)
    }
    app.routes[aliasHash] = routingTable
    log.Printf("Registered named route: %s:%s as %s hash: %x.\n", v.method, v.uri, v.alias, aliasHash)
  }
  log.Println("Routing table compiled for:", controllerAlias)
  log.Println("Route count: ", len(app.routes))

  if module == nil {
    log.Println("Controller registered without options:", controllerAlias)
    return
  }

  var moduleCollection* ModuleCollection = &ModuleCollection{
    Controller: controllerAlias,
    Header: &ModuleOption{
      Name: module.Title,
      Route: module.RouteName,
    },
  }

  for _, category := range controller.RegisterOptions() {
    var categoryCollection* ModuleCategoryCollection = &ModuleCategoryCollection{Name: category.Title}
    for _, option := range category.Options {
      log.Printf("Registered option %s [%s]\n", option.Name, category.Title)
      categoryCollection.Collection = append(categoryCollection.Collection, option)
    }
    moduleCollection.Options = append(moduleCollection.Options, categoryCollection)
  }
  app.modules = append(app.modules, moduleCollection)
  log.Println("Controller registered:", controllerAlias)
}

func (app *InitiumApp) ProcessRouting(req *InitiumRequest) (*RequestParameters, error) {
  log.Println("Application route count: ", len(app.routes), app.routes)

  for _, route := range app.routes {
    log.Println("Current route: ", route)

    if ((route.method != "" && route.method == req.Request.Method) || (route.method == "" && req.Request.Method == "GET")) && route.expr.MatchString(req.Request.URL.Path) {
      uriScheme := route.expr.FindStringSubmatch(req.Request.URL.Path)
      if uriScheme[0] != req.Request.URL.Path {
        continue
      }
      req.Route = route

      if len(route.params) > 0 {
        var reqParams = &RequestParameters{}
        for index, value := range uriScheme[1:] {
          if value == "" {
            continue
          }
          reqParams.params = append(reqParams.params, &InitiumParameter{Name: route.params[index], Value: value})
        }

        return reqParams, nil
      } else {
        return nil, nil
      }
    }
  }
  return nil, CreateError("No route for this address.", 505)
}

func (app *InitiumApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" && strings.Contains(r.URL.Path, ".") {
    log.Print("File request ", r.Method, ": ", r.URL.Path)
    http.ServeFile(w, r, "public" + r.URL.Path)
    return
  }

  log.Print("Router request ", r.Method, ": ", r.URL.Path)
  var request = &InitiumRequest{Writer: w, Request: r}

  var parameters, err = app.ProcessRouting(request)
  if err != nil {
    log.Println("Error occured:", err);
    return
  }

  err = app.sessions.StartSession(request)
  if err != nil {
    log.Println("Session error:", err)
    return
  }

  err = app.sessions.SessionAuthenticate(request)
  if err != nil {
    log.Println("Session authenticate error:", err)
    return
  }

  if !request.HasAccess(request.Route) {
    log.Println("Session has no permissions to view this route.")
    return
  }

  log.Println("Starting handler from controller:", request.Route.controller)
  err = request.Route.handler(request, parameters)
  if err != nil {
    log.Println("Handler error:", err)
    return
  }
}

func (app *InitiumApp) AuthenticateLogin(user, pass string, session ApplicationSession) error {
  log.Println("Starting Authenticate login")

  var db *sql.DB = app.GetDatabase()
  if db == nil {
    return CreateError("Can not connect to database.", 201)
  }

  var user_row = db.QueryRow("SELECT id FROM users WHERE email=? AND password=?", user, pass)
  var user_id int

  var err = user_row.Scan(&user_id)
  if err == sql.ErrNoRows {
    return CreateError("Username or password incorrect.", 301)
  } else if err != nil {
    return err
  }

  for {
    var auth_string string = app.GenerateUUID(6)

    err = db.QueryRow("SELECT id FROM users WHERE auth_token=?", auth_string).Scan(&user_id)
    if err == nil {
      continue
    } else if err != sql.ErrNoRows {
      log.Println("Error occured while database query:", err)
      return err
    }

    _, err = db.Exec("UPDATE users SET auth_token=? WHERE id=?", auth_string, user_id)
    if err != nil {
      return err
    }
    session.SetValue(Session_AuthKey, auth_string)
    break
  }
  return nil
}
