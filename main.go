package main

import (
  "regexp"
  "fmt"
  "os"
  userFactory "wmerfalen/ghetto-orm/users/userfactory"
  "encoding/json"
)

type Configuration struct {
  DataDir string    `json:"data_dir"`
  PKIDFile string   `json:"pkid_file"`
}


func check(e error) {
  if e != nil {
    panic(e)
  }
}

func file_exists(file string) bool {
  f, _ := os.Open(file)
  if f != nil {
    f.Close()
    return true
  }
  return false
}
func count_json_files_in_dir(dir string) int {
  var count int = 0
  handle,_ := os.ReadDir(dir)
  if handle == nil {
    return -1
  }
  regex, _ := regexp.Compile(".json$")
  for _,entry := range handle {
    if regex.MatchString(entry.Name()) {
      count += 1
    }
  }
  return count
}

func main(){
  json_schema, json_error := os.ReadFile("conf.json")
  check(json_error)

	conf := Configuration{}
	json.Unmarshal(json_schema,&conf)
	data_dir := conf.DataDir
	pkid_file := data_dir + conf.PKIDFile

  handle, _ := os.ReadDir(data_dir)
  if handle == nil {
    os.MkdirAll(data_dir,0755)
  }

  count := count_json_files_in_dir(data_dir)
  if count < 0 {
    count = 0
  }
  userFactory.Generate(50,pkid_file,data_dir + "users-" + fmt.Sprint(count + 1) + ".json")
  return
}
