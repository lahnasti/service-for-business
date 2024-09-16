package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/internal/models"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725726028-team-79521/zadanie-6105/src/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTenders(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем тестовый режим для Gin
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)
	srv := &Server{
		Db:    m,
		log:   zerolog.New(os.Stdout),
		Valid: validator.New(),
	}
	r := gin.Default()
	r.GET("/api/tenders", srv.GetAllTendersHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code    int
		tenders string
	}
	type test struct {
		name    string
		request string
		filter  string
		method  string
		tender  []models.Tender
		err     error
		want    want
	}
	tests := []test{
		{
			name:    "Test 'GetAllTendersHandler' #1; Default call with service type",
			request: "/api/tenders",
			filter:  "?serviceType=it",
			method:  http.MethodGet,
			tender: []models.Tender{
				{ID: 1, Name: "tender #1", Description: "new", ServiceType: "it", Status: models.PublishedT, OrganizationID: 1, CreatorUsername: "user1", Version: 1},
				{ID: 2, Name: "tender #2", Description: "new", ServiceType: "it", Status: models.PublishedT, OrganizationID: 2, CreatorUsername: "user2", Version: 1},
			},
			want: want{
				code:    http.StatusOK,
				tenders: `{"message":"List of tenders","tenders":[{"id":1,"name":"tender #1","description":"new","serviceType":"it","status":"PUBLISHED","organizationId":1,"creatorUsername":"user1","version":1},{"id":2,"name":"tender #2","description":"new","serviceType":"it","status":"PUBLISHED","organizationId":2,"creatorUsername":"user2","version":1}]}`,
			},
		},
		{
			name:    "Test 'GetAllTendersHandler' #2; Default call without service type",
			request: "/api/tenders",
			filter:  "",
			method:  http.MethodGet,
			tender: []models.Tender{
				{ID: 1, Name: "tender #1", Description: "new", ServiceType: "it", Status: models.PublishedT, OrganizationID: 1, CreatorUsername: "user1", Version: 1},
				{ID: 2, Name: "tender #2", Description: "new", ServiceType: "beauty", Status: models.PublishedT, OrganizationID: 2, CreatorUsername: "user2", Version: 1},
			},
			want: want{
				code:    http.StatusOK,
				tenders: `{"message":"List of tenders","tenders":[{"id":1,"name":"tender #1","description":"new","serviceType":"it","status":"PUBLISHED","organizationId":1,"creatorUsername":"user1","version":1},{"id":2,"name":"tender #2","description":"new","serviceType":"beauty","status":"PUBLISHED","organizationId":2,"creatorUsername":"user2","version":1}]}`,
			},
		},
		{
			name:    "Test 'GetAllTenders' #3; No tenders found",
			request: "/api/tenders",
			filter:  "",
			method:  http.MethodGet,
			tender:  []models.Tender{},
			err:     nil,
			want: want{
				code:    http.StatusOK,
				tenders: `{"message":"List of tenders","tenders":[]}`,
			},
		},
		{
			name:    "Test 'GetAllTenders' #4; Failed to fetch tenders",
			request: "/api/tenders",
			filter:  "",
			method:  http.MethodGet,
			tender:  nil,
			err:     errors.New("db error"),
			want: want{
				code:    http.StatusInternalServerError,
				tenders: `{"message":"Failed to fetch tenders","error":"db error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.filter != "" {
				m.EXPECT().GetAllTenders(tt.filter).Return(tt.tender, tt.err)
			} else {
				m.EXPECT().GetAllTenders("").Return(tt.tender, tt.err)
			}
			srv.Db = m
			if httpSrv.URL == "" {
				t.Fatal("Test server is not running")
			}
			req := resty.New().R().SetQueryParam("serviceType", tt.filter)
			req.Method = tt.method
			req.URL = httpSrv.URL + tt.request
			resp, err := req.Send()

			assert.NoError(t, err)
			assert.JSONEq(t, tt.want.tenders, string(resp.Body()))
			assert.Equal(t, tt.want.code, resp.StatusCode())
		})
	}
}

func TestGetTendersByUser(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем тестовый режим для Gin
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)
	srv := &Server{
		Db:    m,
		log:   zerolog.New(os.Stdout),
		Valid: validator.New(),
	}
	r := gin.Default()
	r.GET("/api/tenders/my", srv.GetTendersByUser)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close() // Закрываем сервер в конце теста
	type want struct {
		code    int
		tenders string
	}
	type test struct {
		name    string
		request string
		filter  string
		method  string
		tender  []models.Tender
		err     error
		want    want
	}
	tests := []test{
		{
			name:    "Test 'GetTendersByUser' #1; Default call",
			request: "/api/tenders/my",
			filter:  "user1",
			method:  http.MethodGet,
			tender: []models.Tender{
				{ID: 1, Name: "tender #1", Description: "new", ServiceType: "it", Status: models.PublishedT, OrganizationID: 1, CreatorUsername: "user1", Version: 1},
				{ID: 2, Name: "tender #2", Description: "new", ServiceType: "beauty", Status: models.PublishedT, OrganizationID: 1, CreatorUsername: "user1", Version: 1},
			},
			want: want{
				code:    http.StatusOK,
				tenders: `[{"id":1,"name":"tender #1","description":"new","serviceType":"it","status":"PUBLISHED","organizationId":1,"creatorUsername":"user1","version":1},{"id":2,"name":"tender #2","description":"new","serviceType":"beauty","status":"PUBLISHED","organizationId":1,"creatorUsername":"user1","version":1}]`,
			},
		},
		{
			name:    "Test 'GetTendersByUser' #2; No tenders found",
			request: "/api/tenders/my",
			filter:  "user1",
			method:  http.MethodGet,
			tender:  []models.Tender{},
			err:     nil,
			want: want{
				code:    http.StatusOK,
				tenders: `[]`,
			},
		},
		{
			name:    "Test 'GetTendersByUser' #3; Failed to fetch tenders",
			request: "/api/tenders/my",
			filter:  "",
			method:  http.MethodGet,
			tender:  nil,
			err:     errors.New("invalid username"),
			want: want{
				code:    http.StatusBadRequest,
				tenders: `{"message":"Invalid params"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().GetTendersByUser(tt.filter).Return(tt.tender, tt.err)
			srv.Db = m
			req := resty.New().R().SetQueryParam("username", tt.filter)
			req.Method = tt.method
			req.URL = httpSrv.URL + tt.request
			resp, err := req.Send()

			assert.NoError(t, err)
			assert.JSONEq(t, tt.want.tenders, string(resp.Body()))
			assert.Equal(t, tt.want.code, resp.StatusCode())
		})
	}
}

func TestCreateTenderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)
	srv := &Server{
		Db:    m,
		log:   zerolog.New(os.Stdout),
		Valid: validator.New(),
	}
	r := gin.Default()
	r.POST("/api/tenders/new", srv.CreateTenderHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		body    string
		err     error
		dbFlag  bool
		want    want
	}

	tests := []test{
		{
			name:    "Test 'CreateTenderHandler' #1; Valid request",
			request: "/api/tenders/new",
			method:  http.MethodPost,
			body:    `{"name":"tender #1","description":"new","serviceType":"it","organizationId":1,"creatorUsername":"user1"}`,
			err:     nil,
			dbFlag:  true,
			want: want{
				code:   http.StatusOK,
				answer: `{"message":"Tender created successfully","tender":{"id":1,"name":"tender #1","description":"new","serviceType":"it","status":"CREATED","organizationId":1,"creatorUsername":"user1","version":1}}`,
			},
		},
		{
			name:    "Test 'CreateTenderHandler' #2; Invalid request (missing required field)",
			request: "/api/tenders/new",
			method:  http.MethodPost,
			body:    `{"name":"tender","description":"new","serviceType":"it","organizationId":1,"creatorUsername":"user1"`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"message":"Invalid request body"}`,
			},
		},
		{
			name:    "Test 'CreateTenderHandler' #3; Failed to validate request",
			request: "/api/tenders/new",
			method:  http.MethodPost,
			body:    `{"name":"tender #1","description":"new","serviceType":"it","organizationId":1}`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Key: 'Tender.CreatorUsername' Error:Field validation for 'CreatorUsername' failed on the 'required' tag"}`, // Ожидаемая ошибка валидации
			},
		},
		{
			name:    "Test 'CreateTenderHandler' #4; Failed to create tender",
			request: "/api/tenders/new",
			method:  http.MethodPost,
			body:    `{"name":"tender #1","description":"new","serviceType":"it","organizationId":1,"creatorUsername":"user1"}`,
			err:     errors.New("db error"),
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"message":"Failed to add tender","error":"db error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbFlag {
				var tender models.Tender
				if tt.want.code == http.StatusOK {
					tender = models.Tender{
						ID:              1,
						Name:            "tender #1",
						Description:     "new",
						ServiceType:     "it",
						Status:          "CREATED",
						OrganizationID:  1,
						CreatorUsername: "user1",
						Version:         1,
					}
				}
				m.EXPECT().CreateTender(gomock.Any()).Return(tender, tt.err)
			}
			req := resty.New().R()
			req.Method = tt.method
			req.Body = tt.body
			req.URL = httpSrv.URL + tt.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, resp.StatusCode())
			assert.JSONEq(t, tt.want.answer, string(resp.Body()))
		})
	}
}

func TestSetTenderStatusHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	srv := &Server{
		Db:    m,
		log:   zerolog.New(os.Stdout),
		Valid: validator.New(),
	}
	r := gin.Default()
	r.PATCH("/api/tenders/:id/status", srv.SetTenderStatusHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()
	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		body    string
		err     error
		dbFlag  bool
		want    want
	}
	tests := []test{
		{
			name:    "Test 'SetTenderStatusHandler' #1; Valid request",
			request: "/api/tenders/1/status",
			method:  http.MethodPatch,
			body:    `{"status":"PUBLISHED"}`,
			err:     nil,
			dbFlag:  true,
			want: want{
				code:   http.StatusOK,
				answer: `{"message":"Tender status updated successfully"}`,
			},
		},
		{
			name:    "Test 'SetTenderStatusHandler' #2; Invalid tender ID",
			request: "/api/tenders/abc/status",
			method:  http.MethodPatch,
			body:    `{"status":"PUBLISHED"}`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid tender ID"}`,
			},
		},
		{
			name:    "Test 'SetTenderStatusHandler' #3; Invalid request body (missing required field)",
			request: "/api/tenders/1/status",
			method:  http.MethodPatch,
			body:    `{`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid request body"}`,
			},
		},
		{
			name:    "Test 'SetTenderStatusHandler' #4; Invalid status",
			request: "/api/tenders/1/status",
			method:  http.MethodPatch,
			body:    `{"status":"OPEN"}`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid status"}`,
			},
		},
		{
			name:    "Test 'SetTenderStatusHandler' #5; Failed to update tender status",
			request: "/api/tenders/1/status",
			method:  http.MethodPatch,
			body:    `{"status":"PUBLISHED"}`,
			err:     errors.New("db error"),
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"db error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbFlag {
				m.EXPECT().SetTenderStatus(gomock.Any(), "PUBLISHED").Return(tt.err)
			}
			req := resty.New().R()
			req.Method = tt.method
			req.Body = tt.body
			req.URL = httpSrv.URL + tt.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, resp.StatusCode())
			assert.JSONEq(t, tt.want.answer, string(resp.Body()))
		})
	}
}

func TestEditTenderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	srv := &Server{
		Db:    m,
		log:   zerolog.New(os.Stdout),
		Valid: validator.New(),
	}
	r := gin.Default()
	r.PATCH("/api/:id/tenders/", srv.EditTenderHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()
	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		body    string
		err     error
		dbFlag  bool
		want    want
	}
	tests := []test{
		{
			name:    "Test 'EditTenderHandler' #1; Valid request",
			request: "/api/1/tenders/",
			method:  http.MethodPatch,
			body:    `{"name":"tender #1 updated","description":"updated"}`,
			err:     nil,
			dbFlag:  true,
			want: want{
				code:   http.StatusOK,
				answer: `{"id":1,"name":"tender #1 updated","description":"updated","serviceType":"it","status":"PUBLISHED","organizationId":1,"creatorUsername":"user1","version":2}`,
			},
		},
		{
			name:    "Test 'EditTenderHandler' #2; Invalid tender ID",
			request: "/api/abc/tenders/",
			method:  http.MethodPatch,
			body:    `{"name":"tender #1 updated","description":"updated"}`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid tender ID"}`,
			},
		},
		{
			name:    "Test 'EditTenderHandler' #3; Failed validation",
			request: "/api/1/tenders",
			method:  http.MethodPatch,
			body:    `{"description":"updated"`,
			err:     errors.New("validation error"),
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid request body"}`,
			},
		},
		{
			name:    "Test 'EditTenderHandler' #4; Invalid request body (missing required field)",
			request: "/api/1/tenders/",
			method:  http.MethodPatch,
			body:    `{`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid request body"}`,
			},
		},
		{
			name:    "Test 'EditTenderHandler' #5; Failed to update tender",
			request: "/api/1/tenders/",
			method:  http.MethodPatch,
			body:    `{"name":"tender #1 updated","description":"updated"}`,
			err:     errors.New("db error"),
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"db error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbFlag {
				tender := models.Tender{
					ID:              1,
					Name:            "tender #1 updated",
					Description:     "updated",
					ServiceType:     "it",
					Status:          "PUBLISHED",
					OrganizationID:  1,
					CreatorUsername: "user1",
					Version:         2,
				}
				m.EXPECT().EditTender(gomock.Any(), gomock.Any(), gomock.Any()).Return(tender, tt.err)
			}
			req := resty.New().R()
			req.Method = tt.method
			req.Body = tt.body
			req.URL = httpSrv.URL + tt.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, resp.StatusCode())
			assert.JSONEq(t, tt.want.answer, string(resp.Body()))
		})
	}
}

func TestRollbackTenderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockRepository(ctrl)
	srv := &Server{
		Db:    m,
		log:   zerolog.New(os.Stdout),
		Valid: validator.New(),
	}
	r := gin.Default()
	r.PUT("/api/:tenderID/rollback/:version", srv.RollbackTenderHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()
	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		err     error
		dbFlag  bool
		want    want
	}
	tests := []test{
		{
			name:    "Test 'RollbackTenderHandler' #1; Default call",
			request: "/api/1/rollback/1",
			method:  http.MethodPut,
			err:     nil,
			dbFlag:  true,
			want: want{
				code:   http.StatusOK,
				answer: `{"id":1,"name":"tender #1 updated","description":"updated","serviceType":"it","status":"PUBLISHED","organizationId":1,"creatorUsername":"user1","version":1}`,
			},
		},
		{
			name:    "Test 'RollbackTenderHandler' #2; Invalid tender ID",
			request: "/api/abc/rollback/1",
			method:  http.MethodPut,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid tender ID"}`,
			},
		},
		{
			name:    "Test 'RollbackTenderHandler' #3; Invalid version",
			request: "/api/1/rollback/abc",
			method:  http.MethodPut,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"Invalid version"}`,
			},
		},
		{
			name:    "Test 'RollbackTenderHandler' #4; Failed to rollback tender",
			request: "/api/1/rollback/1",
			method:  http.MethodPut,
			err:     errors.New("db error"),
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"db error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbFlag {
				tender := models.Tender{
					ID:              1,
					Name:            "tender #1 updated",
					Description:     "updated",
					ServiceType:     "it",
					Status:          "PUBLISHED",
					OrganizationID:  1,
					CreatorUsername: "user1",
					Version:         1,
				}
				m.EXPECT().RollbackTender(gomock.Any(), gomock.Any()).Return(tender, tt.err)
			}
			req := resty.New().R()
			req.Method = tt.method
			req.URL = httpSrv.URL + tt.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tt.want.code, resp.StatusCode())
			if tt.want.answer != "" {
				assert.JSONEq(t, tt.want.answer, string(resp.Body()))
			}
		})
	}
}
