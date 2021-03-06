// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mindreader

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"go.uber.org/zap"
)

func RunGRPCServer(s *grpc.Server, listenAddr string) error {
	zlog.Info("starting grpc listener", zap.String("listen_addr", listenAddr))
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", listener)
	}

	serverError := make(chan error, 1)

	go func() {
		if err := s.Serve(listener); err != nil {
			serverError <- err
		}
		zlog.Info("grpc server terminated")
	}()

	select {
	case <-time.After(1 * time.Second):
		zlog.Info("grpc server listener ready")
	case err := <-serverError:
		return err
	}

	return nil
}
