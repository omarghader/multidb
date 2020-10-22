package multidb

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func ToMap(input interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("%s: %s\n", EXCEPTION_JSON_MARSHAL, err)
		return nil, err
	}

	err = json.Unmarshal(inputBytes, &res)
	if err != nil {
		logrus.Errorf("%s: %s\n", EXCEPTION_JSON_UNMARSHAL, err)
		return nil, err
	}

	return res, nil
}

func ToMapString(input interface{}) (map[string]string, error) {
	res := map[string]string{}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("%s: %s\n", EXCEPTION_JSON_MARSHAL, err)
		return nil, err
	}

	err = json.Unmarshal(inputBytes, &res)
	if err != nil {
		logrus.Errorf("%s: %s\n", EXCEPTION_JSON_UNMARSHAL, err)
		return nil, err
	}

	return res, nil
}

//------------------------------------------------------------
func ToStruct(input interface{}, res interface{}) error {

	inputBytes, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("%s: %s\n", EXCEPTION_JSON_MARSHAL, err)
		return err
	}

	err = json.Unmarshal(inputBytes, res)
	if err != nil {
		logrus.Errorf("%s: %s\n", EXCEPTION_JSON_UNMARSHAL, err)
		return err
	}

	return nil
}
