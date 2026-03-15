package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/HelenaBlack/anti-bruteforce/api/gen" //nolint:depguard
	"github.com/spf13/cobra"                            //nolint:depguard
	"google.golang.org/grpc"                            //nolint:depguard
	"google.golang.org/grpc/credentials/insecure"       //nolint:depguard
)

var (
	addr   string
	conn   *grpc.ClientConn
	client pb.AntibruteforceClient
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "admin",
		Short: "Anti-Bruteforce Admin CLI",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			var err error
			conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			client = pb.NewAntibruteforceClient(conn)
		},
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			if conn != nil {
				_ = conn.Close()
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&addr, "addr", "localhost:50051", "server address")

	// Reset Command
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset buckets for login and IP",
		Run: func(cmd *cobra.Command, _ []string) {
			login, _ := cmd.Flags().GetString("login")
			ip, _ := cmd.Flags().GetString("ip")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.Reset(ctx, &pb.ResetRequest{Login: login, Ip: ip})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("Buckets reset successfully")
		},
	}
	resetCmd.Flags().String("login", "", "Login to reset")
	resetCmd.Flags().String("ip", "", "IP to reset")

	// Blacklist Commands
	blacklistCmd := &cobra.Command{Use: "blacklist", Short: "Manage blacklist"}
	blAddCmd := &cobra.Command{
		Use:  "add [subnet]",
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.AddToBlacklist(ctx, &pb.SubnetRequest{Subnet: args[0]})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Subnet %s added to blacklist\n", args[0])
		},
	}
	blRemoveCmd := &cobra.Command{
		Use:  "remove [subnet]",
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.RemoveFromBlacklist(ctx, &pb.SubnetRequest{Subnet: args[0]})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Subnet %s removed from blacklist\n", args[0])
		},
	}
	blacklistCmd.AddCommand(blAddCmd, blRemoveCmd)

	// Whitelist Commands
	whitelistCmd := &cobra.Command{Use: "whitelist", Short: "Manage whitelist"}
	wlAddCmd := &cobra.Command{
		Use:  "add [subnet]",
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.AddToWhitelist(ctx, &pb.SubnetRequest{Subnet: args[0]})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Subnet %s added to whitelist\n", args[0])
		},
	}
	wlRemoveCmd := &cobra.Command{
		Use:  "remove [subnet]",
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.RemoveFromWhitelist(ctx, &pb.SubnetRequest{Subnet: args[0]})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("Subnet %s removed from whitelist\n", args[0])
		},
	}
	whitelistCmd.AddCommand(wlAddCmd, wlRemoveCmd)

	rootCmd.AddCommand(resetCmd, blacklistCmd, whitelistCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
