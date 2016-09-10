package main

const (
  Alert_Information = 0x00
  Alert_Success     = 0x01
  Alert_Warning     = 0x02
  Alert_Error       = 0x04
)

type InitiumAlert struct {
  Type string
  Title string
  Message string
}

func (request *InitiumRequest) PullAlerts() []*InitiumAlert {
  var alerts = request.Session.PopValue(Session_Alerts)
  if alerts == nil {
    return []*InitiumAlert{}
  }
  return alerts.([]*InitiumAlert)
}

func (request *InitiumRequest) AddAlert(class string, title string, message string) {
  var alerts []*InitiumAlert
  var sessionAlerts = request.Session.GetValue(Session_Alerts)
  if alerts != nil {
    alerts = sessionAlerts.([]*InitiumAlert)
  }

  alerts = append(alerts, &InitiumAlert{Type: class, Title: title, Message: message})
  request.Session.SetValue(Session_Alerts, alerts)
}