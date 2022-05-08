package main

import (
	"errors"
	"regexp"
  "fmt"
	"log"
	"io/fs"
  "os"
  userFactory "wmerfalen/ghetto-orm/users/userfactory"
  wmFs "wmerfalen/ghetto-orm/users/filesystem"
  "encoding/json"
	"strconv"
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

var stderr = log.New(os.Stderr,"",0)

func bootstrap(config_file string) (bool, Configuration) {
	conf := Configuration{}
  json_schema, json_error := os.ReadFile(config_file)
	if json_error != nil {
		return false,conf
	}

	json.Unmarshal(json_schema,&conf)

  handle, _ := os.ReadDir(conf.DataDir)
  if handle == nil {
    os.MkdirAll(conf.DataDir,0755)
  }
	return true,conf
}

func get_json_file_count(config Configuration) int {
  count := wmFs.CountJsonFilesInDir(config.DataDir)

  if count < 0 {
    count = 0
	}
	return count
}

var usage string = "Usage: ./users [-c file] [--config-file=file] generate <count>\n" +
"Usage: ./users [-c file] [--config-file=file] show <id>\n" +
"Usage: ./users [-c file] [--config-file=file] filter <field> <value>\n" +
"Usage: ./users [-c file] [--config-file=file] delete <id>\n"


func main(){
	var config_file string = "conf.json"
	var mode string = "none"
	var integral_argument int = -1

	r,compile_issue := regexp.Compile("^--config-file=(.*)$")
	if compile_issue != nil {
		stderr.Println("ERROR: failed to compile regular expression: ",compile_issue)
		return
	}
	for i := 1; i < len(os.Args); i++ {
		//
		// --config-file=<FILE>
		//
		match := r.FindAllStringSubmatch(os.Args[i],-1)
		if len(match) > 0 {
			if len(match[0]) > 1 {
				stderr.Println("Using config file:",match[0][1])
				config_file = match[0][1]
				continue
			}
		}
		//
		// -c config-file-name
		//
		if os.Args[i] == "-c" {
			if i + 1 == len(os.Args) {
				stderr.Println("ERROR: expected an argument to -c")
				return
			}
			config_file = os.Args[i+1]
			i += 1
			continue
		}
		//
		// ./users generate <N>
		//
		if os.Args[i] == "generate" {
			mode = "generate"
			if len(os.Args) == i + 1 {
				stderr.Println("ERROR: please specify a number of records to generate")
				return
			}
			row_count, conversion_error := strconv.Atoi(os.Args[i+1])
			if conversion_error != nil {
				stderr.Println("ERROR: please specify a valid integer")
				return
			}
			if row_count <= 0 {
				stderr.Println("ERROR: generate expects a non-zero positive integer")
				return
			}
			integral_argument = row_count
			i += 1
			continue
		}
		//
		// ./users delete <ID>
		//
		if os.Args[i] == "delete" {
			mode = "delete"
			if len(os.Args) == i + 1 {
				stderr.Println("ERROR: delete expects an id")
				return
			}
			row_id, conversion_error := strconv.Atoi(os.Args[i+1])
			if conversion_error != nil {
				stderr.Println("ERROR: please specify a valid integer")
				return
			}
			if row_id < 0 {
				stderr.Println("ERROR: please specify a non-zero positive integer")
				return
			}
			integral_argument = row_id
			i += 1
			continue
		}
	}
	ok,conf := bootstrap(config_file)
	if !ok {
		stderr.Println("ERROR: couldn't bootstrap the app. Does",config_file,"exist?")
		return
	}
	switch mode {
		case "delete":
			delete_row(conf, integral_argument)
		case "generate":
			generate(conf, integral_argument)
		case "none":
			stderr.Println("ERROR: please choose a mode")
	}
}
func is_json_file(file_name string) bool {
	matched, _ := regexp.MatchString(".json$",file_name)
	return matched
}
func delete_row(conf Configuration, id int){
	root := conf.DataDir
	handle := os.DirFS(root)
	fs.WalkDir(handle, ".",func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if is_json_file(root + path) {
			fmt.Println(root + path)
			file_contents,file_error := os.ReadFile(root + path)
			if file_error != nil {
				fmt.Println("WARNING: issue reading",root+path,":",file_error," skipping...")
				return nil
			}
			var users []userFactory.Person
			json_error := json.Unmarshal(file_contents,&users)
			if json_error != nil {
				fmt.Println("WARNING: issue parsing json file:",root+path,":",json_error," skipping...")
				return nil
			}
			var keepers []userFactory.Person
			for _,user := range users {
				if user.Id == id {
					fmt.Println("Found id:",user)
					for _,filtered_user := range users {
						if filtered_user.Id == id {
							continue
						}
						keepers = append(keepers,filtered_user)
					}
					userFactory.SaveUsersToJsonFile(root+path,keepers)
					// TODO: there HAS to be a better way to signal a stop... right..?
					return errors.New("stop")
				}
			}
		}
		return nil
	})
	os.ReadDir(conf.DataDir)
}
func generate(conf Configuration, count int){
	json_file_count := get_json_file_count(conf)
  userFactory.Generate(count,conf.PKIDFile,conf.DataDir + "users-" + fmt.Sprint(json_file_count + 1) + ".json")
  return
}
