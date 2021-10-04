package main

import (
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/common/util"
	"github.com/mongodb/mongo-tools/mongodump"
	"log"
	"os"
)

type MongoDB struct {
}

func (m *MongoDB) dump() (string, error) {
	opts := options.New("mongodump", "1.0.0", mongodump.Usage, "", true, options.EnabledOptions{Auth: true, Connection: true, Namespace: true})
	inputOpts := &mongodump.InputOptions{}
	opts.AddOptions(inputOpts)
	outputOpts := &mongodump.OutputOptions{}
	opts.AddOptions(outputOpts)
	opts.ConnectionString = os.Getenv("MONGO_URI")

	_, err := opts.ParseArgs(os.Args)
	if err != nil {
		return "", err
	}

	dump := mongodump.MongoDump{
		ToolOptions:   opts,
		OutputOptions: outputOpts,
		InputOptions:  inputOpts,
	}

	if err := dump.Init(); err != nil {
		log.Printf("Failed: %v", err)
		os.Exit(0xA)
	}

	if err := dump.Dump(); err != nil {
		log.Printf("Failed: %v", err)
		if err == util.ErrTerminated {
			os.Exit(0xB)
		}
		os.Exit(0xC)
	}

	return "dump", nil
}
