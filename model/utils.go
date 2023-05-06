package msi

type FileCredential struct {
	Key       string
	Path      string
	ShortName string
}

func (c *Component) SetKeyPath(name string) bool {
	if f, ok := c.Files[name]; ok {
		c.KeyPath = OptString(f.Key)
		return true
	}
	return false
}

func (c *Component) GetFileKey(name string) string {
	if f, ok := c.Files[name]; ok {
		return f.Key
	}
	return ""
}
