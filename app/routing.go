package app

import "log"
import "regexp"

// App routing element.
type AppRoute struct {
  expr *regexp.Regexp
  controller uint64
}
