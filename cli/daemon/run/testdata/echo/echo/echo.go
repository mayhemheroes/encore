package echo

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	encore "encore.dev"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
)

type Data[K any, V any] struct {
	Key   K
	Value V
}

type NonBasicData struct {
	// Header
	HeaderString string `header:"X-Header-String"`
	HeaderNumber int    `header:"X-Header-Number"`

	// Body
	Struct       Data[*Data[string, string], int]
	StructPtr    *Data[int, uint16]
	StructSlice  []*Data[string, string]
	StructMap    map[string]*Data[string, float32]
	StructMapPtr *map[string]*Data[string, string]
	AnonStruct   struct{ AnonBird string }
	NamedStruct  *Data[string, float64] `json:"formatted_nest"`
	RawStruct    json.RawMessage

	// Query
	QueryString string `query:"string"`
	QueryNumber int    `query:"no"`

	// Path Parameters
	PathString string
	PathInt    int
	PathWild   string

	// Auth Parameters
	AuthHeader string
	AuthQuery  []int

	// Unexported fields
	unexported string
}

type EmptyData struct {
	OmitEmpty Data[string, string] `json:"OmitEmpty,omitempty"`
	NullPtr   *string
	Zero      Data[string, string]
}

type BasicData struct {
	String      string
	Uint        uint
	Int         int
	Int8        int8
	Int64       int64
	Float32     float32
	Float64     float64
	StringSlice []string
	IntSlice    []int
	Time        time.Time
}

type HeadersData struct {
	Int    int    `header:"X-Int"`
	String string `header:"X-String"`
}

// Echo echoes back the request data.
//encore:api public
func Echo(ctx context.Context, params *Data[string, int]) (*Data[string, int], error) {
	return params, nil
}

// EmptyEcho echoes back the request data.
//encore:api public
func EmptyEcho(ctx context.Context, params EmptyData) (EmptyData, error) {
	return params, nil
}

// NonBasicEcho echoes back the request data.
//encore:api auth path=/NonBasicEcho/:pathString/:pathInt/*pathWild
func NonBasicEcho(ctx context.Context, pathString string, pathInt int, pathWild string, params *NonBasicData) (*NonBasicData, error) {
	data := auth.Data().(*AuthParams)
	params.PathString = pathString
	params.PathInt = pathInt
	params.PathWild = pathWild
	params.AuthQuery = data.Query
	params.AuthHeader = data.Header
	return params, nil
}

// BasicEcho echoes back the request data.
//encore:api public method=GET,POST
func BasicEcho(ctx context.Context, params *BasicData) (*BasicData, error) {
	return params, nil
}

// HeadersEcho echoes back the request headers
//encore:api public method=GET,POST
func HeadersEcho(ctx context.Context, params *HeadersData) (*HeadersData, error) {
	return params, nil
}

// Noop does nothing
//encore:api public method=GET
func Noop(ctx context.Context) error {
	return nil
}

// MuteEcho absorbs a request
//encore:api public method=GET
func MuteEcho(ctx context.Context, params Data[string, string]) error {
	log.Printf("Absorbing %v\n", params)
	return nil
}

// Pong returns a bird tuple
//encore:api public method=GET
func Pong(ctx context.Context) (Data[string, string], error) {
	return Data[string, string]{"woodpecker", "kingfisher"}, nil
}

type EnvResponse struct {
	Env []string
}

// Env returns the environment.
//encore:api public
func Env(ctx context.Context) (*EnvResponse, error) {
	return &EnvResponse{Env: os.Environ()}, nil
}

type AppMetadata struct {
	AppID      string
	APIBaseURL string
	EnvName    string
	EnvType    string
}

// AppMeta returns app metadata.
//encore:api public
func AppMeta(ctx context.Context) (*AppMetadata, error) {
	md := encore.Meta()
	return &AppMetadata{
		AppID:      md.AppID,
		APIBaseURL: md.APIBaseURL.String(),
		EnvName:    md.Environment.Name,
		EnvType:    string(md.Environment.Type),
	}, nil
}

type AuthParams struct {
	Header        string `header:"X-Header"`
	Authorization string `header:"Authorization"`
	Query         []int  `query:"query"`
	NewAuth       bool   `query:"new-auth"`
}

//encore:authhandler
func AuthHandler(ctx context.Context, params *AuthParams) (auth.UID, *AuthParams, error) {
	if params.Authorization == "Bearer tokendata" && params.NewAuth == false {
		return "user", params, nil
	}

	// Check headers and query strings work by adding them together to calculate the answer in the header field
	if params.Header != "" && params.NewAuth {
		ans := 0
		for _, v := range params.Query {
			ans += v
		}

		if strconv.FormatInt(int64(ans), 10) == params.Header {
			return "second_user", params, nil
		}
	}

	return "", nil, &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "invalid token",
	}
}
