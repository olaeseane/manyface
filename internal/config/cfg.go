package config

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

type Config struct {
	Rest   Address `valid:"required"`
	Grpc   Address `valid:"required"`
	Matrix Address `valid:"required"`
	DB     Path    `valid:"required,"`
	Log    Path    `valid:"required"`
}

type Address struct {
	Protocol string `valid:"required,in(http|https)"`
	Host     string `valid:"host,required"`
	Port     string `valid:"numeric,required"`
}

type Path struct {
	File string `valid:"required"`
}

func Read(appName string, config interface{}) error {
	v := viper.New()
	v.SetConfigName(appName)                 // name of config file (without extension)
	v.SetConfigType("yaml")                  // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath("/etc/" + appName + "/") // path to look for the config file in
	v.AddConfigPath(".")                     // optionally look for config in the working directory
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	if config != nil {
		if err := v.Unmarshal(config); err != nil {
			return fmt.Errorf("unable to unmarshall the config file: %w", err)
		}
	}

	if _, err := govalidator.ValidateStruct(config); err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			var data []byte
			for _, fld := range allErrs.Errors() {
				data = []byte(fmt.Sprintf("field: %#v\n\n", fld))
			}
			return errors.New(string(data))
		}
		return err
	}
	return nil
}
