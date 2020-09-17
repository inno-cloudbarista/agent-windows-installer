ifeq ($(OS),Windows_NT)
	delfile	:=	del
else
	delfile	:=	rm
endif

appname := cbinstaller
sources := $(wildcard *.go)

build = cd main && go env GOOS=$(1) GOARCH=$(2) && go build -o $(appname)$(3)
zip = zip $(1)_$(2).zip main/$(appname)$(3) ./build
rmexe = cd main && $(delfile) $(appname)$(1)

.PHONY: installer

##### WINDOWS BUILDS #####
installer: build/windows_386.zip build/windows_amd64.zip

build/windows_386.zip: $(sources)
	$(call build,cbinstaller_windows,386,.exe)
	$(call zip,cbinstaller_windows,386,.exe)
	$(call rmexe,.exe)

build/windows_amd64.zip: $(sources)
	$(call build,cbinstaller_windows,amd64,.exe)
	$(call zip,cbinstaller_windows,amd64,.exe)
	$(call rmexe,.exe)