A super fun Rails-like autoloading daemon that creates a .

Just call

```go
import (

	devo "github.com/xanderflood/dev-o"
)

func main() {
	// if we're in dev mode, start dev-o
	if len(os.Env("DEVO_DAEMON")) > 0 {
		devo.Autoreload(
			WithTarget("github.com/xanderflood/my-app/cmd/start"), //the `main` package to target
			WhileMonitoring( //directories to monitor
        "github.com/xanderflood/my-app/cmd",
        "github.com/xanderflood/my-app/lib",
        "github.com/xanderflood/my-app/vendor",
		)
	}
```
