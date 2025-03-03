package req

type User struct {
	Username string `json:"username" binding:"required,max=16,min=2" msg:"用户名6~12位"`
	Password string `json:"password" binding:"required,max=32,min=6" msg:"密码6~32位"`
	AvatarId string `json:"avatar_id" binding:"required,numeric" msg:"头像id请输入数字"`
}
