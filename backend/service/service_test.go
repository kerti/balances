package service_test

import "github.com/kerti/balances/backend/model"

func getDefaultPageInfo() model.PageInfoOutput {
	return model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}
}
