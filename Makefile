build:
	vgo build github.com/cention-mujibur-rahman/msgoraph
	vgo build github.com/cention-mujibur-rahman/msgoraph/client
	vgo build github.com/cention-mujibur-rahman/msgoraph/common
	vgo build github.com/cention-mujibur-rahman/msgoraph/internal
	vgo build github.com/cention-mujibur-rahman/msgoraph/scopes
	vgo build github.com/cention-mujibur-rahman/msgoraph/users

docs:
	@echo "http://localhost:6060/pkg/github.com/cention-mujibur-rahman/msgoraph/"
	godoc -http=:6060
