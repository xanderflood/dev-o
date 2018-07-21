A super fun Rails-y autoloading daemon for go servers. The basic idea is that `dev-o` will create a goroutine that monitors your go source files. Once it detects a change, `dev-o` begins attempting to rebuild the target package, and once it successfully rebuilds, `dev-o` will `exec  the current process into the new executable preserving command-line arguments and environment variables.

A minimal-ish example can be found [here](https://github.com/xanderflood/dev-o/blob/master/cmd/test.go).
