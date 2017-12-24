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

func getLatestFileNum(prj string) (max int) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	rprj := regexp.MustCompile(fmt.Sprintf("%s.[0-9]+.go$", prj))
	rnum := regexp.MustCompile(".[0-9]+.go$")
	max = 0
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
	return max
}

func getLatestFileName(prj string) (filename string) {
	max := getLatestFileNum(prj)
	filename = fmt.Sprintf("%s.%d.go", prj, max)
	_, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	return filename
}

func genNewFileName(prj string) (filename string) {
	filename = fmt.Sprintf("%s.go", prj)
	if _, err := os.Stat(fmt.Sprintf("%s.go", prj)); err == nil {
		max := getLatestFileNum(prj)
		filename = fmt.Sprintf("%s.%d.go", prj, max+1)
	}
	return filename
}

func genPrjName(defaultPrj string, c *do.Context) (prj string) {
	prj = defaultPrj
	if len(c.Args.NonFlags()) == 1 {
		name := c.Args.NonFlags()[0]
		if path.Ext(name) == ".go" {
			prj = name[0 : len(name)-3]
		} else {
			prj = name
		}
	}
	return prj
}

const DEFAULT_PRJ_NAME = "play"

func tasks(p *do.Project) {

	p.Task("default", do.S{"play"}, nil)

	p.Task("run", nil, func(c *do.Context) {
		prj := genPrjName(DEFAULT_PRJ_NAME, c)
		var filename string
		if len(c.Args.NonFlags()) == 1 {
			name := c.Args.NonFlags()[0]
			if path.Ext(name) == ".go" {
				filename = name
			} else {
				filename = getLatestFileName(prj)
			}
		} else {
			filename = getLatestFileName(prj)
		}
		fmt.Printf("Run: %s\n", filename)
		//force goimports
		c.Run(fmt.Sprintf("goimports -w %s", filename))
		c.Run(fmt.Sprintf("go run %s", filename))
	})

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
		prj := genPrjName(DEFAULT_PRJ_NAME, c)

		// defalut is {{prj}].go
		filename := genNewFileName(prj)
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
