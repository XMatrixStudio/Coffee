package controllers

func (c *ContentController) PostAlbum() (res CommonRes) {
	c.Ctx.SetMaxRequestBodySize(1024 * 1024 * 1024) // 1G
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	if err := c.Service.AddAlbum(c.Ctx, c.Session.GetString("id")); err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}