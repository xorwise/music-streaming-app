package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ssov1 "github.com/xorwise/music-streaming-service/gen"
	"github.com/xorwise/music-streaming-service/tests/auth/suite"
)

const passDefaultLen = 10

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	username := gofakeit.Username()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Username: username,
		Password: pass,
	})

	require.NoError(t, err)
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetId(), int64(claims["uid"].(float64)))
	assert.Equal(t, username, claims["username"].(string))
	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegister_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	username := gofakeit.Username()
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: pass,
	})

	require.Error(t, err)
	assert.Empty(t, respReg.GetId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		username    string
		password    string
		expectedErr string
	}{
		{
			name:        "empty username",
			username:    "",
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			expectedErr: "username and password are required",
		},
		{
			name:        "empty password",
			username:    gofakeit.Username(),
			password:    "",
			expectedErr: "username and password are required",
		},
		{
			name:        "empty username and password",
			username:    "",
			password:    "",
			expectedErr: "username and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Username: tt.username,
				Password: tt.password,
			})

			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		username    string
		password    string
		expectedErr string
	}{
		{
			name:        "empty username",
			username:    "",
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			expectedErr: "username and password are required",
		},
		{
			name:        "empty password",
			username:    gofakeit.Username(),
			password:    "",
			expectedErr: "username and password are required",
		},
		{
			name:        "empty username and password",
			username:    "",
			password:    "",
			expectedErr: "username and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Username: tt.username,
				Password: tt.password,
			})

			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
