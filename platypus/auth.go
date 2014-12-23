package platypus

import (
	"errors"
)

func (p Platypus) Login(username string, password string) error {
	params := Parameters{
		Logintype: "Staff",
		Datatype:  "XML",
		Username:  username,
		Password:  password,
	}

	res, err := p.Exec("Login", params)
	if err != nil {
		return err
	}

	if res.Success == 0 {
		return errors.New(res.ResponseText)
	}

	return nil
}
