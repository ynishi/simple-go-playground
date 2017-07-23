package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"

	do "gopkg.in/godo.v2"
)

const tmpl = `// %s
package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello")
}
`

func tasks(p *do.Project) {

	p.Task("default", do.S{"play"}, nil)

	p.Task("play", nil, func(c *do.Context) {
		if c.FileEvent != nil {
			fmt.Printf("Do: %s\n", c.FileEvent.Path)
			//force goimports
			c.Run(fmt.Sprintf("goimports -w %s", c.FileEvent.Path))
			c.Run(fmt.Sprintf("go run %s", c.FileEvent.Path))
		}
	}).Src("*.go")

	p.Task("new", nil, func(c *do.Context) {
		// defalut is play
		prj := "play"
		if len(c.Args.NonFlags()) == 1 {
			name := c.Args.NonFlags()[0]
			if path.Ext(name) != "go" {
				prj = name
			} else {
				prj = name[0 : len(name)-3]
			}
		}
		// defalut is {{prj}].go
		filename := fmt.Sprintf("%s.go", prj)
		if _, err := os.Stat(fmt.Sprintf("%s.go", prj)); err == nil {
			files, err := ioutil.ReadDir(".")
			if err != nil {
				log.Fatal(err)
			}
			rprj := regexp.MustCompile(fmt.Sprintf("%s.[0-9]+.go$", prj))
			rnum := regexp.MustCompile(".[0-9]+.go$")
			max := 0
			for _, file := range files {
				if rprj.MatchString(file.Name()) {
					s := rnum.FindString(file.Name())
					i, err := strconv.Atoi(s[1 : len(s)-3])
					if err != nil {
						log.Fatalf("Cannot parse file num: %s", err)
					}
					if max < i {
						max = i
					}
				}
			}
			filename = fmt.Sprintf("%s.%d.go", prj, max+1)
		}
		fmt.Printf("Create: %s\n", filename)
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.Write(([]byte)(fmt.Sprintf(tmpl, filename)))
	})
}

func main() {
	do.Godo(tasks)
}
