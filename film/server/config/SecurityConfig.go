package config

const PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBPAIBAAJBANNnshoUaT2gFNrihmFdmC1cBCs1XLFc5Fn3MfNOR3aOGDO0ohXl
bku6Ir/qITN/yeH5pY34WEcETet3YhESpE8CAwEAAQJBAI7Ekdfg/u26RTtJDd2F
WrcPVFVl1TKGfERxl08sB0D9HLvUSBfAEg/UpfWSQ57aSJ9b0gVKmDhgF8FymuUV
v2kCIQDzXXSZ/oeKmqObwad0Fa82IFof3LeZdpbrjyz3w45JDQIhAN5hdmuW+y2w
UgSy0o4zGFsEG/RBZsvVnSSfkdR47dPLAiEA2XbPNLQu5fnc7NeVDLQ7xsAOCJ6w
KR/BKGjeI9/JCxkCIQCjMkU0ec2FXxMhzZXFs2uZR6+4FdL5nZ9ABDaCBekK9wIg
XEfd11qabi9jPrbsOVNZCTk51B7Ug0ZwGyn0BA8Jlo0=
-----END RSA PRIVATE KEY-----
`

const PublicKey = `-----BEGIN RSA PUBLIC KEY-----
MEgCQQDTZ7IaFGk9oBTa4oZhXZgtXAQrNVyxXORZ9zHzTkd2jhgztKIV5W5LuiK/
6iEzf8nh+aWN+FhHBE3rd2IREqRPAgMBAAE=
-----END RSA PUBLIC KEY-----
`

const (
	Issuer           = "GoFilm"
	AuthTokenExpires = 10 * 24 // 单位 h
	UserTokenKey     = "User:Token:%d"
)
