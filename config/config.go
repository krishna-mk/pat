package config

import (
	"flag"
	"io/ioutil"
	"os"
	"fmt"
	"launchpad.net/goyaml"
)

type Config interface {
	StringVar(target *string, name string, defaultValue string, description string)
	IntVar(target *int, name string, defaultValue int, description string)
	BoolVar(target *bool, name string, defaultValue bool, description string)
	EnvVar(target *string, name string, defaultValue string, description string)
	Parse(args []string) error
}

type f struct {
	flagSet *flag.FlagSet
	envVars []env
	targets map[string]interface{}
}

type env struct {
	target       *string
	name         string
	defaultValue string
	description  string
}

func NewConfig() *f {
	return &f{flag.NewFlagSet(os.Args[0], flag.ExitOnError), make([]env, 0), make(map[string]interface{})}
}

var ConfigAndFlags = NewConfig()

func (f *f) StringVar(target *string, name string, defaultValue string, description string) {
	f.allowDoubleSetting(target, name, func() {
		f.flagSet.StringVar(target, name, defaultValue, description)
	})
}

func (f *f) IntVar(target *int, name string, defaultValue int, description string) {
	f.allowDoubleSetting(target, name, func() {
		f.flagSet.IntVar(target, name, defaultValue, description)
	})
}

func (f *f) BoolVar(target *bool, name string, defaultValue bool, description string) {
	f.allowDoubleSetting(target, name, func() {
		f.flagSet.BoolVar(target, name, defaultValue, description)
	})
}

func (f *f) EnvVar(target *string, name string, defaultValue string, description string) {
	f.envVars = append(f.envVars, env{target, name, defaultValue, description})
}

func (f *f) allowDoubleSetting(target interface{}, name string, fn func()) {
	if existing := f.flagSet.Lookup(name); existing == nil {
		f.targets[name] = target
		fn()
	} else if f.targets[name] != target {
		panic("Tried to redefine flag: " + name)
	}
}

func (f *f) Parse(args []string) error {
	
	fmt.Println("config/config.go ..Parse - .checking the config file...")
	config := f.flagSet.String("config", "", "YML file containing configuration parameters")

	fmt.Println("1111....checking the config file..%v.",config)
	if err := f.ParseEnv(); err != nil {


	fmt.Println("config/config.go*****1111**********....checking the config file...")

		return err
	}

	fmt.Println("config/config.go 2222....checking the config file..%v.",args)
	

	f.flagSet.Parse(args)

	fmt.Println("config/config.go 2222+++++....checking the config file..%v.",args)
	if len(*config) > 0 {


	// fmt.Println("config/config.go********22222******....checking the config file...")
		if err := f.ParseConfig(*config); err != nil {
		
	fmt.Println("config/config.go***********3333****....checking the config file...")


	return err
		}
	}

	fmt.Println("config/config.go...3333....checking the config file...")
	return nil
}

func (f *f) ParseEnv() error {
	for _, e := range f.envVars {

fmt.Println("confgi/config.go..Inside Parse Event.....%v.",e.name)


		if value := os.Getenv(e.name); value != "" {
			*e.target = value
fmt.Println("config/config.go....value %v", value)
		} else {


fmt.Println("config/config.go..defaultivalue .. %v", e.defaultValue)
			*e.target = e.defaultValue
		}
	}

	return nil
}

func (f *f) ParseConfig(path string) error {
	file, err := ioutil.ReadFile(path)
	 fmt.Println("..config/config.go...parsing. config... ..%v, %v",file, path)
	if err != nil {
		return err
	}

	yml := make(map[string]string)

	err = goyaml.Unmarshal(file, &yml)
fmt.Println("config/config.go .....yml...%v",yml)
	if err != nil {
		return err
	}

	f.flagSet.Visit(func(flag *flag.Flag) {
		delete(yml, flag.Name)

	fmt.Println("config/config.go..delete flag name... %v",flag.Name)

	})


  var invalid_string string


	for k, v := range yml {
		flag := f.flagSet.Lookup(k)

	fmt.Println("config/config.go...print before panic.....k, v...%v %v", k,v)
	fmt.Println("config/config.go...print flag  %v",flag)

		// if wrong parameters are passed, it will ignore and the default values will be taken.
	if flag != nil { 
		flag.Value.Set(v)
		} else { 
                invalid_string = invalid_string + "," + k
                fmt.Println("config/config.go -- Default value will be used. Not a valid Parameter : %v",k)
               }
	
		}

        if len(invalid_string) > 0 {
                return fmt.Errorf("invalid strings passed %s", invalid_string)
        }

	return nil
}
