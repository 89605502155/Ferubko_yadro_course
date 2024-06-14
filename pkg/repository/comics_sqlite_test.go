package repository

import (
	"errors"
	"fmt"
	"testing"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"xkcd/pkg/xkcd"
)

func TestGenerate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	r := NewComicsSQLite(db)

	type mockBechavior func(data map[string]xkcd.ComicsInfo)
	firstsMap := map[string]xkcd.ComicsInfo{}
	firstsMap["1"] = xkcd.ComicsInfo{
		Url:      "https://google.com",
		Keywords: []string{"google", "yadro", "podolsky"},
	}
	// secondMap := map[string]xkcd.ComicsInfo{}
	// secondMap[""] = xkcd.ComicsInfo{
	// 	Url:      "http://me.to.com",
	// 	Keywords: []string{"vsxcode", "kdp", "liNO3"},
	// }
	testTable := []struct {
		name          string
		data          map[string]xkcd.ComicsInfo
		mockBechavior mockBechavior
		expectedError error
	}{
		{
			name: "OK",
			data: firstsMap,
			mockBechavior: func(data map[string]xkcd.ComicsInfo) {
				mock.ExpectBegin()
				insertQuery := fmt.Sprintf("INSERT INTO %s (comics_id, url, keywords) VALUES ", comicsTable)
				values := ""
				for key, value := range data {
					for _, v := range value.Keywords {
						values += fmt.Sprintf("('%s', '%s','%s'),", key, value.Url, v)
					}
				}
				if len(values) > 0 {
					values = values[:len(values)-1]
				}
				insertQuery += values
				mock.ExpectExec(insertQuery).WillReturnResult(sqlxmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
	}
	for _, test := range testTable {
		t.Run(
			test.name, func(t *testing.T) {
				test.mockBechavior(test.data)
				err := r.Generate(test.data)
				if err != nil && errors.Is(err, test.expectedError) {
				}
			},
		)
	}
}
