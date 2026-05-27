.PHONY: test smoke

test:
	GOWORK=off go test ./...

smoke:
	bash scripts/dev/public_smoke.sh
