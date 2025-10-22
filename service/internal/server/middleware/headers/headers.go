package headers

import (
	"context"
	"strings"

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

				if at := withBearer(accessToken); at != "" {
					hdr.Set("Authorization", at)
				}
				if rt := withBearer(refreshToken); rt != "" {
					hdr.Set("Refresh", rt)
				}

				// expose headers for browsers
				hdr.Add("Access-Control-Expose-Headers", "Authorization, Refresh")
			}

		case transport.KindGRPC:
			// Build metadata only with non-empty tokens
			pairs := []string{}
			if at := withBearer(accessToken); at != "" {
				// gRPC metadata keys are conventionally lowercase
				pairs = append(pairs, "authorization", at)
			}
			if rt := withBearer(refreshToken); rt != "" {
				pairs = append(pairs, "refresh", rt)
			}
			if len(pairs) > 0 {
				md := metadata.Pairs(pairs...)
				_ = grpc_std.SetHeader(ctx, md)
				_ = grpc_std.SetTrailer(ctx, md) // optional duplicate in trailer
			}
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
			return htr.Request().Header.Get("Refresh")
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

func withBearer(tok string) string {
	tok = strings.TrimSpace(tok)
	if tok == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(tok), "bearer ") {
		return tok
	}
	return "Bearer " + tok
}
