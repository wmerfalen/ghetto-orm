package main

import (
  "regexp"
  "fmt"
  "os"
  u "wmerfalen/ghetto-orm/users/userfactory"
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
  handle,err := os.ReadDir(dir)
  if handle == nil {
    panic(err)
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
  const data_dir = "/Users/will/code/golang/hello/users/"
  const pkid_name = "pkid"

  handle, _ := os.ReadDir(data_dir)
  if handle == nil {
    os.MkdirAll(data_dir,0755)
  }

  count := count_json_files_in_dir(data_dir)
  var conf Configuration
  conf.DataDir = data_dir
  conf.PKIDFile = data_dir + pkid_name
  u.Generate(50,conf.PKIDFile,data_dir + "users-" + fmt.Sprint(count + 1) + ".json")
  return
}
