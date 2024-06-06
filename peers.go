package GeeCache

import pb "GeeCache/geecachepb"

// PeerPicker 根据输入的key返回对应的PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 客户端
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
