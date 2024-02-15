package test

import (
	"fmt"
	"strings"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"
)

type MatchCallback[T any] func(actual T) (string, bool)

type genericMatcher[T any] struct {
	results        []string
	matchCallbacks []MatchCallback[T]
}

func (p genericMatcher[T]) Matches(x interface{}) bool {
	if msg, ok := x.(T); ok {
		for _, cb := range p.matchCallbacks {
			result, ok := cb(msg)
			p.results = append(p.results, result)
			if !ok {
				return false
			}
		}

		return true
	}
	return false
}

func (p genericMatcher[T]) String() string {
	return strings.Join(p.results, "\n")
}

func defaultProtoEqual(expected proto.Message) MatchCallback[proto.Message] {
	return func(actual proto.Message) (string, bool) {
		isEq := proto.Equal(actual, expected)
		if !isEq {
			return fmt.Sprintf("%v\nActual: %v", expected, actual), false
		}
		return "proto messages are equal", true
	}
}

func ProtoEq(x proto.Message) gomock.Matcher {
	return genericMatcher[proto.Message]{[]string{}, []MatchCallback[proto.Message]{defaultProtoEqual(x)}}
}

func MatchBy[T any](matchCallbacks ...MatchCallback[T]) gomock.Matcher {
	return genericMatcher[T]{[]string{}, matchCallbacks}
}
