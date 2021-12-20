package config

import "github.com/spf13/viper"

// JsonConfigModleMapping mapping the given json config file to the given struct.
// Using `json:"<name>"` tag to customize the field name
func JsonConfigModleMapping(stc interface{}, filePath string) error {
	vp := viper.New()
	vp.SetConfigFile(filePath)
	vp.SetConfigType("json")
	err := vp.ReadInConfig()
	if err != nil {
		return err
	}
	err = vp.Unmarshal(stc)
	if err != nil {
		return err
	}
	return nil
}
