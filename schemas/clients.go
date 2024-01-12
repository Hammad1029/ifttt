package schemas

type AddClient struct {
	Name string `json:"name" binding:"required"`
}
