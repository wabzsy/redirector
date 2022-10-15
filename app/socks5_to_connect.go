package app

//
//func (a *App) Socks2Connect() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:          "socks2connect",
//		Short:        "socks5 to connect",
//		Aliases:      []string{"s2c"},
//		SilenceUsage: true,
//	}
//
//	cmd.RunE = a.Socks2ConnectHandler
//	return cmd
//}
//
//func (a *App) Socks2ConnectHandler(_ *cobra.Command, args []string) error {
//	leftAddr := args[0]
//	rightAddr := args[0]
//	log.Println("Socks5 to Connect:", leftAddr, "--", rightAddr)
//	if server, err := socks5.New(&socks5.Config{}); err == nil {
//		go func(server *socks5.Server) {
//			if err = server.ListenAndServe("tcp", leftAddr); err != nil {
//				log.Println(err)
//			}
//		}(server)
//		a.Connect2ConnectHandler(h1, p1, h2, p2)
//	} else {
//		log.Println(err)
//	}
//	return nil
//}
