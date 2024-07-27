package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const CommandSET = "SET"

type Command interface {
}

type SetCommand struct {
	key   string
	value string
}

func parseMessage(msg string) (Command, error) {

	rd := resp.NewReader(bytes.NewBufferString(string(msg)))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Read %s\n", v.Type())
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSET:
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("Error while parsing set")
					}

					cmd := SetCommand{
						key:   v.Array()[1].String(),
						value: v.Array()[2].String(),
					}
					return cmd, nil

				}

			}
		}
	}

	return nil, nil

}
