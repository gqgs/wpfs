// Code generated by argsgen.
// DO NOT EDIT!
package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
)

func (o *options) flagSet() *flag.FlagSet {
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.IntVar(&o.port, "port", o.port, "server port")
    flagSet.StringVar(&o.apiKey, "apiKey", o.apiKey, "auth key to access api")
    return flagSet
}

// Parse parses the arguments in os.Args
func (o *options) Parse() error {
    flagSet := o.flagSet()
    args := os.Args[1:]
    for len(args) > 0 {
        if err := flagSet.Parse(args); err != nil {
            return err
        }

        if remaining := flagSet.NArg(); remaining > 0 {
            posIndex := len(args) - remaining
            args = args[posIndex+1:]
            continue
        }
        break
    }

    if o.port == 0 {
        return errors.New("argument 'port' is required")
    }
    if o.apiKey == "" {
        return errors.New("argument 'apiKey' is required")
    }
    return nil
}

// MustParse parses the arguments in os.Args or exists on error
func (o *options) MustParse() {
    if err := o.Parse(); err != nil {
        o.flagSet().PrintDefaults()
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
