package cookie

import "fmt"

type Cookie struct {
	name  string
	value string
	//if maxAge > 0, then the cookie will be invalid after maxAge seconds
	//if maxAge < 0, then the cookie is a temporary cookie, it only be valid before the current window is closed
	//if maxAge = 0, then browser will delete the cookie
	maxAge int64
}

func NewCookie(name string, value string) *Cookie {
	return &Cookie{
		name:  name,
		value: value,
	}
}

func (c *Cookie) GetName() string {
	return c.name
}

func (c *Cookie) GetValue() string {
	return c.value
}

func (c *Cookie) SetMaxAge(age int64) {
	c.maxAge = age
}

//Response 的 Header 中 Set-Cookie 需要满足下面的格式，任选一种即可。
//Set-Cookie: <cookie-name>=<cookie-value>
//Set-Cookie: <cookie-name>=<cookie-value>;Expire=<date>
//Set-Cookie: <cookie-name>=<cookie-value>;Max-Age=<number>
//Set-Cookie: <cookie-name>=<cookie-value>;Domain=<domain-value>
//Set-Cookie: <cookie-name>=<cookie-value>;Path=<path-value>
//Set-Cookie: <cookie-name>=<cookie-value>;Secure
func (c *Cookie) ToCookieHeaderName() string {
	return "Set-Cookie"
}

func (c *Cookie) ToCookieHeaderValue() string {
	return fmt.Sprintf("%s=%s;Max-Aget=%d", c.name, c.value, c.maxAge)
}
