mockgen -source=player.go -destination player_mock.go -package interfaces 
mockgen -source=desk.go -destination desk_mock.go -package interfaces -imports steve_proto_gaterpc=steve/structs/proto/gate_rpc
mockgen -source=message_sender.go -destination message_sender_mock.go -package interfaces -imports steve_proto_gaterpc=steve/structs/proto/gate_rpc


