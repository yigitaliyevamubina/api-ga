CURRENT_DIR=$(shell pwd)
DB_URL="postgres://postgres:mubina2007@db:5432/userdb?sslmode=disable"

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}
	
swag-gen:
	~/go/bin/swag init -g ./api/router.go -o api/docs

g:
	go run cmd/main.go

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-force:
	migrate -path migrations -database "$(DB_URL)" -verbose force 1

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_comments_table