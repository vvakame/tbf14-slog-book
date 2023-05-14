package code

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

// range:sensitiveDef
var _ slog.LogValuer = SensitiveString("")

type SensitiveString string

func (SensitiveString) LogValue() slog.Value {
	return slog.StringValue("**censored**")
}

// range.end

// range:groupDef
var _ slog.LogValuer = (*User)(nil)

type User struct {
	ID       string
	Password SensitiveString
}

func (user *User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", user.ID),
		slog.Any("password", user.Password),
	)
}

// range.end

func Test_logValuer(t *testing.T) {
	defaultLogger := slog.Default()
	defer func() {
		slog.SetDefault(defaultLogger)
	}()
	h := slog.NewJSONHandler(os.Stdout, nil)

	logger := slog.New(h)
	slog.SetDefault(logger)

	ctx := context.Background()

	// range:emit
	slog.InfoCtx(
		ctx, "print LogValuer value",
		slog.Any("user", &User{
			ID:       "123",
			Password: "53c237",
		}),
	)
	// range.end

	// range:resolve
	v := slog.AnyValue(&User{
		ID:       "123",
		Password: "53c237",
	})
	// Kind = LogValuer
	t.Log(v.Kind())

	v = v.Resolve()
	// Kind = Group
	t.Log(v.Kind())
	// range.end
}
