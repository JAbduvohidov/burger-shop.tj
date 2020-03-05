package services

import "os"

func FlagOrEnv(fl *string, envName string) (value string, ok bool) {
	if fl != nil {
		return *fl, true
	}

	return os.LookupEnv(envName)
}

//postgres://nurzyxgxduryxt:a91147d43b56869a99a0815d324323f5f22071d6dfa17cdd789c93388a392072@ec2-52-86-73-86.compute-1.amazonaws.com:5432/dc1ns5rpr9g4e3
//postgres://user:pass@localhost:5432/app
