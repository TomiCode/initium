package main

import "regexp"
import "strings"
import "log"
import "net/http"
// import "reflect"
import "html/template"
import "os"
import "fmt"
import "math/rand"
import "io/ioutil"
import "path/filepath"
import "runtime"
import "time"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

/* Initium permissions. Those values are valid in header entries as of user permissions. */
/* With the difference, that in the header entry None, will be shown to everyone. */
const (
  Permission_None       = 0x00
  Permission_Auth_None  = 0x10
  Permission_Auth_User  = 0x11
  Permission_Auth_Mod   = 0x12
  Permission_Auth_Admin = 0x14
  Permission_Auth_Owner = 0x18
  Permission_NoAuth     = 0x20
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

type PermissionNode struct {
  Controller string
  Value uint8
}

type ControllerPermissions []*PermissionNode

func (permissions ControllerPermissions) Access(controller string) uint8 {
  for _, permission := range permissions {
    if controller == permission.Controller {
      return permission.Value
    }
  }
  return 0x00
}

type InitiumRequest struct {
  Permissions ControllerPermissions
  Controller string
  Session ApplicationSession
  Request *http.Request
  Writer http.ResponseWriter
  User *InitiumUser

  vars map[string]string
}

func (request *InitiumRequest) Redirect(url string) error {
  http.Redirect(request.Writer, request.Request, url, http.StatusFound)
  return nil
}

type RequestFunction func(*InitiumRequest) error

type ControllerRoute struct {
  uri string
  alias string
  access uint8
  method string
  call RequestFunction
}

type InitiumController interface {
  RegisterModule() *InitiumModule
  RegisterRouting() []*ControllerRoute
  RegisterOptions() []*InitiumModuleCategory
}

type ApplicationInterface interface {
  AuthenticateLogin (string, string, ApplicationSession) error
  RenderTemplate(*InitiumRequest, string, interface{}) error
  GetDatabase() (*sql.DB)

  Route(string, ...interface{}) string
  // RouteRedirect(*InitiumRequest, string,  ...interface{}) error
}

type RoutingCollection struct {
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

type ModuleCategoryCollection struct {
  Name string
  Collection []*ModuleOption
}

type ModuleCollection struct {
  Header *ModuleOption
  Options []*ModuleCategoryCollection
  Controller string
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

type InitiumHeader struct {
  Current *ModuleCollection
  Elements []*ModuleOption
}

type TemplateParameter struct {
  Header* InitiumHeader
  User* InitiumUser
  SessionId string
  Self interface{}

  Debug bool
  AuthToken string
}

type InitiumApp struct {
  routes map[string]*RoutingCollection
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
  app.routes = make(map[string]*RoutingCollection, 0)
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

func (app *InitiumApp) OpenDatabase(connection string) {
  var err error
  log.Println("Opening database connection.")
  app.database, err = sql.Open("mysql", connection)

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
  log.Print("Update memory: Alloc: ", (app.Stats.Alloc / 1024), " KB, System: ", (app.Stats.Sys / 1024), " KB")

  time.AfterFunc(time.Duration(time.Second * 16), app.UpdateMemoryStats)
}

func (app *InitiumApp) Route(name string, params ...interface{}) string {
  route, valid := app.routes[name]
  if !valid {
    log.Println("Route does not exists:", name)
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

func (app *InitiumApp) RenderTemplate(request *InitiumRequest, name string, data interface{}) error {
  log.Println("Requesting template:", name)

  var appHeader* InitiumHeader = &InitiumHeader{}

  for _, module := range app.modules {
    if module.Controller == request.Controller {
      appHeader.Current = &ModuleCollection{Header: module.Header}

      for _, category := range module.Options {
        var categoryCollection = &ModuleCategoryCollection{Name: category.Name}
        for _, option := range category.Collection {
          categoryCollection.Collection = append(categoryCollection.Collection, option)  
        }
        appHeader.Current.Options = append(appHeader.Current.Options, categoryCollection)
      }
    } else {
      appHeader.Elements = append(appHeader.Elements, module.Header)
    }
  }

  var templateParam = &TemplateParameter{
    Header: appHeader,
    User: request.User,
    Self: data,

    SessionId: request.Session.GetSessionId(),
    Debug: app.Debug,
  }

  if request.User != nil {
    templateParam.AuthToken = request.User.AuthToken
  }

  var err = app.templates.ExecuteTemplate(request.Writer, name, templateParam)
  if err != nil {
    log.Println("Error occurred while", name, "render:", err);
    return CreateError("Template render error", 104)
  }
  return nil
}

func (app *InitiumApp) RegisterController(controller InitiumController) {
  var module = controller.RegisterModule()
  var controllerAlias string

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
      if strings.HasPrefix(part, "{") {
        params = append(params, part[1:len(part) - 1])
        urlparts[idx] = "%v"
      }
    }

    var abstractUri = strings.Join(urlparts, "/")
    expr, err := regexp.Compile("^" + strings.Replace(abstractUri, "%v", "([^/]*?)", -1) + "$")
    if err != nil {
      log.Println("[Warn] Regular expression error:", v.uri)
      continue;
    }

    var routingTable = &RoutingCollection{
      expr: expr,
      params: params,
      method: v.method,
      handler: v.call,
      abstract: abstractUri,
      permission: v.access,
      controller: controllerAlias,
    }

    if v.alias != "" {
      app.routes[v.alias] = routingTable
      log.Println("Registered named route:", v.method, v.uri, "as", v.alias)
    } else {
      var unamed string = app.GenerateUUID(4)
      app.routes[unamed] = routingTable
      log.Println("Registered unnamed route:", v.method, v.uri, "as", unamed)
    }
  }
  log.Println("Routing table compiled for:", controllerAlias)

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

func (app *InitiumApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" && strings.Contains(r.URL.Path, ".") {
    log.Print("File request ", r.Method, ": ", r.URL.Path)
    http.ServeFile(w, r, "public" + r.URL.Path)
    return
  }
  log.Print("Router request ", r.Method, ": ", r.URL.Path)

  for _, route := range app.routes {
    if route.expr.MatchString(r.URL.Path) && ((route.method != "" && route.method == r.Method) || (route.method == "" && r.Method == "GET")) {
      var params = make(map[string]string, 0)
      var uriScheme = route.expr.FindStringSubmatch(r.URL.Path)
      if uriScheme[0] != r.URL.Path {
        continue;
      }

      if len(route.params) > 0 {
        for value := range uriScheme[1:] {
          log.Println("Parsing parameters", value)
          log.Println(params)
          params[route.params[value]] = uriScheme[value + 1]
        }
      }
      log.Println(params)

      for param, val := range r.URL.Query() {
        params[param] = val[0];
      }

      var requestType = &InitiumRequest{Request: r, Writer: w, Controller: route.controller, vars: params}
      var err error = nil

      err = app.sessions.StartSession(requestType)
      if err != nil {
        log.Println("Session broken, reason:", err)
        break
      }

      err = app.sessions.SessionAuthenticate(requestType)
      if err != nil {
        log.Println("Authorization failed, reason:", err)
        break
      }

      if route.permission == Permission_NoAuth && requestType.User != nil {
        log.Println("This route is not for authorized users only.")
        break;
      }

      if (route.permission & Permission_Auth_None) == Permission_Auth_None {
        if requestType.User == nil {
          log.Println("Only for authorized users.")
          break
        }

        var currentPermission uint8 = requestType.Permissions.Access(requestType.Controller)
        if (route.permission & currentPermission) != currentPermission {
          log.Println("No privileges to access this routing.")
          break;
        }
      }
      log.Println("Routing handled with controller:", requestType.Controller)
      
      err = route.handler(requestType)
      if err != nil {
        app.RenderTemplate(requestType, "error", err)
      }
      break
    }
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
      break
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
