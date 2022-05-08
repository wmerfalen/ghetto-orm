package filesystem

import (
	"regexp"
  "os"
)

func FileExists(file string) bool {
  f, _ := os.Open(file)
  if f != nil {
    f.Close()
    return true
  }
  return false
}
func CountJsonFilesInDir(dir string) int {
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
