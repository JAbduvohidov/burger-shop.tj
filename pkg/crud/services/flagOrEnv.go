package services

import "os"

func FlagOrEnv(flag *string, envKey string) (string, bool) {
	if flag != nil {
		return *flag, true
	}
	return os.LookupEnv(envKey)
}
