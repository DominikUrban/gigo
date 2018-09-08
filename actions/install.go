package actions

import (
  "github.com/DominikUrban/gigo/helpers"
  "github.com/codegangsta/cli"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
  "path"
  "strings"
)

func Install(c *cli.Context) {
  reqs := c.String("r")

  if (reqs == "") {
    installPackages(c.Args())
    return
  }

  _, e := os.Stat(reqs)
  if os.IsNotExist(e) {
    println(reqs + " does not exist, exiting.")
    os.Exit(3)
  }

  file, err := ioutil.ReadFile(reqs)
  if err != nil {
    log.Fatal(err)
    os.Exit(3)
  }

  contents := strings.Split(string(file), "\n")
  installPackages(contents)

}

// this is why GOPATH has to bet set in main.go
func installFromGoGet(srcurl string) error {
  parts := strings.Split(srcurl, "#")
  
  if len(parts) > 2 {
    println("There can only be one hash on a line")
    os.Exit(3)
  }
  
  c := exec.Command("go", "get", parts[0])
  err := c.Run()
  if err != nil {
    log.Fatal(err)
    os.Exit(3)
  }
  
  if len(parts) == 2 {
    os.Chdir(path.Join("src", parts[0]))
    c := exec.Command("git", "checkout", parts[1])
    err := c.Run()
    if err != nil {
      log.Fatal(err)
      os.Exit(3)
    }
  }

  return nil
}

func installPackages(packages []string) {

  for _, elem := range packages {
    if (elem == "" || elem == "\n") {
      continue
    }

    println("Installing " + elem)
    if helpers.IsGoGettable(elem) {
      installFromGoGet(elem)
    } else {
      RcsGet(elem)
    }

  }

}
