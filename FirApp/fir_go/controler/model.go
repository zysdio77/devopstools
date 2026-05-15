package controler

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func SelectDb(ctx context.Context, name *gvar.Var) (string, error) {
	aa, err := g.Model("fir").Where("system_type=?", name).All()
	if err != nil {
		g.Log().Error(ctx, err)
		return "", err
	}
	return aa.Json(), nil
}

func InsertDb(ctx context.Context, info *Infos) (int64, error) {
	result, err := g.Model("fir").Data(info).Insert()
	if err != nil {
		g.Log().Error(ctx, err)
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteDb(ctx context.Context, name string) (int64, error) {
	r, err := g.Model("fir").Where("name", name).Delete()
	if err != nil {
		g.Log().Error(ctx, err)
		return 0, err
	}
	return r.RowsAffected()
}

func PagesData(ctx context.Context, name *gvar.Var, page *gvar.Var) (int, int, string, error) {
	p := gconv.Int(page)
	c, err := g.Model("fir").Where("system_type=?", name).Count()
	if err != nil {
		g.Log().Error(ctx, err)
		return 0, 0, "", err
	}

	perPage := 10
	pages := c/perPage + 1
	if c%perPage == 0 && c > 0 {
		pages = c / perPage
	}
	if pages < 1 {
		pages = 1
	}

	start := (p - 1) * perPage
	aa, err := g.Model("fir").Where("system_type=?", name).
		OrderDesc("create_time").Limit(start, perPage).All()
	if err != nil {
		g.Log().Error(ctx, err)
		return 0, 0, "", err
	}

	return c, pages, aa.Json(), nil
}

func UpdateDb(ctx context.Context, name string, systemtype string) (int64, error) {
	result, err := g.Model("fir").Data(g.Map{"system_type": systemtype}).
		Where("name", name).Update()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}
