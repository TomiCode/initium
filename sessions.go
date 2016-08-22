package main

import (
  "fmt"
  "log"
  "math/rand"
  "time"
  "net/http"
)

type InitiumSession struct {
  /* Debug only field */
  sid string

  values map[string]interface{}
}

type ApplicationSession interface {
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
  sessions map[string]*InitiumSession
  cookie string
  size int
}

func CreateSessionStorage(name string, size int) (*SessionStorage) {
  var storage = &SessionStorage{cookie: name, size: size}
  return storage.InitializeStorage()
}

func (storage *SessionStorage) InitializeStorage() (*SessionStorage) {
  storage.sessions = make(map[string]*InitiumSession, 0)

  log.Println("SessionStorage initialized.")
  return storage
}

func (storage* SessionStorage) NewSession(sessionId string) (*InitiumSession) {
  var newSession = &InitiumSession{sid: sessionId, values: make(map[string]interface{})}
  storage.sessions[sessionId] = newSession

  return newSession
}

func (storage *SessionStorage) GenerateSessionId() string {
  rand.Seed(time.Now().UnixNano())

  var result_id = make([]byte, storage.size)
  rand.Read(result_id)
  return fmt.Sprintf("%x", result_id)
}

func (storage *SessionStorage) StartSession(req* InitiumRequest) error {
  var sessionId string = ""
  if cookie, err := req.Request.Cookie(storage.cookie); err == nil {
    sessionId = cookie.Value
  }

  if sessionId != "" {
    if oldSession, valid := storage.sessions[sessionId]; valid {
      req.Session = oldSession
      log.Println("Found valid session", sessionId)
      return nil
    }
  }
  sessionId = storage.GenerateSessionId()
  cookie := http.Cookie{Name: storage.cookie, Value: sessionId, Path: "/", HttpOnly: true}

  req.Session = storage.NewSession(sessionId)
  http.SetCookie(req.Writer, &cookie)

  log.Println("New session:", sessionId)
  return nil
}
