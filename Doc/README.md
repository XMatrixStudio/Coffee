# Doc

There are docs for Coffee

## API

### Comment
    CommentController 评论

func (c *CommentController) DeleteBy(id string) (res CommonRes)
    DeleteBy DELETE /comment/{commentID} 删除指定评论

func (c *CommentController) GetBy(id string) (res CommentRes)
    GetBy GET /comment/{contentID} 获取指定内容的评论

func (c *CommentController) Post() (res CommonRes)
    Post POST /comment 增加评论

type CommentRes struct {
    State string
    Data  []services.CommentForContent
}
    CommentRes 评论回复

type CommonRes struct {
    State string
    Data  string
}
    CommonRes 返回值

### Content

type ContentController struct {
    Ctx     iris.Context
    Service services.ContentService
    Session *sessions.Session
}
    ContentController 内容

func (c *ContentController) DeleteBy(id string) (res CommonRes)
    DeleteBy DELETE /content/{contentID} 删除指定内容

func (c *ContentController) GetDetailBy(id string) (res ContentRes)
    GetDetailBy GET /content/detail/{contentID} 获取指定内容

func (c *ContentController) GetPublic() (res PublishRes)
    GetPublic GET /content/public 获取公共内容

func (c *ContentController) GetTexts() (res ContentsRes)
    GetTexts GET /content/texts 获取指定用户的所有内容

func (c *ContentController) PatchTextBy(id string) (res CommonRes)
    PatchTextBy PATCH /content/text/{contentID} 修改指定文本内容

func (c *ContentController) PostText() (res CommonRes)
    PostText POST /content/text 增加文本内容

type ContentRes struct {
    State string
    Data  models.Content
    User  services.UserBaseInfo
}
    ContentRes 内容回复

type ContentsRes struct {
    State string
    Data  []models.Content
}
    ContentsRes 内容集合回复

type LikeController struct {
    Ctx     iris.Context
    Service services.LikeService
    Session *sessions.Session
}

### Like
    LikeController Like

func (c *LikeController) Get() (res LikeRes)
    Get GET /like 获取用户点赞列表

func (c *LikeController) PatchBy(id string) (res CommonRes)
    PatchBy PATCH /like/{contentID} 取消用户对某个内容的点赞

func (c *LikeController) PostBy(id string) (res CommonRes)
    PostBy POST /like/{contentID} 对某个内容点赞

type LikeRes struct {
    State string
    Data  []string
}
    LikeRes 用户点赞数据返回值

type NotificationController struct {
    Ctx     iris.Context
    Service services.NotificationService
    Session *sessions.Session
}

### Notification

    NotificationController Like

func (c *NotificationController) DeleteBy(id string) (res CommonRes)
    DeleteBy DELETE /notificaiton/{NotificationID} 删除指定通知

func (c *NotificationController) GetAll() (res NotificationRes)
    GetAll GET /notification/all 获取用户所有通知

func (c *NotificationController) GetUnread() (res CommonRes)
    GetUnread GET /notification/unerad 获取未读通知数

func (c *NotificationController) PatchReadBy(id string) (res CommonRes)
    PatchReadBy PATCH /notification/read/{NotificationID} 标记指定通知为已读

type NotificationRes struct {
    State        string
    Notification []services.NotificationData
}
    NotificationRes 通知集合返回值

type PublishRes struct {
    State string
    Data  []services.PublishData
}
    PublishRes 公共内容返回值

### User

type UserInfoRes struct {
    State        string
    Email        string
    Name         string
    Class        int
    Info         models.UserInfo
    LikeNum      int64
    MaxSize      int64    // 存储库使用最大上限 -1为无上限 单位为KB
    UsedSize     int64    // 存储库已用大小 单位为KB
    SingleSize   int64    // 单个资源最大上限 -1为无上限
    FilesClass   []string // 文件分类
    ContentCount int64    // 内容数量
}
    UserInfoRes 用户信息返回

type UsersController struct {
    Ctx     iris.Context
    Service services.UserService
    Session *sessions.Session
}
    UsersController Users控制

func (c *UsersController) GetInfo() (res UserInfoRes)
    GetInfo GET /user/info 获取用户信息

func (c *UsersController) GetLogin() (res CommonRes)
    GetLogin GET /user/login 获取登陆页面链接

func (c *UsersController) PostInfo() (res CommonRes)
    PostInfo POST /user/info 更新用户信息

func (c *UsersController) PostLogin() (res CommonRes)
    PostLogin POST /user/login 用户登陆

func (c *UsersController) PostLogout() (res CommonRes)
    PostLogout POST /user/logout 退出登陆

func (c *UsersController) PostName() (res CommonRes)
    PostName POST /user/name 更新用户名
