# url.Values parse
Parse a url.Values to a given struct.

```golang
var Login struct {
    Username string 
    Password string
    Remember bool `param:"remember,true"` // true means default set to true if the field not exists in post form data
}

loginForm = new(Login)
req := url.Values{
    "username": {"Leon"},
    "password": {"admin"},
}

if err := param.Unmarshal(req, loginForm); err != nil {
    // handler error
}

// &loginForm{Username: "Leon", Password: "admin", Remember: true}
```
