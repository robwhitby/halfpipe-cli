package model

type Options struct {
	ShowVersion bool `short:"v" long:"version" description:"Display version"`
	Args        Args `positional-args:"true"`
}

type Args struct {
	Dir string `positional-arg-name:"directory" description:"Path to process. Defaults to pwd."`
}
