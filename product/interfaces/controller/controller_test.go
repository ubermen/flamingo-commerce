package controller

import (
	"context"
	"errors"
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/testutil"
	"flamingo/framework/web"
	"net/url"
	"testing"
)

type (
	MockProductService struct{}
)

func (mps *MockProductService) Get(ctx context.Context, marketplacecode string) (domain.BasicProduct, error) {
	if marketplacecode == "fail" {
		return nil, errors.New("fail")
	}

	return domain.SimpleProduct{
		BasicProductData: domain.BasicProductData{Title: "My Product Title", MarketPlaceCode: marketplacecode},
	}, nil
}

func TestViewController_Get(t *testing.T) {
	var redirectedTo, redirectedName string
	var tplname string
	var errorHappened bool

	vc := &View{
		ProductService: new(MockProductService),
		RedirectAware: &testutil.MockRedirectAware{
			CbRedirect: func(name string, args map[string]string) web.Response {
				redirectedTo = "product.view"
				redirectedName = args["name"]
				return nil
			},
		},
		RenderAware: &testutil.MockRenderAware{
			CbRender: func(context web.Context, tpl string, data interface{}) web.Response {
				tplname = tpl
				return nil
			},
		},
		ErrorAware: &testutil.MockErrorAware{
			CbError: func(context web.Context, err error) web.Response {
				errorHappened = true
				return nil
			},
		},
		Template: "product/product",
		Router: &router.Router{
			RouterRegistry: router.NewRegistry(),
		},
	}
	u, _ := url.Parse(`http://test/`)
	vc.Router.SetBase(u)
	vc.Router.RouterRegistry.Route("/", `product.view(marketplacecode?="test", name?="test", variant?="test")`)
	ctx := web.NewContext()

	ctx.LoadParams(router.P{"marketplacecode": "test", "name": "testname"})
	response := vc.Get(ctx)

	expectedUrlTitle := "my-product-title"

	if redirectedTo != "product.view" {
		t.Errorf("Expected redirect to product.view, not %q", redirectedTo)
	}

	if redirectedName != expectedUrlTitle {
		t.Errorf("Expected redirect to name %s, not %q", expectedUrlTitle, redirectedName)
	}

	if response != nil {
		t.Errorf("Expected mocked response to be nil, not %T", response)
	}

	ctx.LoadParams(router.P{"marketplacecode": "test", "name": expectedUrlTitle})
	response = vc.Get(ctx)

	if errorHappened {
		t.Error("expected to not error for 'test' product")
	}

	if tplname != "product/product" {
		t.Errorf("expected to render product/product not %q", tplname)
	}

	if response != nil {
		t.Errorf("Expected mocked response to be nil, not %T", response)
	}

	ctx.LoadParams(router.P{"marketplacecode": "fail", "name": "fail"})
	response = vc.Get(ctx)

	if !errorHappened {
		t.Error("expected to error for 'fail' product")
	}

	if response != nil {
		t.Errorf("Expected mocked response to be nil, not %T", response)
	}

}
