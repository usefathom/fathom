package api

import (
  "log"
)

// log fatal errors
func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
