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

  Permission_Auth_None  = 0x01
  Permission_Auth_User  = 0x11
  Permission_Auth_Mod   = 0x21
  Permission_Auth_Admin = 0x41
  Permission_Auth_Owner = 0x81

  Permission_NoAuth     = 0x02
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


/*
type ControllerAccess struct {
  CId string
  Value uint8
}

type ControllerPermission struct {
  permission []ControllerAccess
}

func (controllerPerm *ControllerPermission) GetAccess(cid string) uint8 {
  if controllerPerm == nil {
    return InitiumPermission_None
  }

  for _, permission := range controllerPerm.permission {
    if permission.CId == cid {
      return permission.Value
    }
  }
  return InitiumPermission_None
}

func (controllerPerms *ControllerPermission) HasAccess(cid string, perm uint8) bool {
  log.Println("HasAccess:", cid, perm, controllerPerms)
  if controllerPerms == nil {
    if (perm & InitiumPermission_Auth) != 0 {
      return false
    } else {
      return true
    }
  }

  if (perm & InitiumPermission_NoAuth) != 0 {
    return false
  }

  for _, permNode := range controllerPerms.permission {
    if permNode.CId == cid {
      log.Println("Checking permission for", cid, "with", perm, "compared", permNode.Value)
      if (perm & permNode.Value) == (perm & 0x0F) {
        return true
      } else {
        return false
      }
    }
  }
  return false
}
*/

// type ControllerPermission struct {
//   Node string
//   Value uint8
// }

type InitiumRequest struct {
  // Permission *ControllerPermission
  Controller string
  Session ApplicationSession
  Request *http.Request
  Writer http.ResponseWriter
  User *InitiumUser

  vars map[string]string
  // cid string
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

  Route(name string, params ...interface{}) string
  Redirect(request *InitiumRequest, url string) error
  // ToRoute(request *InitiumRequest, name string) error
}

type RoutingCollection struct {
  controller string 

  // cid string
  // name string
  expr *regexp.Regexp
  method string
  params []string
  handler RequestFunction

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
  // cid string

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

type InitiumApp struct {
  // permnodes map[string]string

  routes map[string]*RoutingCollection
  modules []*ModuleCollection
  sessions *SessionStorage
  database *sql.DB
  templates *template.Template
  // SessionSize int
  // SessionCookie string

  Debug bool
  Stats *runtime.MemStats
}

type TemplateParameter struct {
  Header* InitiumHeader
  User* InitiumUser
  SessionId string
  Self interface{}

  Debug bool
  AuthToken string
}

func CreateInitium(debug bool) (*InitiumApp) {
  var app = &InitiumApp{Debug: debug}
  return app.Initialize()
}

func (app *InitiumApp) Initialize() (*InitiumApp) {
  // app.permnodes = make(map[string]string)
  app.routes = make(map[string]*RoutingCollection, 0)
  app.CreateSessionStorage()
  app.Stats = &runtime.MemStats{}
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
    log.Println("Requested routing does not exists:", name)
    return ""
  }

  var appRoute string = fmt.Sprintf(route.abstract, params...)
  var err = strings.Index(appRoute, "%!")
  if err != -1 {
    log.Println("Error while generating abstract route:", appRoute)
    return appRoute[0:err]
  }
  return appRoute
}

func (app *InitiumApp) Redirect(request *InitiumRequest, url string) error {
  http.Redirect(request.Writer, request.Request, url, http.StatusFound)

  return nil
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
  // var moduleId = app.GenerateUUID(5)

  if module != nil && module.ControllerAlias != "" {
    //app.permnodes[moduleId] = module.PermissionNode
    //app.permnodes[module.PermissionNode] = moduleId
    //log.Println("Registered permission node", module.PermissionNode, "for", moduleId)

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
      log.Println("[Warn] Can not compile regular expression for route:", v.uri)
      continue;
    }

    var routingTable = &RoutingCollection{
      // cid: moduleId,
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

      // if (route.permission & InitiumPermission_NoAuth) != 0 && requestType.User != nil {
      //   log.Println("Route only for not authorized.")
      //   break
      // } else if (route.permission & InitiumPermission_Auth) == InitiumPermission_Auth && requestType.User == nil {
      //   log.Println("Route only for authorized users.")
      //   break
      // } else if route.permission != 0 {
      //   if requestType.User == nil {
      //     log.Println("Not authorized")
      //     break;
      //   }

      //   err = app.sessions.SessionPermission(requestType)
      //   if err != nil {
      //     log.Println("Error while requesting permissions.")
      //     break;
      //   }

      //   if requestType.Permission.HasAccess(route.cid, route.permission) {
      //     log.Println("Access restricted. You have no access into this route.")
      //     break
      //   }
      //   log.Println("User has access to authorized routing.")
      // }

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
