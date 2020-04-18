package my_graphql

import (
	"errors"
	"mgo-gin/middlewares"
	"strings"
)

func RequireAuthorization(token string,auths...string) error {
	if token=="" {
		return errors.New("Unauthorized")
	}
	jwtToken :=strings.Split(token,"Bearer ")
	if len(jwtToken)<=1{
		return errors.New("Token is invalid")
	}
		roles := middlewares.GetRolesFromToken(jwtToken[1])
	if len(roles) <= 0 {
		return errors.New("Not have permission")
	}
	isAccessible := false
	if len(roles) < len(auths) || len(roles) == len(auths) {
		for _, auth := range auths {
			for _, role := range roles {
				if role == auth {
					isAccessible = true
					break
				}
			}
		}
	}
	if len(roles) > len(auths) {
		for _, role := range roles {
			for _, auth := range auths {
				if auth == role {
					isAccessible = true
					break
				}
			}
		}
	}
	if isAccessible == false {
		return errors.New("Not have permission")
	}
	return nil
}