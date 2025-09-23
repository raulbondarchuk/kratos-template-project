package meta

import "service/api/common/v1"

func MetaOK(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_OK, Message: msg}
}
func MetaNoContent(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_NO_CONTENT, Message: msg}
}
func MetaBadRequest(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_BAD_REQUEST, Message: msg}
}
func MetaUnauthorized(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_UNAUTHORIZED, Message: msg}
}
func MetaForbidden(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_FORBIDDEN, Message: msg}
}
func MetaNotFound(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_NOT_FOUND, Message: msg}
}
func MetaConflict(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_CONFLICT, Message: msg}
}
func MetaUnprocessable(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_UNPROCESSABLE_ENTITY, Message: msg}
}
func MetaInternal(msg string) *common.MetaResponse {
	return &common.MetaResponse{Code: common.ResponseCode_RESPONSE_CODE_INTERNAL_SERVER_ERROR, Message: msg}
}

// With details (to avoid creating map each time)
func WithDetails(meta *common.MetaResponse, kv map[string]string) *common.MetaResponse {
	if meta == nil {
		meta = &common.MetaResponse{}
	}
	if meta.Details == nil {
		meta.Details = map[string]string{}
	}
	for k, v := range kv {
		meta.Details[k] = v
	}
	return meta
}
