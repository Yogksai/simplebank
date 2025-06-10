package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Yogksai/simplebank/db/mock"
	db "github.com/Yogksai/simplebank/db/sqlc"
	"github.com/Yogksai/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0, // Invalid account ID
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0) // No call to GetAccount expected
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStabs(mockStore)
			//start test server
			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err, "Error creating request")
			server.router.ServeHTTP(recorder, request)
			fmt.Printf("Response: %s, Request: %s", recorder.Body.String(), request.URL.Path)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.Randomint(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	var gotAccount db.Account
	data, err := io.ReadAll(body)
	require.NoError(t, err, "Error reading response body")
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err, "Error unmarshalling response body")
	require.Equal(t, account, gotAccount, "Account IDs do not match")
}

func TestListAccountsAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mockdb.NewMockStore(ctrl)
	accounts := []db.Account{
		randomAccount(),
		randomAccount(),
		randomAccount(),
		randomAccount(),
		randomAccount()}
	pageID := int32(1)
	pageSize := int32(5)

	mockStore.EXPECT().
		ListAccounts(gomock.Any(), db.ListAccountsParams{
			Limit:  pageSize,
			Offset: (pageID - 1) * pageSize,
		}).
		Times(1).
		Return(accounts, nil)

	server := NewServer(mockStore)

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", pageID, pageSize)
	request, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	var gotAccounts []db.Account
	err = json.Unmarshal(recorder.Body.Bytes(), &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
