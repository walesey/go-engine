COVER_DIR = cover

test:
	go test github.com/walesey/go-engine/actor
	go test github.com/walesey/go-engine/controller
	go test github.com/walesey/go-engine/effects
	go test github.com/walesey/go-engine/engine
	go test github.com/walesey/go-engine/physics/physicsAPI
	go test github.com/walesey/go-engine/ui
	go test github.com/walesey/go-engine/util
	go test github.com/walesey/go-engine/vectormath

coverage:
	mkdir -p $(COVER_DIR)
	go test github.com/walesey/go-engine/controller -coverprofile=$(COVER_DIR)/controller.cover.out && \
	go test github.com/walesey/go-engine/effects -coverprofile=$(COVER_DIR)/effects.cover.out && \
	go test github.com/walesey/go-engine/physics -coverprofile=$(COVER_DIR)/physics.cover.out && \
	go test github.com/walesey/go-engine/physics/gjk -coverprofile=$(COVER_DIR)/gjk.cover.out && \
	go test github.com/walesey/go-engine/util -coverprofile=$(COVER_DIR)/util.cover.out && \
	go test github.com/walesey/go-engine/vectormath -coverprofile=$(COVER_DIR)/vectormath.cover.out && \
		rm -f $(COVER_DIR)/coverage.out && \
		echo 'mode: set' > $(COVER_DIR)/coverage.out && \
		cat $(COVER_DIR)/*.cover.out | sed '/mode: set/d' >> $(COVER_DIR)/coverage.out && \
		go tool cover -html=$(COVER_DIR)/coverage.out -o=$(COVER_DIR)/coverage.html && \
		rm $(COVER_DIR)/*.cover.out
