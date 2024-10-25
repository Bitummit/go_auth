package main

import (
	run "github.com/Bitummit/go_auth/internal"
)


func main() {
	run.Run()
}


// go func() {
// 	defer wg.Done()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			grpcServer.GracefulStop()
// 		default:
// 			if err = grpcServer.Serve(listener); err != nil {
// 				server.Log.Error("error starting server", logger.Err(err))
// 			}
// 		}
// 	}
// }()