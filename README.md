[![Go Report Card](https://goreportcard.com/badge/github.com/smurfless1/pathlib)](https://goreportcard.com/report/github.com/smurfless1/pathlib) [![license](https://img.shields.io/github/license/smurfless1/pathlib.svg)](https://github.com/smurfless1/pathlib/blob/master/LICENSE)

# pathlib

A golang path library, it is easy to use. Similar to Python pathlib.

# Installation

```
go get -u github.com/smurfless1/pathlib
```

# But why?

I have a large python codebase that I'm porting to golang, and plainly I missed having a familiar interface to 
filesystem paths.

I saw someone had started a project to provide a similar implementaiton, but my QA feelers went very, very red when 
I missed tests, interfaces, mocks, etc.

```go

package main

import "github.com/smurfless1/pathlib"

func main () {
	p := New("test.txt")

	fmt.Println(p.Absolute())
	fmt.Println(p.Cwd())
	fmt.Println(p.Parent())
	fmt.Println(p.Touch())

	fmt.Println(p.Unlink())
	fmt.Println(p.MkDir(os.ModePerm, true))
	fmt.Println(p.RmDir())
	fmt.Println(p.Open())
	fmt.Println(p.Chmod(os.ModePerm))
	fmt.Println(p.Chmod(os.ModePerm))

	fmt.Println(p.Exists())
	fmt.Println(p.IsDir())
	fmt.Println(p.IsFile())
	fmt.Println(p.IsAbs())
}

```