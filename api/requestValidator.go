package api

import (
	"errors"
)

func validateTemplateId(templateId string) error {
	if templateId == "" {
		return errors.New("templateId can't be empty!")
	}
	return nil
}

func validateUuid(uuid string) error {
	if uuid == "" {
		return errors.New("uuid can't be empty!")
	}
	if len(uuid) < 15 {
		return errors.New("instanceId has to be longer than 15 characters!")
	}
	return nil
}
