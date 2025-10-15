package authorization

import (
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewEnforcer)

func NewEnforcer() (*casbin.Enforcer, error) {
	// 获取项目根目录 - 从 internal/authorization/provider.go 到 lushop 项目根目录
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "../..")

	modelPath := filepath.Join(projectRoot, "internal/pkg/middleware/casbin/model.conf")
	policyPath := filepath.Join(projectRoot, "internal/pkg/middleware/casbin/policy.csv")

	e, err := casbin.NewEnforcer(
		modelPath,
		fileadapter.NewAdapter(policyPath),
	)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}
	return e, nil
}
