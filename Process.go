package et

import "os"

func ProcessCurrentId() int{

	return os.Getpid()
}

func ProcessParentId() int{

	return os.Getppid()
}

