package price

import (
	"flamingo.me/dingo"
	pricegraphql "flamingo.me/flamingo-commerce/v3/price/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/price/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/core/locale"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/graphql"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "commercePriceFormat", new(templatefunctions.CommercePriceFormatFunc))
	injector.BindMulti(new(graphql.Service)).To(pricegraphql.Service{})
}

// Depends adds our dependencies
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(locale.Module),
	}
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	alias := make(map[string]string)
	for _, v := range []string{
		"locale.locale",
		"locale.accounting.default.precision",
		"locale.accounting.default.decimal",
		"locale.accounting.default.thousand",
		"locale.accounting.default.formatZero",
		"locale.accounting.default.format",
		"locale.accounting.default.formatLong",
		"locale.numbers.decimal",
		"locale.numbers.thousand",
		"locale.numbers.precision",
		"locale.date.dateFormat",
		"locale.date.timeFormat",
		"locale.date.dateTimeFormat",
		"locale.date.location",
	} {
		alias[v] = "core." + v
	}
	return alias
}
