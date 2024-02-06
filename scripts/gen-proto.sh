CURRENT_DIR=$(pwd)

protoc -I /usr/local/include \
       -I $GOPATH/src/github.com/gogo/protobuf/gogoproto \
       -I "$CURRENT_DIR/protos/" \
        --gofast_out=plugins=grpc:"$CURRENT_DIR/genproto/" \
        "$CURRENT_DIR/protos/comment_service/comment.proto";