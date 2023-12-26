package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/indramhrt/simplebank/db/mock"
	db "github.com/indramhrt/simplebank/db/sqlc"
	"github.com/indramhrt/simplebank/util"
	"github.com/stretchr/testify/require"
)

func UserEq(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParams{arg, password}
}

func (e eqCreateUserParams) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParams) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

type eqCreateUserParams struct {
	arg      db.CreateUserParams
	password string
}

func TestCreateUserAPI(t *testing.T) {
	user, password := createRandomUser(t)
	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				hashedPassword, err := util.HashedPassword(password)
				require.NoError(t, err)

				arg := db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: hashedPassword,
					FullName:       user.FullName,
					Email:          user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), UserEq(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(10)
	hashedPassword, err := util.HashedPassword(password)
	require.Empty(t, err)

	return db.User{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
