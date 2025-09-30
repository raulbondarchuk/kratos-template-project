package headers

import (
	"context"

	grpc_std "google.golang.org/grpc"

	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc/metadata"
)

// SetTokens installs tokens in the response for both transport types (HTTP/gRPC)
func SetTokens(ctx context.Context, accessToken, refreshToken string) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		switch tr.Kind() {
		case transport.KindHTTP:
			if htr, ok := tr.(*khttp.Transport); ok {
				hdr := htr.ReplyHeader()
				hdr.Set("Authorization", "Bearer "+accessToken)
				if refreshToken != "" {
					hdr.Set("Refresh", "Bearer "+refreshToken)
				}
				hdr.Add("Access-Control-Expose-Headers", "Authorization, Refresh")
			}
		case transport.KindGRPC:
			// gRPC metadata
			md := metadata.Pairs(
				"authorization", "Bearer "+accessToken,
			)
			if refreshToken != "" {
				md.Append("refresh", "Bearer "+refreshToken)
			}
			// Send in headers
			_ = grpc_std.SetHeader(ctx, md)
			// (optional) duplicate in trailer
			_ = grpc_std.SetTrailer(ctx, md)
		}
	}
}

func GetAccessTokenFromHeader(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if htr, ok := tr.(*khttp.Transport); ok && htr.Request() != nil {
			return htr.Request().Header.Get("Authorization")
		}
	}
	return ""
}

func GetRefreshTokenFromHeader(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if htr, ok := tr.(*khttp.Transport); ok && htr.Request() != nil {
			return htr.Request().Header.Get("Authorization")
		}
	}
	return ""
}

func GetSecretFromHeader(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		if htr, ok := tr.(*khttp.Transport); ok && htr.Request() != nil {
			return htr.Request().Header.Get("X-Secret-Access")
		}
	}
	return ""
}
