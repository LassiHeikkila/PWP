package auth

import (
	"math/rand" // use math/rand to not rely on hardware to supply entropy
	"testing"
	"time"
)

func TestPasswordRoundTrip(t *testing.T) {
	plain := "correct horse battery staple"
	hashed := HashPassword(plain)
	if !PasswordEqualsHashed(plain, hashed) {
		t.Fatal("equality check failed")
	}
}

func TestBadPasswordCheck(t *testing.T) {
	realPW := "myp@ssw0rd"
	hashed := HashPassword(realPW)

	wrong := "mypassword"
	if PasswordEqualsHashed(wrong, hashed) {
		t.Fatal("equality check passed for wrong password")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	// intentionally using math/rand here // skipcq: GSC-G404
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// use 32 character passwords as test, even if it is a bit generous
	buf := make([]byte, 32)

	for n := 0; n < b.N; n++ {
		_, _ = r.Read(buf)

		h := HashPassword(string(buf))
		_ = h
	}
}

func BenchmarkComparePasswordGood(b *testing.B) {
	plain := "correct horse battery staple"
	hashed := HashPassword(plain)

	for n := 0; n < b.N; n++ {
		e := PasswordEqualsHashed(plain, hashed)
		_ = e
	}
}

func BenchmarkComparePasswordBad(b *testing.B) {
	plain := "correct horse battery staple"
	hashed := HashPassword(plain)
	bad := "something different"

	for n := 0; n < b.N; n++ {
		e := PasswordEqualsHashed(bad, hashed)
		_ = e
	}
}
