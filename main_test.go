package main

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	datepicker "github.com/ethanefung/bubble-datepicker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedTable struct {
	mock.Mock
}

func (t *MockedTable) AddRow(s ...string) {
	t.Called(s)
}

func (t *MockedTable) DeleteRows() {
	t.Called()
}

func (t *MockedTable) SetHeader(s ...string) {
	t.Called(s)
}

func (t *MockedTable) SetPadding(i int) error {
	args := t.Called(i)
	return args.Error(0)
}

func (t *MockedTable) SetRowSeparator(s string) {
	t.Called(s)
}

func (t *MockedTable) SetColumnSeparator(s string) {
	t.Called(s)
}

func (t *MockedTable) SetHeaderSeparator(s string) {
	t.Called(s)
}

func (t *MockedTable) ToString() string {
	args := t.Called()
	return args.String(0)
}

func TestModel_View(t *testing.T) {
	mockedTable := new(MockedTable)
	cases := []struct {
		name                string
		m                   Model
		tableToStringReturn string
		expected            string
	}{
		{
			name: "view overview no users",
			m: Model{
				users:         nil,
				cursor:        0,
				table:         mockedTable,
				currentScreen: screenOverview,
				err:           nil,
			},
			expected: "No users yet. Press 'a' to add one.\n\nKeys: ↑/↓ move | a add | Enter Add Date | q quit\n",
		},
		{
			name: "view overview with users",
			m: Model{
				users:         Users{User{"Marc", nil}, User{"Spongebob", nil}},
				cursor:        0,
				table:         mockedTable,
				currentScreen: screenOverview,
				err:           nil,
			},
			tableToStringReturn: "Some user data",
			expected:            "Some user data\nKeys: ↑/↓ move | a add | Enter Add Date | q quit\n",
		},
		{
			name: "view overview error",
			m: Model{
				users:         Users{User{"Marc", nil}, User{"Spongebob", nil}},
				cursor:        0,
				table:         mockedTable,
				currentScreen: screenOverview,
				err:           fmt.Errorf("something went wrong"),
			},
			tableToStringReturn: "Some user data",
			expected:            "⚠ something went wrong\n\nSome user data\nKeys: ↑/↓ move | a add | Enter Add Date | q quit\n",
		},
		{
			name: "view add user",
			m: Model{
				users:         Users{User{"Marc", nil}, User{"Spongebob", nil}},
				cursor:        0,
				input:         textinput.Model{},
				currentScreen: screenAddUser,
				err:           nil,
			},
			expected: "Add user\n\nName:\n \n\nEnter save | Esc cancel\n",
		},
		{
			name: "view add date",
			m: Model{
				users:         Users{User{Name: "Patrick"}},
				cursor:        0,
				datepicker:    datepicker.Model{},
				currentScreen: screenAddDate,
				err:           nil,
			},
			expected: "Add Date for Patrick:\n    January 1  \n              \nSuMoTuWeThFrSa\n  010203040506\n07080910111213\n14151617181920\n21222324252627\n28293031      \nEnter save | Esc/q cancel",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.tableToStringReturn != "" {
				mockedTable.On("ToString").Return(c.tableToStringReturn)
				defer mockedTable.On("ToString").Unset()
			}
			result := c.m.View()
			assert.Equal(t, c.expected, result)
		})
	}
}

func TestModel_Init(t *testing.T) {
	m := Model{}

	result := m.Init()

	assert.Nil(t, result)
}

func TestUpdateTable(t *testing.T) {
	mockedTable := new(MockedTable)
	cases := []struct {
		name string
		m    Model
	}{
		{
			name: "Update table no users",
			m: Model{
				users: Users{},
				table: mockedTable,
			},
		},
		{
			name: "Update table with users",
			m: Model{
				users: Users{User{"Spongebob", nil}, User{"Patrick", nil}},
				table: mockedTable,
			},
		},
	}
	mockedTable.On("DeleteRows").Return()
	mockedTable.On("AddRow", mock.Anything, mock.Anything).Return()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.m.updateTable()
			mockedTable.AssertCalled(t, "DeleteRows")
			mockedTable.AssertNumberOfCalls(t, "AddRow", len(c.m.users))
		})
	}
}
