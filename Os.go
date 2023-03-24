package et

import (
	"os"
	"runtime"
)

const OS_WINDOWS = "windows"
const OS_LINUX = "linux"
const OS_DARWIN = "darwin"

func OsName() string{
	return runtime.GOOS
}

func OsEnvVar(strKey string) string{
	return os.Getenv(strKey)
}
