package configure

import (
	"encoding/json"
	. "github.com/mchirico/go_who/util"
)

// GetSecret returns SecretStruct
func GetSecret(file string) (SecretStruct, error) {
	secStr := SecretStruct{}

	data, err := ReadFile(file)
	if err != nil {
		return secStr, err
	}
	err = json.Unmarshal([]byte(data), &secStr)
	if err != nil {
		return secStr, err
	}
	return secStr, nil

}
