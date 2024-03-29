package tests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/service/app/services/sales-api/handlers/crud/productgrp"
	"github.com/ardanlabs/service/business/data/dbtest"
	"github.com/ardanlabs/service/business/web/errs"
	"github.com/google/go-cmp/cmp"
)

func productUpdate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/products/%s", sd.Users[0].Products[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Model: &productgrp.AppUpdateProduct{
				Name:     dbtest.StringPointer("Guitar"),
				Cost:     dbtest.FloatPointer(10.34),
				Quantity: dbtest.IntPointer(10),
			},
			Resp: &productgrp.AppProduct{},
			ExpResp: &productgrp.AppProduct{
				ID:          sd.Users[0].Products[0].ID.String(),
				UserID:      sd.Users[0].ID.String(),
				Name:        "Guitar",
				Cost:        10.34,
				Quantity:    10,
				DateCreated: sd.Users[0].Products[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].Products[0].DateCreated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*productgrp.AppProduct)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*productgrp.AppProduct)
				gotResp.DateUpdated = expResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func productUpdate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:       "bad-input",
			URL:        fmt.Sprintf("/v1/products/%s", sd.Users[0].Products[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Model: &productgrp.AppUpdateProduct{
				Cost:     dbtest.FloatPointer(-1.0),
				Quantity: dbtest.IntPointer(0),
			},
			Resp: &errs.Response{},
			ExpResp: &errs.Response{
				Error:  "data validation error",
				Fields: map[string]string{"cost": "cost must be 0 or greater", "quantity": "quantity must be 1 or greater"},
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func productUpdate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/products/%s", sd.Users[0].Products[0].ID),
			Token:      "",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Resp:       &errs.Response{},
			ExpResp:    &errs.Response{Error: "Unauthorized"},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        fmt.Sprintf("/v1/products/%s", sd.Users[0].Products[0].ID),
			Token:      sd.Users[0].Token + "A",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Resp:       &errs.Response{},
			ExpResp:    &errs.Response{Error: "Unauthorized"},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        fmt.Sprintf("/v1/products/%s", sd.Admins[0].Products[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Model: &productgrp.AppUpdateProduct{
				Name:     dbtest.StringPointer("Guitar"),
				Cost:     dbtest.FloatPointer(10.34),
				Quantity: dbtest.IntPointer(10),
			},
			Resp:    &errs.Response{},
			ExpResp: &errs.Response{Error: "Unauthorized"},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
