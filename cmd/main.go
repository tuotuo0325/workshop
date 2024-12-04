package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var commands = []Command{
	NewProducerCommand(),
	NewConsumerCommand(),
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// 设置信号处理
	setupSignalHandler()

	// 查找并执行命令
	cmdName := os.Args[1]
	for _, cmd := range commands {
		if cmd.Name() == cmdName {
			if err := cmd.Run(os.Args[2:]); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	// 未找到命令
	if cmdName == "help" || cmdName == "-h" || cmdName == "--help" {
		printUsage()
		return
	}

	fmt.Printf("Unknown command: %s\n", cmdName)
	printUsage()
	os.Exit(1)
}

func setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal. Shutting down gracefully...")
		os.Exit(0)
	}()
}

func printUsage() {
	fmt.Println("Usage: airbnb-cli <command> [options]")
	fmt.Println("\nCommands:")
	for _, cmd := range commands {
		fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Usage())
	}
	fmt.Println("  help       Show this help message")
	fmt.Println("\nProducer options:")
	fmt.Println("  --data string    Tasks data file (default \"tasks.json\")")
	fmt.Println("  --queue string   Queue server address (default \"file://queue.data\")")
	fmt.Println("\nConsumer options:")
	fmt.Println("  --workers int    Number of worker threads (default 10)")
	fmt.Println("  --queue string   Queue server address (default \"file://queue.data\")")
	fmt.Println("  --storage string Storage file path (default \"data/hotels.json\")")
}
