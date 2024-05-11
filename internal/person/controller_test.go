package person_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"qore-be/internal/domain/dto"
	"qore-be/internal/mocks"
	"qore-be/internal/person"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPersonController(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		ctrl, err := person.NewController(person.WithService(mocks.NewPersonService(t)))
		assert.NoError(t, err)
		assert.NotNil(t, ctrl)
	})

	t.Run("with missing person service", func(t *testing.T) {
		ctrl, err := person.NewController()
		assert.Error(t, err)
		assert.Nil(t, ctrl)
	})

	t.Run("with nil person service", func(t *testing.T) {
		ctrl, err := person.NewController(person.WithService(nil))
		assert.Error(t, err)
		assert.Nil(t, ctrl)
	})
}

func TestNewPersonController_Create(t *testing.T) {
	validReq := dto.PersonDTO{
		Name:    "name",
		Age:     15,
		Number:  "111-111-1111",
		City:    "city",
		State:   "state",
		Street1: "str1",
		Street2: "str2",
		Zip:     "1234",
	}

	cases := []struct {
		name           string
		svc            func(*testing.T) person.Service
		in             interface{}
		expectedStatus int
	}{
		{
			name: "successfully",
			svc: func(t *testing.T) person.Service {
				s := mocks.NewPersonService(t)
				s.On("Create", mock.Anything, mock.Anything).
					Return(dto.PersonDTO{Name: "name"}, nil)

				return s
			},
			in:             validReq,
			expectedStatus: 200,
		},
		{
			name: "with invalid request",
			svc: func(t *testing.T) person.Service {
				return mocks.NewPersonService(t)
			},
			in:             "name=user;phone=111-111-1111",
			expectedStatus: 400,
		},
		{
			name: "with internal error",
			svc: func(t *testing.T) person.Service {
				s := mocks.NewPersonService(t)
				s.On("Create", mock.Anything, mock.Anything).
					Return(dto.PersonDTO{}, fmt.Errorf("error"))

				return s
			},
			in:             validReq,
			expectedStatus: 500,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, err := person.NewController(person.WithService(tc.svc(t)))
			require.NoError(t, err)

			srv := gin.Default()
			gin.SetMode(gin.TestMode)

			srv.POST("/", ctrl.Create)

			rec := httptest.NewRecorder()

			data, err := json.Marshal(tc.in)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))
			require.NoError(t, err)

			srv.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}

func TestNewPersonController_GetByID(t *testing.T) {
	cases := []struct {
		name           string
		svc            func(*testing.T) person.Service
		req            string
		expectedStatus int
	}{
		{
			name: "successfully",
			svc: func(t *testing.T) person.Service {
				s := mocks.NewPersonService(t)
				s.On("GetByID", mock.Anything, mock.Anything).
					Return(&dto.PersonDTO{Name: "name"}, nil)

				return s
			},
			req:            "1",
			expectedStatus: 200,
		},
		{
			name: "with invalid request",
			svc: func(t *testing.T) person.Service {
				return mocks.NewPersonService(t)
			},
			req:            "$$",
			expectedStatus: 400,
		},
		{
			name: "with internal error",
			svc: func(t *testing.T) person.Service {
				s := mocks.NewPersonService(t)
				s.On("GetByID", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("error"))

				return s
			},
			req:            "1",
			expectedStatus: 500,
		},
		{
			name: "with not-found error",
			svc: func(t *testing.T) person.Service {
				s := mocks.NewPersonService(t)
				s.On("GetByID", mock.Anything, mock.Anything).
					Return(nil, person.ErrRecordNotFound)

				return s
			},
			req:            "1",
			expectedStatus: 404,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, err := person.NewController(person.WithService(tc.svc(t)))
			require.NoError(t, err)

			srv := gin.Default()
			gin.SetMode(gin.TestMode)

			srv.GET("/:id", ctrl.GetByID)

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/"+tc.req, nil)
			require.NoError(t, err)

			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.NotEmpty(t, resp)
		})
	}
}
