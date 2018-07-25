mockgen -source=flow.go -destination flow_mock.go -package interfaces 
mockgen -source=state.go -destination state_mock.go -package interfaces 
mockgen -source=transition.go -destination transition_mock.go -package interfaces 