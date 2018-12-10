package param

import (
	"net/url"
	"testing"
)

// Common request params
type Common struct {
	AppName string
	Sig     string
	Source  string
}

// SmsSend sms end params
type SmsSend struct {
	Common
	Mobile string `param:"mobile,13500000000"`
	Msg    string
	Age    int `param:"age,18"`
}

var req = url.Values{
	"appName": {"usercenter"},
	"sig":     {"sig1111"},
	"source":  {"mix"},
	"mobile":  {"13512341234"},
	"msg":     {"message info"},
	// "age":     {"1"},
}

func TestUnmarshal(t *testing.T) {
	o := new(SmsSend)
	err := Unmarshal(req, o)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", o)
}

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < 100000; i++ {
		o := new(SmsSend)
		err := Unmarshal(req, o)
		if err != nil {
			b.Error(err)
		}
		// b.Logf("%+v", o)
		// 1.779 s
	}
}
