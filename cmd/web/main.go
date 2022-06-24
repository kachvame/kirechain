package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/kachvame/kirechain/chain"
	"github.com/kachvame/kirechain/web"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	order := os.Getenv("ORDER")
	path := os.Getenv("PATH")

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed opening entries, %w", err)
	}

	orderNumber, err := strconv.Atoi(order)
	if err != nil {
		return fmt.Errorf("failed parsing order, %w", err)
	}

	h := web.New(nil)
	go func() {
		ch, err := chain.New(orderNumber, file)
		if err != nil {
			log.Println(fmt.Errorf("failed to build chain, %w", err))
			os.Exit(1)
		}
		h.Chain = &ch
	}()

	addr := ":8080"
	log.Println("kachvam na ", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		return err
	}
	return nil
}
