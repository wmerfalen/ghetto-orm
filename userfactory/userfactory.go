package userfactory

import (
  "time"
  "fmt"
  "math/rand"
  "os"
  "encoding/json"
  "strconv"
)

// create a database of users using a struct

type Person struct {
  Id int            `json:"id"`
  Name string       `json:"name"`
  Age int           `json:"age"`
  Birthday [3]int   `json:"birthday"`
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func age() int {
  return rand.Intn(66)
}

func month() int {
  var month int = 0
  for month == 0 {
    month = int(rand.Intn(12))
  }
  return month
}
func day() int {
  var d int = 0
  for d == 0 {
    d = int(rand.Intn(28))
  }
  return d
}
func year() int {
  return int(1980 + rand.Intn(15))
}

func last_name() string {
  names := []string{
    "Johnson",
    "Jones",
    "Machida",
    "Rua",
    "Benevidez",
    "Schevchenko",
    "Pierre",
    "Hughes",
    "Penn",
  }
  
  return names[rand.Intn(len(names))]

}
func first_name() string {
  names := []string{
    "John",
    "Jeff",
    "Mary",
    "Harry",
    "Larry",
    "Gordon",
    "Fred",
    "Wednesday",
  }
  
  return names[rand.Intn(len(names))]

}

var person_id int = 0
func Create() (Person) {
  var dude Person
  person_id += 1
  dude.Id = person_id
  dude.Name = first_name() + " " + last_name()
  dude.Age = age()
  dude.Birthday[0] = month()
  dude.Birthday[1] = day()
  year, _, _ := time.Now().Date()
  dude.Birthday[2] = year - dude.Age
  return dude
}

func format_birthdate(Birthday [3]int) string {
  months := []string{
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  }
  return months[Birthday[0]]+" "+ fmt.Sprint(Birthday[1]) +", " + fmt.Sprint(Birthday[2])
}

func Print_user(person Person) {
    fmt.Printf("[Id]: %d                                \n",person.Id)
    fmt.Printf("[full_name]: %s                         \n",person.Name)
    fmt.Printf("[Age]: %d                               \n",person.Age)
    fmt.Printf("[Birthday]: %s                          \n",format_birthdate(person.Birthday))
    fmt.Printf("----------------------------------------\n")
}
func Print_map(persons map[int]Person) {
  for _,person := range(persons) {
    Print_user(person)
  }
}

func read_next_pkid(file string) {
  person_id = 0
  dat, err := os.ReadFile(file)
  check(err)
  var guts string = string(dat)
  intvar, err := strconv.Atoi(guts)
  check(err)
  person_id = intvar
}

func save_next_pkid(file string){
  var Id string = fmt.Sprint(person_id)
  bytes := []byte(Id)
  err := os.WriteFile(file,bytes,0644)
  check(err)
}

func print_users(people []Person) {
  for _,value := range people {
    Print_user(value)
  }
}
func save_users_to_json_file(file string,users []Person) (bool, string)  {
  a,error := json.Marshal(users)
  check(error)
  bytes := []byte(a)
  err := os.WriteFile(file,bytes,0644)
  if err != nil {
    return false, string(err.Error())
  }
  return true, "ok"
}

// TODO: move this to file module
func file_exists(file string) bool {
  f, _ := os.Open(file)
  if f != nil {
    f.Close()
    return true
  }
  return false
}

func Generate(count int, pkid_file string, outfile_name string){
  if !file_exists(pkid_file) {
    save_next_pkid(pkid_file)
  }
  read_next_pkid(pkid_file)
  var users []Person
  for i := 0; i < count; i++ {
    users = append(users,Create())
  }
  save_next_pkid(pkid_file)
  save_users_to_json_file(outfile_name, users)
}

