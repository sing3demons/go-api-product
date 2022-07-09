package controllers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/app/database"
	"github.com/sing3demons/app/models"
	"gorm.io/gorm"
)

//Users - receiver adater
type Users struct {
	DB *gorm.DB
}

type createUserForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type updateUserForm struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Name     string `json:"name"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email" `
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

func setUserImage(c *gin.Context, user *models.User) (imgUrl *string, err error) {
	file, err := c.FormFile("avatar")
	if file == nil || err != nil {
		return nil, err
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	// filename := strconv.Itoa(int(user.ID)) + "_" + file.Filename
	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, os.ModePerm)

	filename := path + "_" + "avatar"
	if err := c.SaveUploadedFile(file, filename); err != nil {
		return nil, err
	}

	url, err := cloudinaryUpload(filename)
	// cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: public_id})

	pwd, _ := os.Getwd()
	fmt.Println(pwd + "/" + filename)
	os.Remove(pwd + "/" + filename)

	return url, nil
}

func _setUserImage(ctx *gin.Context, user *models.User) error {
	file, _ := ctx.FormFile("avatar")
	if file == nil {
		return nil
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, os.ModePerm)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return nil
	}

	db := database.GetDB()
	user.Avatar = os.Getenv("HOST") + "/" + filename
	db.Save(user)

	return nil
}

func (u *Users) findUserByID(ctx *gin.Context) (*models.User, error) {
	id := ctx.Param("id")
	var user models.User

	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func NewCloudinary(filename string) (*cloudinary.Cloudinary, string, error) {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("API_KEY"), os.Getenv("API_SECRET"))
	if err != nil {
		return nil, "", err
	}

	public_id := "docs/sdk/go/" + filename
	return cld, public_id, nil
}

func cloudinaryUpload(filename string) (url *string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cld, public_id, err := NewCloudinary(filename)

	_, err = cld.Upload.Upload(ctx, filename, uploader.UploadParams{
		PublicID:       public_id,
		Transformation: "c_crop,g_center/q_auto/f_auto",
		Tags:           []string{"fruit"},
	})
	if err != nil {
		return nil, err
	}

	// Instantiate an object for the asset with public ID "my_image"
	my_image, err := cld.Image(public_id)
	if err != nil {
		return nil, err
	}
	fmt.Println(my_image.String())
	result, err := my_image.String()
	if err != nil {
		return nil, err
	}

	return &result, nil
}
