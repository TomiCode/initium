package models

type User struct {
  Username string
  Password string
}

func (*User) Create(interface{}) error {
  return nil
}

func (*User) Update(interface{}) error {
  return nil
}

func (*User) Destroy(interface{}) error {
  return nil
}