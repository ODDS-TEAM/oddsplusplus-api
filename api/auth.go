package api

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/plus1/backend-go/model"
	"google.golang.org/api/oauth2/v1"
	"gopkg.in/mgo.v2/bson"
)

// Login get token decode and create user
var clientID = "822425732761-q99fsf9aue9kkf3d71dn6be83f6l45vq.apps.googleusercontent.com"

func (db *MongoDB) Login(c echo.Context) error {
	login := &model.TokenGoogle{}
	if err := c.Bind(login); err != nil {
		return err
	}

	tokenInfo, err := getInfo(login.Token)
	if err != nil {
		return err
	}
	user, firstLogin := db.CheckUser(tokenInfo)

	token, err := genToken(user)
	if err != nil {
		return err
	}

	res := &model.TokenRes{
		Token:      token,
		FirstLogin: firstLogin,
	}

	return c.JSON(http.StatusOK, res)

}

// CheckUser detect the user in database and returns the user's status.
func (db *MongoDB) CheckUser(tokenInfo *oauth2.Tokeninfo) (*model.User, bool) {
	firstLogin := false
	user, status := db.FindUser(tokenInfo)
	if status == "new" {
		// Create a new user to database
		user = db.CreateUser(tokenInfo)
		firstLogin = true
	} else if status == "notComplete" {
		firstLogin = true
	}

	return user, firstLogin
}
func (db *MongoDB) Register(c echo.Context) error {
	u := &model.User{}
	if err := c.Bind(u); err != nil {
		return err
	}
	uid := GetIDFromToken(c)
	db.UpdateUser(uid, u)
	return c.JSON(http.StatusOK, &u)
}

// UpdateUser or register information of user to database.
func (db *MongoDB) UpdateUser(uid bson.ObjectId, u *model.User) {
	q := bson.M{
		"_id":   uid,
		"email": u.Email,
	}
	ch := bson.M{
		"$set": bson.M{
			"name":   u.Name,
			"imgUrl": u.ImgUrl,
		},
	}

	if err := db.UCol.Update(q, &ch); err != nil {
		return
	}
}
func GetIDFromToken(c echo.Context) bson.ObjectId {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := claims["id"].(string)
	id := bson.ObjectIdHex(uid)
	return id
}
func genToken(user *model.User) (string, error) {
	// Set custom claims
	claims := &model.JwtCustomClaims{
		user.UserId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(
		[]byte("sMJuczqQPYzocl1s6SLj"),
	)
	if err != nil {
		return "", err
	}

	return t, nil
}

// CreateUser created by ID and email in database.
func (db *MongoDB) CreateUser(tokeninfo *oauth2.Tokeninfo) *model.User {
	u := &model.User{
		UserId: bson.NewObjectId(),
		Email:  tokeninfo.Email,
	}
	if err := db.UCol.Insert(&u); err != nil {
		return nil
	}

	return u
}
func verifyAudience(aud string) bool {
	return aud == clientID
}
func getInfo(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(&http.Client{})
	if err != nil {
		return nil, err
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfo, err := tokenInfoCall.IdToken(idToken).Do()
	if err != nil {
		return nil, err
	}

	if !verifyAudience(tokenInfo.Audience) {
		log.Println("token expire!!")
		return nil, nil
	}

	return tokenInfo, nil
}
