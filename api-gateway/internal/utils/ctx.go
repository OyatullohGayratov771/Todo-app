package utils

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

func InjectTokenToContext(r *http.Request) (context.Context, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	var token string
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		token = authHeader // faqat token yuborilgan bo‘lsa ham ishlatamiz
	}

	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(r.Context(), md), nil
}
