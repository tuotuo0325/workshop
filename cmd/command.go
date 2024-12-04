package main

// Command 表示一个命令行命令
type Command interface {
	// Name 返回命令名称
	Name() string
	// Usage 返回命令用法说明
	Usage() string
	// Run 运行命令
	Run(args []string) error
}
