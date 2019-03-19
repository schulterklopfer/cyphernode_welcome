
package main

import (
  "bytes"
  "cyphernode_status/cnAuth"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/op/go-logging"
  "github.com/spf13/viper"
  "html/template"
  "io/ioutil"
  "net/http"
  "os"
)

type Status struct {
  Foo string  `json:"foo"`
}

var auth *cnAuth.CnAuth
var statsKeyLabel string
var indexTemplate string
var statusUrl string
var log = logging.MustGetLogger("main")

func RootHandler(w http.ResponseWriter, _ *http.Request) {
  t,_ := template.ParseFiles(indexTemplate)
  t.Execute(w, nil)
}

func generateAuthHeaders() map[string]string {
  result := make( map[string]string, 0 )
  return result
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {

  req, err := http.NewRequest("GET", statusUrl, nil)
  if err != nil {
    w.WriteHeader(503 )
    return
  }

  res, err := http.DefaultClient.Do(req)

  if err != nil {
    w.WriteHeader(503 )
    return
  }

  defer res.Body.Close()

  body, _ := ioutil.ReadAll(res.Body)

  //status := new (Status)
  //jsonData, _ := json.Marshal(status)
  w.Header().Set("Content-Type", "application/json")
  bearer, err := auth.BearerFromKey(statsKeyLabel)
  if err != nil {
    w.WriteHeader(503 )
    return
  }
  w.Header().Set("Authorization", bearer )

  fmt.Fprint(w, bytes.NewBuffer(body))
  return

}

func main() {

  viper.SetConfigName("config")
  viper.AddConfigPath("data")

  err := viper.ReadInConfig()

  if err != nil {
    log.Error(err)
  }

  keysFilePath := viper.GetString("gatekeeper.key_file")
  statsKeyLabel = viper.GetString("gatekeeper.key_label")
  statusUrl = viper.GetString("gatekeeper.status_url")
  listenTo := viper.GetString("server.listen")
  indexTemplate = viper.GetString("server.index_template")

  file, err := os.Open(keysFilePath)

  if err != nil {
    log.Error(err)
  }

  auth, err = cnAuth.NewCnAuthFromFile( file )
  file.Close()

  if err != nil {
    log.Error(err)
  }

  log.Infof("Started cyphernode status page backend. URL Port [%v] ",listenTo)

  router := mux.NewRouter()
  router.HandleFunc("/", RootHandler)
  router.HandleFunc("/status", StatusHandler)
  router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

  http.Handle("/", router)

  log.Fatal(http.ListenAndServe(listenTo, nil))
}
