package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	// this is a factor of how much computation it takes to create the hash
	// so its also a factor of how hard it is to crack
	// default is 10, but let's increase it a bit to make cracking slower
	// it does slow down pw verification for us as well, but its acceptable.

	// with cost = 14, hashing takes about 800-900ms on developers laptop.
	// on a low power VPS it might take too long

	/*
	   $ go test -run="Benchmark.*" -bench=. -benchtime 10s
	   goos: linux
	   goarch: amd64
	   pkg: github.com/LassiHeikkila/taskey/internal/auth
	   cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
	   BenchmarkHashPassword-4          	      12	 847085260 ns/op
	   BenchmarkComparePasswordGood-4   	      12	 911934717 ns/op
	   BenchmarkComparePasswordBad-4    	      12	 914731851 ns/op
	   PASS
	   ok  	github.com/LassiHeikkila/taskey/internal/auth	48.281s
	*/

	// TODO: decrease if its too slow on target hardware
	const cost = 14
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return ""
	}
	return string(b)
}

func PasswordEqualsHashed(plain string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
