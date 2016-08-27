package main

import (
  "log"
  "net/http"
  "database/sql"
)

const SessionAuthKey = "user_auth"

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
  GetSessionId() string
  SetValue(string, interface{})
  GetValue(string) interface{}
  IsValid(string) bool
}

/* A bit bullshit.. */
func (request* InitiumRequest) IsAuthorized() bool {
  if request.User != nil {
    return true
  }
  return false
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
  if cookie, err := request.Request.Cookie(storage.app.SessionCookie); err == nil {
    sessionToken = cookie.Value
  }

  if sessionToken != "" {
    if oldSession, valid := storage.sessions[sessionToken]; valid {
      request.Session = oldSession
      log.Println("Found valid session", sessionToken)
      return nil
    }
  }
  sessionToken = storage.app.GenerateUUID(storage.app.SessionSize)
  cookie := http.Cookie{Name: storage.app.SessionCookie, Value: sessionToken, Path: "/", HttpOnly: true}

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

  if !request.Session.IsValid(SessionAuthKey) {
    return nil
  }
  var auth_token string = request.Session.GetValue(SessionAuthKey).(string)
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
  var db *sql.DB = storage.app.GetDatabase()
  if db == nil {
    return CreateError("Database connection not exists", 901)
  }
  var auth_token = request.Session.GetValue(SessionAuthKey).(string)
  var err = db.QueryRow("SELECT permissions.value FROM users JOIN permissions ON users.id = permissions.user_id WHERE users.auth_token=? AND permissions.controller=?", auth_token, request.Permission.Node).Scan(&request.Permission.Value)
  if err == sql.ErrNoRows {
    log.Println("No valid permission for node:", request.Permission.Node)
    request.Permission.Value = 0
    return nil
  } else if err != nil {
    return err
  }
  log.Println("Found user permission node:", request.Permission.Node, request.Permission.Value)
  return nil
}
