package response_token

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	grpc_std "google.golang.org/grpc"
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
					hdr.Set("Refresh", refreshToken)
				}
				hdr.Add("Access-Control-Expose-Headers", "Authorization, Refresh")
			}
		case transport.KindGRPC:
			// gRPC metadata
			md := metadata.Pairs(
				"authorization", "Bearer "+accessToken,
			)
			if refreshToken != "" {
				md.Append("refresh", refreshToken)
			}
			// Send in headers
			_ = grpc_std.SetHeader(ctx, md)
			// (optional) duplicate in trailer
			_ = grpc_std.SetTrailer(ctx, md)
		}
	}
}
