GOCMD=go
GOBUILD=$(GOCMD) build
DOCKER=sudo docker
APP=nrm
GOFMT=gofmt -w
GOLINT=golint

define makeit 
	$(GOFMT) $(1)/*.go
	$(GOLINT) $(1)/*.go
	$(GOBUILD) $(1)/*.go
endef


build:
	$(call makeit,vpc)
	$(call makeit,nbapi)
	$(call makeit,config)
	$(call makeit,sbapi)
	$(call makeit,psm)
	$(call makeit,apiutils)
	$(GOBUILD) 
	$(DOCKER) build -t $(APP):latest .

run:
	$(DOCKER) run -d -p 8080:8080 -it $(APP) 

clean:
	rm -f $(APP)
