package server_logic

import (
	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/azazel3ooo/keeper/internal/models/testing_repos_server"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCheckToken(t *testing.T) {
	expiration := 1.0
	realID := "1"
	testToken, _ := GenerateToken(realID, expiration)
	expiredToken, _ := GenerateToken(realID, 0.0)

	tests := []struct {
		description string
		token       string
		want        string
		wantErr     error
	}{
		{
			description: "real token",
			token:       testToken,
			want:        realID,
			wantErr:     nil,
		},
		{
			description: "bad token",
			token:       testToken[:len(testToken)-5],
			want:        "",
			wantErr:     models.ErrInvalidToken,
		},
		{
			description: "bad token",
			token:       expiredToken,
			want:        "",
			wantErr:     models.ErrExpiredToken,
		},
	}

	for _, tt := range tests {
		id, err := CheckToken(tt.token)
		assert.Equalf(t, tt.want, id, tt.description)
		assert.Equalf(t, tt.wantErr, err, tt.description)
	}
}

func TestCheckUser(t *testing.T) {
	var storage testing_repos_server.TestingServerStorage

	storage.Init()
	user := "test_user"
	userPas := "test_pas"
	userID, err := storage.CreateUser(user, userPas)
	if err != nil {
		log.Println(err)
	}

	tests := []struct {
		description string
		req         models.UserRequest
		want        string
		wantErr     error
	}{
		{
			description: "success check",
			req:         models.UserRequest{Login: user, Password: userPas},
			want:        userID,
			wantErr:     nil,
		},
		{
			description: "data conflict check",
			req:         models.UserRequest{Login: user, Password: userPas[:len(userPas)-1]},
			want:        "",
			wantErr:     models.ErrUserDataConflict,
		},
	}
	for _, tt := range tests {
		res, err := CheckUser(tt.req, storage)
		assert.Equalf(t, tt.want, res, tt.description)
		assert.Equalf(t, tt.wantErr, err, tt.description)
	}
}

func TestDelete(t *testing.T) {
	var s testing_repos_server.TestingServerStorage
	s.Init()

	id := "tmp_ID"
	u := "tmp_u"
	s.SetData(models.UserData{ID: id, Data: "asfsdaf", Comment: "asdf"}, u)

	tests := []struct {
		description string
		req         models.DeleteRequest
		user        string
		wantErr     bool
	}{
		{
			description: "success delete",
			req:         models.DeleteRequest{ID: id},
			user:        u,
			wantErr:     false,
		},
		{
			description: "unknown id delete",
			req:         models.DeleteRequest{ID: id[:len(id)-1]},
			user:        u,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		err := Delete(tt.req, s, tt.user)
		assert.Equalf(t, tt.wantErr, err != nil, tt.description)
	}
}

func TestRegistration(t *testing.T) {
	var storage testing_repos_server.TestingServerStorage

	storage.Init()
	user := "test_user"
	userPas := "test_pas"

	tests := []struct {
		description string
		req         models.UserRequest
		wantErr     error
	}{
		{
			description: "success check",
			req:         models.UserRequest{Login: user, Password: userPas},
			wantErr:     nil,
		},
		{
			description: "data conflict check",
			req:         models.UserRequest{Login: user, Password: userPas},
			wantErr:     models.ErrUserConflict,
		},
	}
	for _, tt := range tests {
		_, err := Registration(tt.req, storage)
		assert.Equalf(t, tt.wantErr, err, tt.description)
	}
}

func TestSet(t *testing.T) {
	var s testing_repos_server.TestingServerStorage
	s.Init()

	id := "tmp_ID"
	u := "tmp_u"

	tests := []struct {
		description string
		req         models.UserData
		user        string
		wantErr     bool
	}{
		{
			description: "success set",
			req:         models.UserData{ID: id, Data: "tmp_data", Comment: "dsgdfg"},
			user:        u,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		err := Set(tt.req, s, tt.user)
		assert.Equalf(t, tt.wantErr, err != nil, tt.description)
	}
}

func TestUpdate(t *testing.T) {
	var s testing_repos_server.TestingServerStorage
	s.Init()

	id := "tmp_ID"
	u := "tmp_u"
	s.SetData(models.UserData{ID: id, Data: "asfsdaf", Comment: "asdf"}, u)
	prevRes, _ := s.GetData(u)

	tests := []struct {
		description string
		req         models.UserData
		user        string
		wantErr     bool
	}{
		{
			description: "success update",
			req:         models.UserData{ID: id, Data: "qqq", Comment: ""},
			user:        u,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		err := Update(tt.req, s, tt.user)
		assert.Equalf(t, tt.wantErr, err != nil, tt.description)
		res, _ := s.GetData(tt.user)
		assert.Equalf(t, true, prevRes[0].Data != res[0].Data, tt.description)
	}
}
