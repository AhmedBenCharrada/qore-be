package person_test

import (
	"context"
	"fmt"
	"qore-be/internal/domain/dto"
	"qore-be/internal/mocks"
	"qore-be/internal/person"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPersonService(t *testing.T) {
	t.Run("successfully", func(t *testing.T) {
		ctrl, err := person.NewService(person.WithRepository(&person.Repo{}))
		assert.NoError(t, err)
		assert.NotNil(t, ctrl)
	})

	t.Run("with opt initializer error", func(t *testing.T) {
		ctrl, err := person.NewService(func(si *person.ServiceImpl) error {
			return fmt.Errorf("error")
		})
		assert.Error(t, err)
		assert.Nil(t, ctrl)
	})

	t.Run("with missing person repo", func(t *testing.T) {
		ctrl, err := person.NewService()
		assert.Error(t, err)
		assert.Nil(t, ctrl)
	})

	t.Run("with nil person repo", func(t *testing.T) {
		ctrl, err := person.NewService(person.WithRepository(nil))
		assert.Error(t, err)
		assert.Nil(t, ctrl)
	})
}

func TestPersonService_Create(t *testing.T) {
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
		name   string
		db     func(*testing.T) person.Repository
		in     dto.PersonDTO
		hasErr bool
	}{
		{
			name: "successfully",
			db: func(t *testing.T) person.Repository {
				d := mocks.NewPersonRepository(t)
				d.On("Add", mock.Anything, mock.Anything).
					Return(dto.PersonDTO{Name: "name"}, nil)

				return d
			},
			in: validReq,
		},
		{
			name: "with error",
			db: func(t *testing.T) person.Repository {
				d := mocks.NewPersonRepository(t)
				d.On("Add", mock.Anything, mock.Anything).
					Return(dto.PersonDTO{}, fmt.Errorf("error"))

				return d
			},
			in:     validReq,
			hasErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := person.NewService(person.WithRepository(tc.db(t)))
			require.NoError(t, err)

			res, err := svc.Create(context.TODO(), tc.in)
			assert.Equal(t, !tc.hasErr, err == nil)
			if !tc.hasErr {
				assert.NotEmpty(t, res)
			}
		})
	}
}

func TestPersonService_GetByID(t *testing.T) {
	cases := []struct {
		name   string
		db     func(*testing.T) person.Repository
		in     int
		hasErr bool
	}{
		{
			name: "successfully",
			db: func(t *testing.T) person.Repository {
				d := mocks.NewPersonRepository(t)
				d.On("GetByID", mock.Anything, mock.Anything).
					Return(dto.PersonDTO{Name: "name"}, nil)

				return d
			},
			in: 1,
		},
		{
			name: "with error",
			db: func(t *testing.T) person.Repository {
				d := mocks.NewPersonRepository(t)
				d.On("GetByID", mock.Anything, mock.Anything).
					Return(dto.PersonDTO{}, fmt.Errorf("error"))

				return d
			},
			in:     1,
			hasErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := person.NewService(person.WithRepository(tc.db(t)))
			require.NoError(t, err)

			res, err := svc.GetByID(context.TODO(), tc.in)
			assert.Equal(t, !tc.hasErr, err == nil)
			if !tc.hasErr {
				assert.NotEmpty(t, res)
			}
		})
	}

}
