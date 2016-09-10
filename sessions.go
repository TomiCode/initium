package main

import (
  "log"
  "net/http"
  "database/sql"
)

const ( 
  Session_AuthKey = "user_auth"
  Session_Cookie  = "__initium_app"
  Session_Size    = 16
)

type InitiumSession struct {
  /* Debug only field */
  sid string
  values map[string]interface{}
}

type InitiumUser struct {
  Name string
  Token int16
  Email string
  AuthToken string
}

type ApplicationSession interface {
  // Debug
  GetSessionId() string

  SetValue(string, interface{})
  GetValue(string) interface{}
  IsValid(string) bool
}

/* Debug function - this should be removed after the testing stages. */
func (session *InitiumSession) GetSessionId() string {
  return session.sid
}


func (session *InitiumSession) GetValue(key string) interface{} {
  if data, valid := session.values[key]; valid {
    return data
  }
  return nil
}

func (session *InitiumSession) IsValid(key string) bool {
  value, valid := session.values[key];
  log.Println("Accessing session:", key, "=>", value)
  return valid
}

func (session *InitiumSession) SetValue(key string, data interface{}) {
  session.values[key] = data
}

type SessionStorage struct {
  app *InitiumApp
  sessions map[string]*InitiumSession
}

func (app *InitiumApp) CreateSessionStorage() {
  log.Println("Creating session storage.")
  app.sessions = &SessionStorage{app: app, sessions: make(map[string]*InitiumSession)}
}

func (storage* SessionStorage) NewSession(sessionToken string) (*InitiumSession) {
  var newSession = &InitiumSession{sid: sessionToken, values: make(map[string]interface{})}
  storage.sessions[sessionToken] = newSession

  return newSession
}

func (storage* SessionStorage) StartSession(request* InitiumRequest) error {
  var sessionToken string = ""
  if cookie, err := request.Request.Cookie(Session_Cookie); err == nil {
    sessionToken = cookie.Value
  }

  if sessionToken != "" {
    if oldSession, valid := storage.sessions[sessionToken]; valid {
      request.Session = oldSession
      log.Println("Found valid session", sessionToken)
      return nil
    }
  }
  sessionToken = storage.app.GenerateUUID(Session_Size)
  cookie := http.Cookie{Name: Session_Cookie, Value: sessionToken, Path: "/", HttpOnly: true}

  request.Session = storage.NewSession(sessionToken)
  http.SetCookie(request.Writer, &cookie)

  log.Println("New session:", sessionToken)
  return nil
}

func (storage* SessionStorage) SessionAuthenticate(request* InitiumRequest) error {
  var db *sql.DB = storage.app.GetDatabase()

  log.Println("Started session authentication.")
  if db == nil {
    return CreateError("Database connection not valid", 901)
  }

  if !request.Session.IsValid(Session_AuthKey) {
    return nil
  }
  var auth_token string = request.Session.GetValue(Session_AuthKey).(string)
  var user_row = db.QueryRow("SELECT login, token, email FROM users WHERE auth_token=?", auth_token)
  log.Println("User auth session token:", auth_token)

  var sessionUser *InitiumUser = &InitiumUser{AuthToken: auth_token}
  var err = user_row.Scan(&sessionUser.Name, &sessionUser.Token, &sessionUser.Email)
  if err == sql.ErrNoRows {
    return nil
  } else if err != nil {
    return err
  }

  request.User = sessionUser
  return nil
}

func (storage* SessionStorage) SessionPermission(request* InitiumRequest) error {
  log.Println("Sstarting session permission managment.")
  var db *sql.DB = storage.app.GetDatabase()
  if db == nil {
    return CreateError("Database connection not exists", 901)
  }

  var auth_token = request.Session.GetValue(Session_AuthKey).(string)
  var rows, err = db.Query("SELECT permissions.controller, permissions.value FROM users JOIN permissions ON users.id = permissions.user_id WHERE users.auth_token=?", auth_token)
  if err != nil {
    log.Println("Error while permission query:", err)
    return err
  }
  defer rows.Close()

  var controllerAlias string
  var permissionValue uint8
  for rows.Next() {
    if err = rows.Scan(&controllerAlias, &permissionValue); err != nil {
      log.Println("Error while appending permission:", err)
      continue
    }

    request.Permissions = append(request.Permissions, &PermissionNode{Controller: controllerAlias, Value: permissionValue})
    log.Println("Permission node for:", controllerAlias, "with", controllerAlias)
  }
  return nil
}
