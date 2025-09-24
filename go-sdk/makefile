.PHONY: spa # make spa
spa:
	cd ../frontend && quasar build
	cd ../go-sdk && rm -rf full/spa && cp -r ../frontend/dist/spa full/spa

.PHONY: android # make android
android: spa
	go get -d golang.org/x/mobile/cmd/gomobile
	gomobile bind  -ldflags '-s -w' -o ../frontend/src-capacitor/android/frontend/libs/full.aar -target=android -androidapi 22 full/android/android
	cd ../frontend && quasar build -m capacitor -T android
	

.PHONY: ios # make ios
ios: spa
# 	go get -d golang.org/x/mobile/cmd/gomobile
# 	gomobile bind  -ldflags '-s -w' -o ../frontend/src-capacitor/ios/full.framework -target=ios full/ios/ios
# 	cd ../frontend && quasar build -m capacitor -T ios

.PHONY: harmony # make harmony
harmony: spa
	rm -rf full/harmony/log-adaptor
	cd full/harmony && git clone https://github.com/harmony-contrib/log-adaptor.git
	cd full/harmony && ./build-win.sh
# 	cd ../frontend && quasar build -m electron && rm -rf dist/full main.go

	

.PHONY: win # make win
win: spa
	go build  -ldflags '-s -w' -o "../frontend/dist/full/full.exe"
	cd ../frontend && quasar build -m electron && rm -rf dist/full
	
.PHONY: linux # make linux
linux: spa
	go build  -ldflags '-s -w' -o "../frontend/dist/full/full"
	cd ../frontend && quasar build -m electron && rm -rf dist/full main.go

.PHONY: mac # make mac
mac: spa
	go build  -ldflags '-s -w' -o "../frontend/dist/full/full"
	cd ../frontend && quasar build -m electron && rm -rf dist/full main.go

	

# .PHONY: all # make all
# all: android win