mockgen -source=socket.go -destination mock_socket.go -package socket
mockgen -source=server.go -destination mock_server.go -package socket
mockgen -source=packer.go -destination mock_packer.go -package socket
