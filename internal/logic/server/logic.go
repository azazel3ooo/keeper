package server_logic

import (
	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/golang-jwt/jwt"
	"time"
)

const JWTSalt = "super_secret_salt"

// Registration выполняет регистрацию пользователя по данным, переданным в models.UserRequest, в хранилище models.Storable4Server
// возвращает id созданного пользователя
func Registration(request models.UserRequest, s models.Storable4Server) (string, error) {
	id, err := s.CreateUser(request.Login, request.Password)
	if id == "exist" {
		err = models.ErrUserConflict
	}

	return id, err
}

// GenerateToken по переданному id создает JWT. Опционально можно задать необходимую длительность жизни токена (по умолчанию 5 минут)
func GenerateToken(id string, duration ...float64) (string, error) {
	var val float64
	if len(duration) == 1 {
		val = duration[0]
	} else {
		val = 5
	}

	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Duration(val) * time.Minute).Unix(),
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(JWTSalt)) // salt for jwt
	if err != nil {
		return "", err
	}

	return token, nil
}

// CheckToken проверяет корректность переданного JWT и возвращает значение поля id из этого токена
func CheckToken(token string) (string, error) {
	type myCl struct {
		jwt.StandardClaims
		Id string `json:"id"`
	}

	cl := myCl{}

	t, err := jwt.ParseWithClaims(token, &cl, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return "", models.ErrInvalidToken
		}
		return []byte(JWTSalt), nil
	})
	if err != nil || !t.Valid {
		return "", models.ErrInvalidToken
	}

	if cl, ok := t.Claims.(*myCl); ok {
		if cl.ExpiresAt <= time.Now().Unix() {
			return "", models.ErrExpiredToken
		}
	} else {
		return "", models.ErrInvalidToken
	}

	return cl.Id, nil
}

// CheckUser проверяет соответствие пароля и логина. Возвращает id пользователя
func CheckUser(request models.UserRequest, s models.Storable4Server) (string, error) {
	id, pas, err := s.CheckUser(request.Login)
	if err != nil {
		return "", err
	}
	if pas != request.Password {
		return "", models.ErrUserDataConflict
	}

	return id, nil
}

// GetAll возвращает массив записей, где id владельца == переданному id
func GetAll(id string, s models.Storable4Server) ([]models.UserData, error) {
	data, err := s.GetData(id)
	if err != nil {
		return nil, err
	}

	// не в одну строку, поскольку оставляется место под crypto(но это не точно)

	return data, nil
}

// Set добавляет данные из models.UserData для пользователя user
func Set(req models.UserData, s models.Storable4Server, user string) error {
	err := s.SetData(req, user)
	if err != nil {
		return err
	}

	return nil
}

// Delete удаляет данные из models.DeleteRequest для пользователя user
func Delete(req models.DeleteRequest, s models.Storable4Server, user string) error {
	err := s.Delete(req, user)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет данные из models.UserData для пользователя user
func Update(req models.UserData, s models.Storable4Server, user string) error {
	err := s.Update(req, user)
	if err != nil {
		return err
	}

	return nil
}
