
DOCKER_IMAGE?=xaionaro2/streampanel-android-builder
DOCKER_CONTAINER_NAME?=CameraView-builder
GOBUILD_FLAGS?=
LINKER_FLAGS?=
LINKER_FLAGS_ANDROID?=$(LINKER_FLAGS)
GOPATH?=$(shell go env GOPATH)

go-run:
	go run ./

CameraView.apk: CameraView-unsigned.apk sign-CameraView.apk

CameraView-unsigned.apk:
	docker pull  $(DOCKER_IMAGE)
	docker start $(DOCKER_CONTAINER_NAME) >/dev/null 2>&1 || \
		( cd ../../ && docker run \
			--detach \
			--init \
			--name $(DOCKER_CONTAINER_NAME) \
			--volume ".:/project" \
			--tty \
			$(DOCKER_IMAGE) >/dev/null 2>&1 ) || /bin/true
	docker exec $(DOCKER_CONTAINER_NAME) make -C /project/cmd/CameraView CameraView-unsigned.apk-in-docker

CameraView-unsigned.apk-in-docker: $(GOPATH)/bin/pkg-config-wrapper
	go mod tidy || /bin/true
	git config --global --add safe.directory /project
	$(eval ANDROID_NDK_HOME=$(shell ls -d /home/builder/lib/android-ndk-* | tail -1))
	PKG_CONFIG_WRAPPER_LOG='/tmp/pkg_config_wrapper.log' \
	PKG_CONFIG_WRAPPER_LOG_LEVEL='trace' \
	PKG_CONFIG_LIBS_FORCE_STATIC='libav*' \
	PKG_CONFIG_ERASE="-fopenmp=*,-landroid,-lcamera2ndk,-lmediandk" \
	PKG_CONFIG='$(GOPATH)/bin/pkg-config-wrapper' \
	PKG_CONFIG_PATH='/data/data/com.termux/files/usr/lib/pkgconfig' \
	CGO_CFLAGS='-I$(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/sysroot/usr/include/ -I/data/data/com.termux/files/usr/include -Wno-incompatible-function-pointer-types -Wno-unused-result -Wno-xor-used-as-pow' \
	CGO_LDFLAGS='-ldl -lc -lcamera2ndk -lmediandk -L$(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/sysroot/usr/lib/ -L$(shell ls -d $(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/sysroot/usr/lib/aarch64-linux-android/[0-9]* | tail -1)/ -L/data/data/com.termux/files/usr/lib' \
	ANDROID_NDK_HOME="$(ANDROID_NDK_HOME)" \
	PATH="${PATH}:${HOME}/go/bin" \
	GOFLAGS="$(GOBUILD_FLAGS) -ldflags=$(shell echo ${LINKER_FLAGS_ANDROID} | tr " " ",")" \
	fyne package $(FYNEBUILD_FLAGS) -release -os android/arm64
	mv CameraView.apk CameraView-unsigned.apk

$(GOPATH)/bin/pkg-config-wrapper:
	go install github.com/xaionaro-go/pkg-config-wrapper@5dd443e6c18336416c49047e2ba0002e26a85278

priv/android-apk.keystore:
	mkdir -p priv
	keytool -genkey -v -keystore priv/android-apk.keystore -alias CameraView -keyalg RSA -keysize 2048 -validity 36500

sign-CameraView.apk: priv/android-apk.keystore
	cp CameraView-unsigned.apk CameraView-signing.apk
	zipalign -p 4 CameraView-signing.apk CameraView-signing-aligned.apk
	rm -f CameraView-signing.apk 
	mv CameraView-signing-aligned.apk CameraView-signing.apk
	zipalign -c 4 CameraView-signing.apk
	apksigner sign --ks-key-alias CameraView --ks priv/android-apk.keystore CameraView-signing.apk
	apksigner verify CameraView-signing.apk
	mv CameraView-signing.apk CameraView.apk

