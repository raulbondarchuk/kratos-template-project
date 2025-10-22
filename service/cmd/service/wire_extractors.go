package main

import conf "service/internal/conf/v1"

func ProvideAppFromBootstrap(b *conf.Bootstrap) *conf.App {
	if b == nil {
		return nil
	}
	return b.App
}

func ProvideServerFromBootstrap(b *conf.Bootstrap) *conf.Server {
	if b == nil {
		return nil
	}
	return b.Server
}

func ProvideDataFromBootstrap(b *conf.Bootstrap) *conf.Data {
	if b == nil {
		return nil
	}
	return b.Data
}

func ProvideWebhooksFromBootstrap(b *conf.Bootstrap) *conf.Webhooks {
	if b == nil {
		return nil
	}
	return b.Webhooks
}
