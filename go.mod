module github.com/colt3k/nglog

go 1.13

require (
	github.com/colt3k/utils/archive v0.0.9
	github.com/go-mail/mail v2.3.1+incompatible
	github.com/mattn/go-isatty v0.0.20
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.17.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1 // indirect
)

replace golang.org/x/net => golang.org/x/net v0.19.0 //CVE-2023-48795
